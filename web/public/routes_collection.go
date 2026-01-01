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

type collectionController struct {
	collectionService app.CollectionService
}

func newCollectionController(collectionService app.CollectionService) *collectionController {
	return &collectionController{
		collectionService: collectionService,
	}
}

func (c *collectionController) Register(r chi.Router) {
	r.Route("/col", func(r chi.Router) {
		r.Get("/{collectionID}", c.books)
	})
}

func (c *collectionController) books(w http.ResponseWriter, r *http.Request) {
	collectionIDRaw := chi.URLParam(r, "collectionID")
	collectionID, slug := olhttp.ParseInt64Slug(collectionIDRaw)
	if collectionID == 0 {
		writeBadRequest(w, r, nil)
	}
	page := olhttp.GetPage(r.URL.Query(), "p")
	pageSize := olhttp.GetPageSize(r.URL.Query(), "pageSize", 1, 100, 20)

	collectionOpt, err := c.collectionService.GetCollection(r.Context(), collectionID)
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}
	if !collectionOpt.Valid {
		// no such collection
		// TODO!
		notFoundHandler(w, r)
		return
	}
	collection := collectionOpt.Value

	if slug != collection.Slug {
		http.Redirect(w, r, fmt.Sprintf("/col/%s-%d", collection.Slug, collection.ID), http.StatusFound)
		return
	}

	result, err := c.collectionService.GetCollectionBooks(r.Context(), app.GetCollectionBooksQuery{
		CollectionID: collectionID,
		Page:         int32(page),
		PageSize:     int32(pageSize),
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	session, isAuthorized := auth.GetSession(r.Context())
	canEdit := isAuthorized && collection.UserID == session.UserID

	olhttp.WriteTemplate(w, r.Context(), templates.Collection(result, collection, canEdit))
}
