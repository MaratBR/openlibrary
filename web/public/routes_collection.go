package public

import (
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
	collectionID, err := olhttp.URLParamInt64(r, "collectionID")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}

	page := olhttp.GetPage(r.URL.Query(), "p")
	pageSize := olhttp.GetPageSize(r.URL.Query(), "pageSize", 1, 100, 20)

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
	canEdit := isAuthorized && result.Collection.UserID == session.UserID

	olhttp.WriteTemplate(w, r.Context(), templates.Collection(result, canEdit))
}
