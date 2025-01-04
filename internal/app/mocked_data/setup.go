package mockeddata

import (
	"log/slog"

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
		bookIds []int64
	)

	slog.Info("creating users...", "count", options.Users)
	if userIds, err = CreateUsers(s.authService, options.Users); err != nil {
		return err
	}

	slog.Info("importing books from ao3...", "location", options.BooksLocation)
	if bookIds, err = MassImportAo3(s.bookManagerService, s.tagsService, options.BooksLocation, userIds); err != nil {
		return err
	}
	for _, bookId := range bookIds {
		if err := CreateReviews(s.reviewsService, userIds, bookId); err != nil {
			return err
		}
	}

	return nil
}
