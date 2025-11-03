package public

import (
	"fmt"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/flash"
	"github.com/MaratBR/openlibrary/internal/i18n"
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
		r.Use(requiresAuthorizationMiddleware)
		r.Get("/", c.index)
		r.Get("/archive", c.archive)
		r.Get("/collections", c.collections)
		r.Post("/collections", c.createCollection)
		r.Get("/collections/{collectionID}/manage", c.manageCollection)
		r.Post("/collections/{collectionID}/manage", c.manageCollectionAct)
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

func (c *libraryController) manageCollectionAct(w http.ResponseWriter, r *http.Request) {
	collectionID, err := olhttp.URLParamInt64(r, "collectionID")
	session := auth.RequireSession(r.Context())
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}
	act := r.URL.Query().Get("act")

	if act == "delete" {
		c.collectionService.DeleteCollection(r.Context(), app.DeleteCollectionCommand{
			ActorUserID:  session.UserID,
			CollectionID: collectionID,
		})
		l := i18n.GetLocalizer(r.Context())
		flash.Add(r, flash.Text(l.T("collection.edit.deleted")))
		http.Redirect(w, r, "/library/collections", http.StatusFound)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/library/collections/%d/manage", collectionID), http.StatusFound)
	}
}

func (c *libraryController) manageCollection(w http.ResponseWriter, r *http.Request) {
	collectionID, err := olhttp.URLParamInt64(r, "collectionID")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}
	col, err := c.collectionService.GetCollection(r.Context(), collectionID)

	if !col.Valid {
		notFoundHandler(w, r)
		return
	}

	olhttp.WriteTemplate(
		w, r.Context(), templates.CollectionManage(col.Value),
	)
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
