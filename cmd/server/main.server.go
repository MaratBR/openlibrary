package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/cache"
	"github.com/MaratBR/openlibrary/internal/csrf"
	elasticstore "github.com/MaratBR/openlibrary/internal/elastic-store"
	i18n "github.com/MaratBR/openlibrary/internal/i18n"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/internal/reqid"
	"github.com/MaratBR/openlibrary/internal/session"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/MaratBR/openlibrary/web/admin"
	"github.com/MaratBR/openlibrary/web/frontend"
	"github.com/MaratBR/openlibrary/web/public"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knadh/koanf/v2"
	"github.com/redis/go-redis/v9"
	"golang.org/x/text/language"
)

type cliParams struct {
	Dev            bool
	BypassTLSCheck bool
	StaticDir      string
}

func mainServer(
	cliParams cliParams,
	config *koanf.Koanf,
) {

	if cliParams.Dev {
		app.GlobalFeatureFlags.DisableCache = true
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	var err error

	// --------------------------------------
	// initialize server dependencies
	// --------------------------------------

	slog.Debug("initializing localizer provider")
	localizerProvider := i18n.NewLocaleProvider(
		language.English,
		cliParams.Dev,
		map[language.Tag][]string{
			language.English: {
				"translations/en.toml",
			},
		},
	)

	slog.Debug("connecting to database")
	db := connectToDatabase(config)

	slog.Debug("connecting to cache")
	cacheInstance := createCache(config)

	slog.Debug("connecting to elasticsearch client")
	esClient := setupElasticsearch(config)

	slog.Debug("initializing csrf handler")
	csrfHandler := csrf.NewHandler("CSRF HANDLER HERE")

	slog.Debug("initializing redis and session store")
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.String("redis.addr"),
	})
	sessionStore := session.NewRedisStore(redisClient)

	// --------------------------------------
	// initialize web server
	// --------------------------------------

	r := chi.NewRouter()
	r.Use(olhttp.ReqCtxMiddleware)
	r.Use(reqid.New())
	r.Use(olhttp.MakeRecoveryMiddleware())
	r.Use(csrfHandler.Middleware)
	r.Use(localizerProvider.Middleware)

	// mount assets
	assetsFS := frontend.AssetsFS(frontend.AssetsConfig{Dev: cliParams.Dev})
	assetsHandler := frontend.Assets(assetsFS)

	embedAssetsFS := frontend.EmbedAssetsFS()
	embedAssetsHandler := frontend.EmbedAssets(embedAssetsFS)

	frontend.AttachAssetsInliningHandler(embedAssetsFS, "embed-assets", r)
	frontend.AttachAssetsInliningHandler(assetsFS, "assets", r)

	r.Mount("/_/embed-assets/", http.StripPrefix("/_/embed-assets/", embedAssetsHandler))
	r.Mount("/_/assets/cats", http.StripPrefix("/_/assets/cats", http.FileServer(http.Dir("./cats"))))
	r.Mount("/_/assets/", http.StripPrefix("/_/assets/", assetsHandler))

	// --------------------------------------
	// create and start background services
	// --------------------------------------
	bgServices := app.NewBackgroundServices(db, esClient)
	err = bgServices.Start()
	if err != nil {
		panic("failed to start background services: " + err.Error())
	}
	defer bgServices.Stop()

	uploadService := app.NewUploadServiceFromApplicationConfig(config)

	publicUIHandler := public.NewHandler(db, config, cacheInstance, csrfHandler, bgServices, uploadService, esClient)
	adminHandler := admin.NewHandler(db, config, cacheInstance, bgServices, esClient)

	err = publicUIHandler.Start()
	if err != nil {
		panic("failed to start public ui handler: " + err.Error())
	}

	err = adminHandler.Start()
	if err != nil {
		panic("failed to start admin handler: " + err.Error())
	}

	// the area where all the fun happens
	r.Group(func(r chi.Router) {
		r.Use(session.Middleware(sessionStore))
		r.Mount("/", publicUIHandler)
		r.Mount("/admin", adminHandler)

	})

	defer publicUIHandler.Stop()
	defer adminHandler.Stop()

	//
	// start background services
	//

	//
	// post-initialization stuff
	//
	if config.Bool("init.create-default-users") {
		go func() {
			authService := app.NewAuthService(db, app.NewSessionService(db))
			err := authService.EnsureAdminUserExists(context.Background())
			if err != nil {
				slog.Error("failed to ensure admin user exists", "err", err)
			}
		}()
	}

	if config.Bool("init.import-predefined-tags") {
		go func() {
			err := app.ImportPredefinedTags(context.Background(), store.New(db))
			if err != nil {
				slog.Error("failed to import predefined tags", "err", err)
			}
		}()
	}

	go func() {
		err := uploadService.InitBuckets(context.Background())
		if err != nil {
			slog.Error("failed to make main bucket", "err", err)
		}
	}()

	go func() {
		ip, err := getPublicIP()
		if err == nil {
			slog.Info("public ip", "ip", ip)
		} else {
			slog.Error("failed to get public ip", "err", err)
		}
	}()

	listenOn := fmt.Sprintf("%s:%d", config.String("server.host"), config.Int("server.port"))
	slog.Info("server listening", "on", listenOn, "url", fmt.Sprintf("http://%s", listenOn))

	srv := &http.Server{
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
		Addr:           listenOn,
		Handler:        r,
	}
	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func connectToDatabase(config *koanf.Koanf) app.DB {
	connectionString := config.String("database.url")
	if connectionString == "" {
		slog.Error("database.url is empty")
		os.Exit(1)
	}
	var (
		db       *pgxpool.Pool
		err      error
		sleeping = time.Second
	)

	for {
		db, err = store.Connect(context.Background(), connectionString)
		if err != nil {
			slog.Error("failed to connect to DB", "err", err, "sleeping", sleeping)
			time.Sleep(sleeping)
			if sleeping < time.Second*60 {
				sleeping *= 2
			}
		} else {
			break
		}
	}

	return db
}

func createCache(config *koanf.Koanf) *cache.Cache {
	cacheBackend, err := cache.CacheBackendFromConfig(config)
	if err != nil {
		panic(err)
	}
	cacheInstance := cache.New(cacheBackend)
	return cacheInstance
}

func setupElasticsearch(config *koanf.Koanf) *elasticsearch.TypedClient {
	elasticsearchURL := config.String("elasticsearch.url")
	if elasticsearchURL == "" {
		slog.Error("elasticsearch.url is empty")
		os.Exit(1)
	}

	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{elasticsearchURL},
	})
	if err != nil {
		panic(err)
	}

	go func() {
		err = elasticstore.Setup(context.Background(), client)
		if err != nil {
			slog.Error("FAILED TO SETUP ELASTIC", "err", err)
		}
	}()

	return client
}
