package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"strings"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	cliParams struct {
		DevProxy bool
	}
)

func main() {
	flag.BoolVar(&cliParams.DevProxy, "dev-frontend-proxy", false, "enable dev frontend proxy")
	flag.Parse()

	db := connectToDatabase()
	authService := app.NewAuthService(db)
	authorizationMiddleware := newAuthorizationMiddleware(authService)
	requiresAuthorization := newRequireAuthorizationMiddleware()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(csrfMiddleware)

	if cliParams.DevProxy {
		fallbackHandler := newAuthorizationMiddlewareConditional(authService, func(r *http.Request) bool {
			if strings.HasPrefix(r.URL.Path, "/node_modules/") || strings.HasPrefix(r.URL.Path, "/src/") || strings.HasPrefix(r.URL.Path, "/@vite/") || r.URL.Path == "/vite.svg" {
				return false
			}
			return true
		})(http.HandlerFunc(devProxyIndex))

		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			fallbackHandler.ServeHTTP(w, r)
		})
	}

	r.Route("/api", func(r chi.Router) {
		authController := newAuthController(authService)

		r.Use(authorizationMiddleware)
		r.Post("/auth/signin", authController.SignIn)

		if cliParams.DevProxy {
			r.HandleFunc("/auth/signin-admin", authController.SignInAdmin)
		}

		r.Handle("/auth/csrf", http.HandlerFunc(refreshCsrfToken))

		tagsService := app.NewTagsService(db)
		bookService := app.NewBookService(db, tagsService)
		bookManagerService := app.NewBookManagerService(db, tagsService)
		bookController := newBookController(bookService)

		r.Get("/books/{id}", bookController.GetBook)
		r.Get("/books/{bookID}/chapters/{chapterID}", bookController.GetChapter)

		r.Route("/tags", func(r chi.Router) {
			tagsController := newTagsController(tagsService)
			r.Get("/search", tagsController.Search)
		})

		r.Route("/manager", func(r chi.Router) {
			r.Use(requiresAuthorization)

			bookManagerController := newBookManagerController(bookManagerService)

			r.Post("/books", bookManagerController.CreateBook)
			r.Get("/books/my-books", bookManagerController.GetMyBooks)
			r.Get("/books/{bookID}", bookManagerController.GetBook)
			r.Post("/books/{bookID}", bookManagerController.UpdateBook)
			r.Post("/books/{bookID}/chapters", bookManagerController.CreateChapter)
			r.Post("/books/{bookID}/chapters/{chapterID}", bookManagerController.UpdateChapter)
		})
	})

	go func() {
		err := authService.EnsureAdminUserExists(context.Background())
		if err != nil {
			slog.Error("failed to ensure admin user exists", "err", err)
		}
	}()

	go func() {
		err := app.ImportPredefinedTags(context.Background(), store.New(db))
		if err != nil {
			slog.Error("failed to import predefined tags", "err", err)
		}
	}()

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}

func connectToDatabase() *pgxpool.Pool {
	db, err := store.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/openlibrary?sslmode=disable")
	if err != nil {
		panic(err)
	}
	return db
}
