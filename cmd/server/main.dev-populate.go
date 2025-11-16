package main

import (
	"log/slog"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/email"
	mockeddata "github.com/MaratBR/openlibrary/internal/app/mocked_data"
	"github.com/knadh/koanf/v2"
)

func mainPopulate(config *koanf.Koanf) {
	db := connectToDatabase(config)
	siteConfig := app.NewSiteConfig(db, config)
	sessionService := app.NewSessionService(db)
	authService := app.NewAuthService(db, sessionService)

	tagsService := app.NewTagsService(db)
	uploadService := app.NewUploadServiceFromApplicationConfig(config)
	userService := app.NewUserService(db)
	bookManagerService := app.NewBookManagerService(db, tagsService, uploadService, userService, app.NewDummyBookReindexService())
	reviewsService := app.NewReviewsService(db, userService, app.NewDummyBookBackgroundService())
	signUpService := app.NewSignUpService(db, config, siteConfig, email.NewBlackhole())

	slog.Info("populating database with random data...")

	setup := mockeddata.NewSetup(tagsService, reviewsService, bookManagerService, authService, signUpService)
	if err := setup.Run(mockeddata.SetupOptions{
		Users:         100,
		BooksLocation: "./rr-books",
	}); err != nil {
		panic(err)
	}

}
