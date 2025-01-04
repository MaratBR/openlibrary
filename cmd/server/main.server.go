package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/MaratBR/openlibrary/cmd/server/csrf"
	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/cache"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/knadh/koanf/v2"
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

	db := connectToDatabase(config)

	// infrastructure layer services
	cacheInstance := createCache(config)

	uploadService := app.NewUploadServiceFromApplicationConfig(config)

	// application layer services
	sessionService := app.NewSessionService(db)
	sessionService = app.NewCachedSessionService(sessionService, cacheInstance)
	authService := app.NewAuthService(db, sessionService)

	favoriteRecalculationBackgroundService := app.NewFavoriteRecalculationBackgroundService(db)
	favoritesService := app.NewFavoriteService(db, favoriteRecalculationBackgroundService)

	tagsService := app.NewTagsService(db)
	readingListService := app.NewReadingListService(db)
	bookService := app.NewBookService(db, tagsService, uploadService, readingListService)
	bookManagerService := app.NewBookManagerService(db, tagsService, uploadService)
	searchService := app.NewSearchService(db, tagsService, uploadService)
	searchService = app.NewCachedSearchService(searchService, cacheInstance)
	userService := app.NewUserService(db)
	reviewsService := app.NewReviewsService(db, userService)
	// reviewsService = app.NewCachedReviewsService(reviewsService, cacheInstance)

	// middlewares
	authorizationMiddleware := newAuthorizationMiddleware(sessionService)
	requiresAuthorization := newRequireAuthorizationMiddleware()
	csrfHandler := csrf.NewHandler("csrf secret here, REPLACE LATER")

	// controllers
	bookController := newBookController(bookService)
	authController := newAuthController(authService, cliParams.BypassTLSCheck, csrfHandler)
	searchController := newSearchController(searchService, tagsService)
	userController := newUserController(userService)
	favoritesController := newFavoritesController(favoritesService)
	bookManagerController := newBookManagerController(bookManagerService)
	settingsController := newSettingsController(userService, sessionService)
	readingListController := newReadingListController(readingListService)
	reviewsController := newReviewsController(reviewsService)

	// create router
	r := chi.NewRouter()
	r.Use(csrfHandler.Middleware)
	r.Use(authorizationMiddleware)

	//
	// init spa and data preload
	//
	spaHandler := newSPAHandler(config, bookService, reviewsService, userService, searchService, tagsService)
	r.NotFound(spaHandler.ServeHTTP)
	{
		// initialize front-end data preload
		// searchUIController := newSearchUIController(searchService, tagsService)
		// r.Get("/search", searchUIController.Search)
	}

	//
	// init api endpoints
	//
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.Logger)

		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
			w.Write([]byte("not found"))
		})

		//
		// public endpoints first
		//
		r.Post("/auth/signin", authController.SignIn)
		r.Post("/auth/signup", authController.SignUp)
		r.Post("/auth/signout", authController.SignOut)
		if cliParams.Dev {
			r.HandleFunc("/auth/signin-admin", authController.SignInAdmin)
		}
		r.Handle("/auth/csrf-check", http.HandlerFunc(csrfHandler.CheckEndpoint))

		r.Get("/users/{userID}", userController.GetUser)
		r.Get("/users/whoami", userController.Whoami)

		r.Get("/books/{id}", bookController.GetBook)
		r.Get("/books/{bookID}/chapters/{chapterID}", bookController.GetChapter)

		r.Get("/search", searchController.Search)
		r.Get("/search/book-extremes", searchController.GetBookExtremes)

		r.Route("/tags", func(r chi.Router) {
			tagsController := newTagsController(tagsService)
			r.Get("/search-tags", tagsController.Search)
			r.Get("/lookup", tagsController.GetByName)
		})

		r.Route("/reviews", func(r chi.Router) {
			r.Get("/{bookID}", reviewsController.GetReviews)
			r.Post("/{bookID}", reviewsController.UpdateOrCreateReview)

		})

		//
		// endpoints requiring authorization grouped
		//
		r.Group(func(r chi.Router) {
			r.Use(requiresAuthorization)

			// reading list
			r.Route("/reading-list", func(r chi.Router) {
				r.Get("/status", readingListController.GetStatus)
				r.Delete("/delete", readingListController.Delete)
				r.Post("/status", readingListController.UpdateStatus)
				r.Post("/start-reading", readingListController.StartReading)

			})

			// users
			r.Post("/users/follow", userController.Follow)
			r.Post("/users/unfollow", userController.Unfollow)

			// favorites
			r.Post("/favorite", favoritesController.SetFavorite)

			// user settings
			r.Get("/settings/about", settingsController.GetAboutSettings)
			r.Put("/settings/about", settingsController.UpdateAboutSettings)
			r.Get("/settings/privacy", settingsController.GetPrivacySettings)
			r.Put("/settings/privacy", settingsController.UpdatePrivacySettings)
			r.Get("/settings/moderation", settingsController.GetModerationSettings)
			r.Put("/settings/moderation", settingsController.UpdateModerationSettings)
			r.Get("/settings/customization", settingsController.GetCustomizationSettings)
			r.Put("/settings/customization", settingsController.UpdateCustomizationSettings)
			r.Get("/settings/sessions", settingsController.GetSessions)

			// book manager
			r.Route("/manager", func(r chi.Router) {
				r.Post("/books", bookManagerController.CreateBook)
				r.Post("/books/ao3-import", bookManagerController.ImportAO3)
				r.Get("/books/my-books", bookManagerController.GetMyBooks)
				r.Get("/books/{bookID}", bookManagerController.GetBook)
				r.Post("/books/{bookID}", bookManagerController.UpdateBook)
				r.Post("/books/{bookID}/cover", bookManagerController.UploadBookCover)
				r.Get("/books/{bookID}/chapters", bookManagerController.GetChapters)
				r.Post("/books/{bookID}/chapters", bookManagerController.CreateChapter)
				r.Post("/books/{bookID}/chapters/reorder", bookManagerController.UpdateChaptersOrder)
				r.Post("/books/{bookID}/chapters/{chapterID}", bookManagerController.UpdateChapter)
				r.Get("/books/{bookID}/chapters/{chapterID}", bookManagerController.GetChapter)
			})
		})

	})

	//
	// post-initialization stuff
	//
	if config.Bool("init.create-default-users") {
		go func() {
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

	favoriteRecalculationBackgroundService.Start()

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
	db, err := store.Connect(context.Background(), connectionString)
	if err != nil {
		panic(err)
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
