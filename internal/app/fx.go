package app

import (
	"github.com/MaratBR/openlibrary/internal/app/analytics"
	"github.com/MaratBR/openlibrary/internal/store"
	"go.uber.org/fx"
)

var FXModule = fx.Module("ol_app", fx.Decorate(),
	analytics.FXModule,

	fx.Provide(
		NewUserService,
		NewAuthService,
		NewSignUpService,
		NewBookService,
		NewBookFullReindexService,
		fx.Annotate(
			NewBookBackgroundService,
			fx.As(fx.Self()),
			fx.As(new(BookRecalculationIngest)),
		),
		NewBookManagerService,
		NewDraftPublishingService,
		NewReviewsService,
		NewSearchService,
		NewCommentsService,
		NewSiteConfig,
		NewUploadServiceFromApplicationConfig,
		NewReadingListService,
		NewReaderPreferencesService,
		NewCollectionsService,
		NewIPLocationService,
		NewSessionService,
		NewModerationAuthorizer,
		NewModerationBookService,
		NewContentModerationRepository,
		NewContentModerationService,
		NewLoginHistoryService,
		NewModerationUserRepository,
		NewModerationUserService,
		NewReportRepository,
		NewReportService,
		NewTagsService,

		// alias for DB -> DBTX
		func(db DB) store.DBTX {
			return db
		},
	),

	cronRunnerModule,

	fx.Invoke(startTagsPopulateService),

	fx.Invoke(func(srv BookBackgroundService, draftPublishing *DraftPublishingService) {}),
)
