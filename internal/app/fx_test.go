package app_test

import (
	"testing"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/dal"
	"github.com/MaratBR/openlibrary/internal/app/email"
	"github.com/MaratBR/openlibrary/internal/app/reader"
	"github.com/MaratBR/openlibrary/internal/app/tags"
	"github.com/MaratBR/openlibrary/internal/app/upload"
	"github.com/MaratBR/openlibrary/internal/app/userdata"
	"github.com/knadh/koanf/v2"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func TestApplicationFeatureWiring(t *testing.T) {
	// Validate the complete graph without opening connections or starting workers.
	err := fx.ValidateApp(
		fx.NopLogger,
		app.FXModule,
		fx.Provide(
			func() app.DB { return nil },
			func() dal.DB { return nil },
			func() *koanf.Koanf { return nil },
			func() *zap.SugaredLogger { return nil },
			func() *redis.Client { return nil },
			func() *opensearchapi.Client { return nil },
			func() email.Service { return nil },
		),
		fx.Invoke(func(
			reader.ReaderPreferencesService,
			reader.ReaderFontService,
			tags.TagsService,
			*upload.UploadService,
			userdata.UserDataService,
			app.BookManagerService,
			app.SearchService,
			app.SignUpService,
			app.SetupService,
		) {
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
}
