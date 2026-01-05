package main

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/csrf"
	"github.com/MaratBR/openlibrary/internal/flash"
	i18n "github.com/MaratBR/openlibrary/internal/i18n"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/internal/reqid"
	"github.com/MaratBR/openlibrary/internal/session"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/MaratBR/openlibrary/web/admin"
	"github.com/MaratBR/openlibrary/web/frontend"
	"github.com/MaratBR/openlibrary/web/public"
	"github.com/MaratBR/openlibrary/web/webfx"
	"github.com/go-chi/chi/v5"
	"github.com/knadh/koanf/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

//go:embed default_robots.txt
var robotsTxt string

type cliParams struct {
	Dev            bool
	BypassTLSCheck bool
	StaticDir      string
}

func mainServer(
	cliParams cliParams,
) {
	if cliParams.Dev {
		app.GlobalFeatureFlags.DisableCache = true
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	fx.New(
		fx.Supply(cliParams),
		infraModule,
		app.FXModule,

		public.FXModule,
		admin.FXModule,
		flash.FXModule,

		fx.Provide(
			newHTTPServer,
			zap.NewExample,
			func(log *zap.Logger) *zap.SugaredLogger {
				return log.Sugar()
			},
			fx.Annotate(
				newRootMux,
				webfx.ProviderMountables(),
			),
		),
		fx.Invoke(postInit, func(*http.Server) {}),
	).Run()
}

func newRootMux(
	// handlers MUST go first
	handlers []webfx.MountableHandler,
	cliParams cliParams,
	config *koanf.Koanf,
	csrfHandler *csrf.Handler,
	sessionStore session.Store,
	localeProvider *i18n.LocaleProvider,

	log *zap.SugaredLogger,
) http.Handler {

	r := chi.NewRouter()

	r.Use(olhttp.ReqCtxMiddleware)
	r.Use(reqid.New())
	r.Use(olhttp.MakeRecoveryMiddleware())
	r.Use(csrfHandler.Middleware)
	r.Use(localeProvider.Middleware)

	// mount assets
	assetsFS := frontend.AssetsFS(frontend.AssetsConfig{Dev: cliParams.Dev})
	assetsHandler := frontend.Assets(assetsFS)

	embedAssetsFS := frontend.EmbedAssetsFS()
	embedAssetsHandler := frontend.EmbedAssets(embedAssetsFS)

	frontend.AttachAssetsInliningHandler(embedAssetsFS, "embed-assets", r)
	frontend.AttachAssetsInliningHandler(assetsFS, "assets", r)

	r.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		s := strings.ReplaceAll(robotsTxt, "{{sitemap}}", fmt.Sprintf("https://%s/sitemap.xml", r.Host))
		w.Write([]byte(s))
	})

	r.Mount("/_/embed-assets/", http.StripPrefix("/_/embed-assets/", embedAssetsHandler))
	r.Mount("/_/assets/cats", http.StripPrefix("/_/assets/cats", http.FileServer(http.Dir("./cats"))))
	r.Mount("/_/assets/", http.StripPrefix("/_/assets/", assetsHandler))

	// the area where all the fun happens
	r.Group(func(r chi.Router) {
		r.Use(session.Middleware(sessionStore))

		for _, h := range handlers {
			log.Infow("registering mounted handler", "at", h.MountAt(), "type", reflect.TypeOf(h).String())
			r.Mount(h.MountAt(), h)
		}
	})

	return r
}

func postInit(cfg *koanf.Koanf, uploadService *app.UploadService, authService app.AuthService, db app.DB, log *zap.SugaredLogger) {
	time.Sleep(time.Second * 2)

	if cfg.Bool("init.create-default-users") {
		go func() {
			err := authService.EnsureAdminUserExists(context.Background())
			if err != nil {
				log.Errorw("failed to ensure admin user exists", "err", err)
			}
		}()
	}

	if cfg.Bool("init.import-predefined-tags") {
		go func() {
			err := app.ImportPredefinedTags(context.Background(), store.New(db))
			if err != nil {
				log.Errorw("failed to import predefined tags", "err", err)
			}
		}()
	}

	go func() {
		err := uploadService.InitBuckets(context.Background())
		if err != nil {
			log.Errorw("failed to make main bucket", "err", err)
		}
	}()

	go func() {
		ip, err := getPublicIP()
		if err == nil {
			log.Infow("public ip", "ip", ip)
		} else {
			log.Errorw("failed to get public ip", "err", err)
		}
	}()

}

func newHTTPServer(lc fx.Lifecycle, handler http.Handler, cfg *koanf.Koanf, log *zap.SugaredLogger) *http.Server {
	addr := fmt.Sprintf("%s:%d", cfg.String("server.host"), cfg.Int("server.port"))

	srv := &http.Server{
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
		Addr:           addr,
		Handler:        handler,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Infow("starting server", "addr", addr, "url", fmt.Sprintf("http://%s", addr))
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("server stopping...")
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
