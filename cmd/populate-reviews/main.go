package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type options struct {
	count          int
	seed           int64
	databaseURL    string
	configPath     string
	includePrivate bool
	dryRun         bool
}

type user struct {
	id uuid.UUID
}

type book struct {
	id int64
}

type pair struct {
	userIndex int
	bookIndex int
}

func main() {
	opts := parseFlags()
	if opts.count < 1 {
		log.Fatal("count must be greater than zero")
	}
	if opts.seed == 0 {
		opts.seed = time.Now().UnixNano()
	}

	databaseURL, err := resolveDatabaseURL(opts)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	db, err := store.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer db.Close()

	users, err := loadUsers(ctx, db)
	if err != nil {
		log.Fatalf("load users: %v", err)
	}
	books, err := loadBooks(ctx, db, opts.includePrivate)
	if err != nil {
		log.Fatalf("load books: %v", err)
	}
	if len(users) == 0 || len(books) == 0 {
		log.Fatal("database must contain at least one eligible user and book")
	}

	rng := rand.New(rand.NewSource(opts.seed))
	pairs := selectPairs(rng, users, books, opts.count)
	if len(pairs) == 0 {
		log.Fatal("no eligible user/book pairs found")
	}
	if len(pairs) < opts.count {
		log.Printf("requested %d reviews, but only %d unique eligible user/book pairs exist", opts.count, len(pairs))
	}

	if opts.dryRun {
		log.Printf("dry run: would upsert %d reviews (seed %d)", len(pairs), opts.seed)
		return
	}

	if err := populate(ctx, db, rng, users, books, pairs); err != nil {
		log.Fatalf("populate reviews: %v", err)
	}
	log.Printf("upserted %d reviews across %d books (seed %d)", len(pairs), distinctBooks(pairs, books), opts.seed)
}

func parseFlags() options {
	var opts options
	flag.IntVar(&opts.count, "count", 250, "number of unique user/book reviews to upsert")
	flag.Int64Var(&opts.seed, "seed", 0, "random seed; zero uses the current time")
	flag.StringVar(&opts.databaseURL, "database-url", os.Getenv("DATABASE_URL"), "PostgreSQL URL; defaults to DATABASE_URL or the config file")
	flag.StringVar(&opts.configPath, "config", "openlibrary.toml", "path to the base TOML configuration")
	flag.BoolVar(&opts.includePrivate, "include-private", false, "include hidden, banned, trashed, and empty books")
	flag.BoolVar(&opts.dryRun, "dry-run", false, "validate selection without writing to the database")
	flag.Parse()
	return opts
}

func resolveDatabaseURL(opts options) (string, error) {
	if strings.TrimSpace(opts.databaseURL) != "" {
		return opts.databaseURL, nil
	}

	config := koanf.New(".")
	if err := config.Load(file.Provider(opts.configPath), toml.Parser()); err != nil {
		return "", fmt.Errorf("load config %q: %w", opts.configPath, err)
	}
	privatePath := strings.TrimSuffix(opts.configPath, ".toml") + ".private.toml"
	if _, err := os.Stat(privatePath); err == nil {
		if err := config.Load(file.Provider(privatePath), toml.Parser()); err != nil {
			return "", fmt.Errorf("load private config %q: %w", privatePath, err)
		}
	}
	url := config.String("database.url")
	if strings.TrimSpace(url) == "" {
		return "", fmt.Errorf("database.url is empty")
	}
	return url, nil
}

type queryer interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

func loadUsers(ctx context.Context, db queryer) ([]user, error) {
	rows, err := db.Query(ctx, `select id from users where not is_banned and role <> 'system' order by id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var value user
		if err := rows.Scan(&value.id); err != nil {
			return nil, err
		}
		users = append(users, value)
	}
	return users, rows.Err()
}

func loadBooks(ctx context.Context, db queryer, includePrivate bool) ([]book, error) {
	query := `select id from books`
	if !includePrivate {
		query += ` where is_publicly_visible and not is_banned and not is_trashed and chapters > 0`
	}
	query += ` order by id`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []book
	for rows.Next() {
		var value book
		if err := rows.Scan(&value.id); err != nil {
			return nil, err
		}
		books = append(books, value)
	}
	return books, rows.Err()
}

func selectPairs(rng *rand.Rand, users []user, books []book, requested int) []pair {
	maxPairs := len(users) * len(books)
	if requested > maxPairs {
		requested = maxPairs
	}

	selectedIndexes := make(map[int]struct{}, requested)
	for candidateMax := maxPairs - requested; candidateMax < maxPairs; candidateMax++ {
		candidate := rng.Intn(candidateMax + 1)
		if _, exists := selectedIndexes[candidate]; exists {
			selectedIndexes[candidateMax] = struct{}{}
		} else {
			selectedIndexes[candidate] = struct{}{}
		}
	}

	result := make([]pair, 0, len(selectedIndexes))
	for index := range selectedIndexes {
		result = append(result, pair{
			userIndex: index / len(books),
			bookIndex: index % len(books),
		})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].bookIndex == result[j].bookIndex {
			return result[i].userIndex < result[j].userIndex
		}
		return result[i].bookIndex < result[j].bookIndex
	})
	return result
}

type beginner interface {
	Begin(context.Context) (pgx.Tx, error)
}

func populate(ctx context.Context, db beginner, rng *rand.Rand, users []user, books []book, pairs []pair) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // no-op after Commit

	affectedBooks := make(map[int64]struct{})
	for _, selected := range pairs {
		userID := users[selected.userIndex].id
		bookID := books[selected.bookIndex].id
		rating := int16(rng.Intn(10) + 1)
		content := generateReview(rng, rating)
		createdAt := randomPastTime(rng)

		if _, err := tx.Exec(ctx, `
            insert into reviews (user_id, book_id, content, created_at)
            values ($1, $2, $3, $4)
            on conflict (user_id, book_id) do update
            set content = excluded.content, last_updated_at = now()`, userID, bookID, content, createdAt); err != nil {
			return fmt.Errorf("upsert review for user %s and book %d: %w", userID, bookID, err)
		}
		if _, err := tx.Exec(ctx, `
            insert into ratings (user_id, book_id, rating, updated_at)
            values ($1, $2, $3, $4)
            on conflict (user_id, book_id) do update
            set rating = excluded.rating, updated_at = now()`, userID, bookID, rating, createdAt); err != nil {
			return fmt.Errorf("upsert rating for user %s and book %d: %w", userID, bookID, err)
		}
		affectedBooks[bookID] = struct{}{}
	}

	bookIDs := make([]int64, 0, len(affectedBooks))
	for bookID := range affectedBooks {
		bookIDs = append(bookIDs, bookID)
	}
	if _, err := tx.Exec(ctx, `
        update books
        set rating = (select avg(rating::float8) from ratings where ratings.book_id = books.id),
            total_ratings = (select count(*) from ratings where ratings.book_id = books.id),
            total_reviews = (select count(*) from reviews where reviews.book_id = books.id)
        where id = any($1::int8[])`, bookIDs); err != nil {
		return fmt.Errorf("recalculate affected books: %w", err)
	}

	return tx.Commit(ctx)
}

func randomPastTime(rng *rand.Rand) time.Time {
	const twoYears = int64(2 * 365 * 24 * time.Hour)
	return time.Now().Add(-time.Duration(rng.Int63n(twoYears)))
}

func generateReview(rng *rand.Rand, rating int16) string {
	stars := float64(rating) / 2
	short := []string{
		"It was fine. Nothing groundbreaking, but I had a decent time with it.",
		fmt.Sprintf("Okay book, %.1f/5.", stars),
		"A quick, readable story. I liked some parts more than others.",
		"Not bad. I would probably read another book by this author.",
	}
	negative := []string{
		"The premise had potential, but the pacing dragged and the characters never felt fully developed.",
		"I wanted to like this more. The early chapters set up interesting ideas, but the payoff did not work for me.",
		"The prose is readable, yet the plot relies too heavily on repetition and convenient decisions.",
		"There are a few strong scenes, but they are buried under uneven pacing and dialogue that often feels forced.",
	}
	mixed := []string{
		"This is an uneven but enjoyable read. The middle slows down, although the ending recovers nicely.",
		"The world and central idea are interesting, while the characterization could use more depth.",
		"Some chapters really worked for me and others felt rushed. Overall, it lands somewhere in the middle.",
		"A solid casual read with a few memorable moments. It does not always come together, but it kept me reading.",
	}
	positive := []string{
		"A genuinely engaging story with confident pacing and characters I quickly became invested in.",
		"The strongest part is the character work. Even the quieter chapters move the relationships forward.",
		"Well written and easy to get lost in. The story balances momentum with enough room for the cast to breathe.",
		"The opening hooked me, and the book maintained that energy without sacrificing emotional weight.",
	}
	glowing := []string{
		"Excellent from beginning to end. The voice is distinctive, the cast feels alive, and the payoff is earned.",
		"One of the most satisfying reads I have had recently. I finished it and immediately wanted more.",
		"The kind of story that makes it difficult to stop after one chapter. Thoughtful, exciting, and beautifully paced.",
		"A standout book with memorable characters and a remarkably strong sense of direction.",
	}

	if rng.Intn(4) == 0 {
		return "<p>" + short[rng.Intn(len(short))] + "</p>"
	}
	var pool []string
	switch {
	case rating <= 3:
		pool = negative
	case rating <= 6:
		pool = mixed
	case rating <= 8:
		pool = positive
	default:
		pool = glowing
	}
	first := pool[rng.Intn(len(pool))]
	second := []string{
		"The editing is mostly clean, and the chapter length makes it easy to keep going.",
		"I especially appreciated that the story takes its time when a scene needs it.",
		"Your mileage may vary, but the book has a clear identity and commits to it.",
		"I would recommend trying the first few chapters to see whether the style clicks.",
	}[rng.Intn(4)]
	return "<p>" + first + "</p><p>" + second + "</p>"
}

func distinctBooks(pairs []pair, books []book) int {
	values := make(map[int64]struct{})
	for _, selected := range pairs {
		values[books[selected.bookIndex].id] = struct{}{}
	}
	return len(values)
}
