package public

import (
	"fmt"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/go-chi/chi/v5"
)

type libraryController struct {
	service           app.ReadingListService
	collectionService app.CollectionService
}

func newLibraryController(service app.ReadingListService, collectionService app.CollectionService) *libraryController {
	return &libraryController{service: service, collectionService: collectionService}
}

func (c *libraryController) Register(r chi.Router) {
	r.Route("/library", func(r chi.Router) {
		r.Use(redirectToLoginOnUnauthorized)
		r.Get("/", c.index)
		r.Get("/archive", c.archive)
		r.Get("/collections", c.collections)
		r.Post("/collections", c.createCollection)

	})
}

func (c *libraryController) index(w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())

	wantToRead, err := c.service.GetReadingListBooks(r.Context(), app.GetReadingListItemsQuery{
		UserID: session.UserID,
		Limit:  12,
		Status: app.ReadingListStatusWantToRead,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	reading, err := c.service.GetReadingListBooks(r.Context(), app.GetReadingListItemsQuery{
		UserID: session.UserID,
		Limit:  12,
		Status: app.ReadingListStatusReading,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	paused, err := c.service.GetReadingListBooks(r.Context(), app.GetReadingListItemsQuery{
		UserID: session.UserID,
		Limit:  12,
		Status: app.ReadingListStatusPaused,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	templates.Library(wantToRead, reading, paused).Render(r.Context(), w)
}

func (c *libraryController) archive(w http.ResponseWriter, r *http.Request) {
	session, isAuthorized := auth.GetSession(r.Context())

	if !isAuthorized {
		http.Redirect(w, r, "/login?next=/library", http.StatusSeeOther)
		return
	}

	read, err := c.service.GetReadingListBooks(r.Context(), app.GetReadingListItemsQuery{
		UserID: session.UserID,
		Limit:  12,
		Status: app.ReadingListStatusRead,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	dnf, err := c.service.GetReadingListBooks(r.Context(), app.GetReadingListItemsQuery{
		UserID: session.UserID,
		Limit:  12,
		Status: app.ReadingListStatusDnf,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	templates.LibraryArchive(read, dnf).Render(r.Context(), w)
}

func (c *libraryController) collections(w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())
	page := olhttp.GetPage(r.URL.Query(), "p")
	pageSize := olhttp.GetPageSize(r.URL.Query(), "pageSize", 1, 100, 15)
	result, err := c.collectionService.GetUserCollections(r.Context(), app.GetUserCollectionsQuery{
		UserID:   session.UserID,
		Page:     int32(page),
		PageSize: int32(pageSize),
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	booksMap, err := c.collectionService.GetCollectionBooksMap(r.Context(), result.Collections)
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	olhttp.WriteTemplate(w, r.Context(), templates.LibraryCollections(result, booksMap))
}

func (c *libraryController) createCollection(w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())
	if err := r.ParseForm(); err != nil {
		writeApplicationError(w, r, err)
		return
	}

	name := r.Form.Get("name")
	collectionID, err := c.collectionService.CreateCollection(r.Context(), app.CreateCollectionCommand{
		UserID: session.UserID,
		Name:   name,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/col/%d", collectionID), http.StatusSeeOther)
}
