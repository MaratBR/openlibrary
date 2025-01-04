package main

import (
	"log/slog"

	"github.com/MaratBR/openlibrary/internal/app"
	mockeddata "github.com/MaratBR/openlibrary/internal/app/mocked_data"
	"github.com/knadh/koanf/v2"
)

func mainPopulate(config *koanf.Koanf) {
	db := connectToDatabase(config)
	sessionService := app.NewSessionService(db)
	authService := app.NewAuthService(db, sessionService)

	tagsService := app.NewTagsService(db)
	uploadService := app.NewUploadServiceFromApplicationConfig(config)
	bookManagerService := app.NewBookManagerService(db, tagsService, uploadService)

	userService := app.NewUserService(db)
	reviewsService := app.NewReviewsService(db, userService)

	slog.Info("populating database with random data...")

	setup := mockeddata.NewSetup(tagsService, reviewsService, bookManagerService, authService)
	if err := setup.Run(mockeddata.SetupOptions{
		Users:         100,
		BooksLocation: "./ao3-books",
	}); err != nil {
		panic(err)
	}

}
