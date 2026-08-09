package populate

import (
	"fmt"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/analytics"
	"github.com/MaratBR/openlibrary/internal/app/dal"
	"github.com/MaratBR/openlibrary/internal/app/email"
	mockeddata "github.com/MaratBR/openlibrary/internal/app/mocked_data"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

func Run(config *koanf.Koanf, db dal.DB, log *zap.SugaredLogger) error {
	siteConfig := app.NewSiteConfig(db, config)
	sessionService := app.NewSessionService(db, app.NewIPLocationService())
	authService := app.NewAuthService(db, sessionService)

	tagsService := app.NewTagsService(db)
	uploadService := app.NewUploadServiceFromApplicationConfig(config)
	userService := app.NewUserService(db)
	bookManagerService := app.NewBookManagerService(db, tagsService, uploadService, userService, app.NewDummyBookReindexService(), analytics.NewDummyMetricService(), log)
	reviewsService := app.NewReviewsService(db, userService, app.NewDummyBookBackgroundService())
	signUpService := app.NewSignUpService(db, config, siteConfig, email.NewBlackhole())

	log.Infow("populating database with random data")
	setup := mockeddata.NewSetup(tagsService, reviewsService, bookManagerService, authService, signUpService)
	if err := setup.Run(mockeddata.SetupOptions{
		Users:         100,
		BooksLocation: "./rr-books",
	}); err != nil {
		return fmt.Errorf("populate database: %w", err)
	}
	return nil
}
