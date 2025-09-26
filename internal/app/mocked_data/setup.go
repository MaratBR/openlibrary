package mockeddata

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"runtime"
	"strings"
	"sync"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/gofrs/uuid"
)

type Setup struct {
	tagsService        app.TagsService
	reviewsService     app.ReviewsService
	bookManagerService app.BookManagerService
	authService        app.AuthService
}

func NewSetup(tagsService app.TagsService, reviewsService app.ReviewsService, bookManagerService app.BookManagerService, authService app.AuthService) Setup {
	return Setup{
		tagsService:        tagsService,
		reviewsService:     reviewsService,
		bookManagerService: bookManagerService,
		authService:        authService,
	}
}

type SetupOptions struct {
	Users         int
	BooksLocation string
}

func (s *Setup) Run(options SetupOptions) error {
	if options.BooksLocation == "" {
		slog.Warn("no books location provided")
	}
	if options.Users == 0 {
		options.Users = 50
	}

	var (
		err     error
		userIds []uuid.UUID
	)

	slog.Info("creating users...", "count", options.Users)
	if userIds, err = CreateUsers(s.authService, options.Users); err != nil {
		return err
	}

	slog.Info("creating a few lorem ipsum books that are very long")
	err = s.createLongBook(userIds)
	if err != nil {
		return err
	}

	slog.Info("importing rr books")

	err = massImport(context.Background(), options.BooksLocation, userIds, s.bookManagerService, s.tagsService, runtime.NumCPU())
	if err != nil {
		return err
	}

	return nil
}

func (s *Setup) createLongBook(userIds []uuid.UUID) error {
	var wg sync.WaitGroup
	for _, size := range []int{100, 200, 500, 1000, 2000, 5000} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bookName := fmt.Sprintf("Very long book with %d chapters", size)
			bookID, err := s.bookManagerService.CreateBook(context.Background(), app.CreateBookCommand{
				Name:              bookName,
				UserID:            (userIds[rand.Int()%len(userIds)]),
				Tags:              nil,
				IsPubliclyVisible: true,
				AgeRating:         app.AgeRatingPG,
				Summary:           fmt.Sprintf("<p>This is an auto generated book with %d chapters</p><p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>", size),
			})

			if err != nil {
				slog.Error("failed to create book", "err", err)
				return
			}

			var longContent string

			{
				sb := strings.Builder{}

				for i := 0; i < size; i++ {
					sb.WriteString(fmt.Sprintf("<p><strong>%d:&nbsp;</strong>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>", i+1))
				}

				longContent = sb.String()
			}

			for i := 0; i < size; i++ {
				_, err := s.bookManagerService.CreateBookChapter(context.Background(), app.CreateBookChapterCommand{
					BookID:  bookID,
					Name:    fmt.Sprintf("Chapter %d", i+1),
					Summary: fmt.Sprintf("Summary of chapter %d", i+1),
					Content: longContent,
				})

				if err != nil {
					slog.Error("failed to create chapter", "err", err)
				}
			}

			// create some reviews for this book
			slog.Info("creating reviews for the long book", "name", bookName, "chapters", size)
			_ = s.createBookReviews(userIds, bookID, size/100)
		}()
	}

	wg.Wait()

	return nil
}
func (s *Setup) createBookReviews(userIds []uuid.UUID, bookID int64, maxReviews int) error {
	var wg sync.WaitGroup
	for i := 0; i < maxReviews; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := s.reviewsService.UpdateReview(context.Background(), app.UpdateReviewCommand{
				BookID:  bookID,
				UserID:  userIds[rand.Int()%len(userIds)],
				Rating:  app.RatingValue(rand.IntN(10) + 1),
				Content: "I really enjoyed reading this book. Highly recommended!",
			})

			if err != nil {
				slog.Error("failed to create review", "err", err)
				return
			}
		}()
	}

	wg.Wait()

	return nil
}
