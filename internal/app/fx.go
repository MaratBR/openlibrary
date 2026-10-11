package app

import (
	"github.com/MaratBR/openlibrary/internal/app/analytics"
	"github.com/MaratBR/openlibrary/internal/app/bookfont"
	"github.com/MaratBR/openlibrary/internal/app/content"
	"github.com/MaratBR/openlibrary/internal/app/reader"
	"github.com/MaratBR/openlibrary/internal/app/siteconfig"
	"github.com/MaratBR/openlibrary/internal/app/tags"
	"github.com/MaratBR/openlibrary/internal/app/upload"
	"github.com/MaratBR/openlibrary/internal/app/userdata"
	"github.com/MaratBR/openlibrary/internal/store"
	"go.uber.org/fx"
)

var FXModule = fx.Module("ol_app", fx.Decorate(),
	analytics.FXModule,
	bookfont.FXModule,
	content.FXModule,
	siteconfig.FXModule,
	tags.FXModule,
	reader.FXModule,
	upload.FXModule,
	userdata.FXModule,

	fx.Provide(
		NewUserService,
		NewAuthService,
		NewSignUpService,
		NewBookService,
		NewAdultContentService,
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
		NewReadingListService,
		NewCollectionsService,
		NewIPLocationService,
		NewSessionService,
		NewModerationAuthorizer,
		NewModerationBookService,
		NewModerationAuditLogService,
		NewContentModerationRepository,
		NewContentModerationService,
		NewLoginHistoryRepository,
		NewLoginHistoryService,
		NewModerationUserRepository,
		NewModerationUserService,
		NewReportRepository,
		NewReportService,
		NewReportActionExecutor,
		NewModerationReportService,

		// alias for DB -> DBTX
		func(db DB) store.DBTX {
			return db
		},
	),

	cronRunnerModule,

	fx.Invoke(func(srv BookBackgroundService, draftPublishing *DraftPublishingService) {}),
)
