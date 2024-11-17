package store

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Int4Range struct {
	Min pgtype.Int4
	Max pgtype.Int4
}

func applyRange(builder *goqu.SelectDataset, column exp.Comparable, int4Range Int4Range) *goqu.SelectDataset {
	if int4Range.Min.Valid {
		builder = builder.Where(column.Gte(int4Range.Min.Int32))
	}
	if int4Range.Max.Valid {
		builder = builder.Where(column.Lte(int4Range.Max.Int32))
	}

	return builder
}

type BookSearchFilter struct {
	Words           Int4Range
	WordsPerChapter Int4Range
	Chapters        Int4Range
	Favorites       Int4Range

	IncludeAuthors []pgtype.UUID
	ExcludeAuthors []pgtype.UUID

	IncludeParentTags []int64
	ExcludeParentTags []int64

	IncludeBanned bool
	IncludeHidden bool
	IncludeEmpty  bool
}

type BookSearchRequest struct {
	BookSearchFilter
	Limit  uint
	Offset uint
}

type BookSearchRow struct {
	ID                 int64              `db:"id"`
	Name               string             `db:"name"`
	Summary            string             `db:"summary"`
	Favorites          int32              `db:"favorites"`
	CreatedAt          pgtype.Timestamptz `db:"created_at"`
	AgeRating          AgeRating          `db:"age_rating"`
	Words              int                `db:"words"`
	Chapters           int                `db:"chapters"`
	TagIds             []int64            `db:"tag_ids"`
	CachedParentTagIds []int64            `db:"cached_parent_tag_ids"`
	AuthorUserID       pgtype.UUID        `db:"author_user_id"`
	AuthorName         string             `db:"author_name"`
}

func applyWhere(query *goqu.SelectDataset, filter *BookSearchFilter) *goqu.SelectDataset {
	query = applyRange(query, goqu.I("books.words"), filter.Words)
	query = applyRange(query, goqu.I("books.chapters"), filter.Chapters)
	query = applyRange(query, goqu.I("books.favorites"), filter.Favorites)
	query = applyRange(query, goqu.L("cast(words as real) / chapters"), filter.WordsPerChapter)

	if len(filter.IncludeAuthors) > 0 {
		query = query.Where(goqu.C("author_user_id").In(filter.IncludeAuthors))
	}

	if len(filter.ExcludeAuthors) > 0 {
		query = query.Where(goqu.C("author_user_id").NotIn(filter.ExcludeAuthors))
	}

	if len(filter.IncludeParentTags) > 0 {
		query = query.Where(goqu.Literal("cached_parent_tag_ids @> ?::int8[]", PGArrayExpr(filter.IncludeParentTags)))
	}

	if len(filter.ExcludeParentTags) > 0 {
		query = query.Where(goqu.Literal("not (cached_parent_tag_ids && ?::int8[])", PGArrayExpr(filter.ExcludeParentTags)))
	}

	if !filter.IncludeBanned {
		query = query.Where(goqu.C("is_banned").IsFalse())
	}

	if !filter.IncludeHidden {
		query = query.Where(goqu.C("is_publicly_visible").IsTrue())
	}

	if !filter.IncludeEmpty {
		query = query.Where(goqu.C("chapters").Gt(0))
	}

	return query
}

var (
	postgresQuery = goqu.Dialect("postgres")
)

func createBookSearchSelect(req *BookSearchRequest) *goqu.SelectDataset {
	if req == nil {
		panic("BookSearchRequest is nil")
	}

	query := postgresQuery.
		Select(
			goqu.I("books.id"),
			goqu.I("books.name"),
			goqu.I("books.summary"),
			goqu.I("books.created_at"),
			goqu.I("books.age_rating"),
			goqu.I("books.words"),
			goqu.I("books.chapters"),
			goqu.I("books.tag_ids"),
			goqu.I("books.favorites"),
			goqu.I("books.cached_parent_tag_ids"),
			goqu.I("books.author_user_id"),
			goqu.I("author.name").As("author_name"),
		).
		From("books").
		Join(goqu.T("users").As("author"), goqu.On(
			goqu.I("books.author_user_id").Eq(goqu.I("author.id"))))

	query = applyWhere(query, &req.BookSearchFilter)
	query = query.Limit(req.Limit).Offset(req.Offset)

	return query
}

func SearchBooks(ctx context.Context, db DBTX, req BookSearchRequest) ([]BookSearchRow, error) {
	sql, params, err := createBookSearchSelect(&req).ToSQL()
	if err != nil {
		return nil, err
	}

	println(sql)
	rows, err := db.Query(ctx, sql, params...)
	defer rows.Close()
	if err != nil {
		slog.Error("failed to execute search query", "err", err, "sql", sql)
		return nil, err
	}

	var books []BookSearchRow
	for rows.Next() {
		row, err := pgx.RowToStructByName[BookSearchRow](rows)
		if err != nil {
			return nil, err
		}
		books = append(books, row)
	}
	if err != nil {
		return nil, err
	}

	if books == nil {
		return []BookSearchRow{}, nil
	}
	return books, nil
}

type BookStats struct {
	ChaptersMin        int64 `db:"chapters_min"`
	ChaptersMax        int64 `db:"chapters_max"`
	WordsMin           int64 `db:"words_min"`
	WordsMax           int64 `db:"words_max"`
	WordsPerChapterMin int64 `db:"words_per_chapter_min"`
	WordsPerChapterMax int64 `db:"words_per_chapter_max"`
	FavoritesMin       int64 `db:"favorites_min"`
	FavoritesMax       int64 `db:"favorites_max"`
}

func GetBooksFilterExtremes(ctx context.Context, db DBTX, req *BookSearchFilter) (BookStats, error) {
	query := postgresQuery.
		Select(
			goqu.MAX(goqu.I("books.chapters")).As("chapters_max"),
			goqu.MIN(goqu.I("books.chapters")).As("chapters_min"),

			goqu.MAX(goqu.I("books.words")).As("words_max"),
			goqu.MIN(goqu.I("books.words")).As("words_min"),

			goqu.MAX(goqu.I("books.favorites")).As("favorites_max"),
			goqu.MIN(goqu.I("books.favorites")).As("favorites_min"),

			goqu.MAX(goqu.L("ceil(cast(books.words as real) / books.chapters)")).As("words_per_chapter_max"),
			goqu.MIN(goqu.L("floor(cast(books.words as real) / books.chapters)")).As("words_per_chapter_min"),
		).
		From("books")
	query = applyWhere(query, req)

	sql, params, err := query.ToSQL()
	if err != nil {
		return BookStats{}, err
	}

	rows, err := db.Query(ctx, sql, params...)
	defer rows.Close()
	if err != nil {
		println("err1")

		return BookStats{}, err
	}
	if rows.Next() {
		stats, err := pgx.RowToStructByName[BookStats](rows)
		if err != nil {
			fmt.Printf("ERR : %s\n", err.Error())

			return BookStats{}, err
		}
		fmt.Printf("OK : %v\n", stats)
		return stats, nil
	} else {
		return BookStats{}, pgx.ErrNoRows
	}
}
