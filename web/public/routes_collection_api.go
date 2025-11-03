package public

import (
	"fmt"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/flash"
	"github.com/MaratBR/openlibrary/internal/i18n"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
)

type apiCollectionController struct {
	collectionService app.CollectionService
}

func newAPICollectionController(collectionService app.CollectionService) *apiCollectionController {
	return &apiCollectionController{
		collectionService: collectionService,
	}
}

func (c *apiCollectionController) Register(r chi.Router) {
	r.Route("/collections", func(r chi.Router) {
		r.Use(apiRequiresAuthorizationMiddleware)
		r.Get("/recent", c.getRecent)
		r.Get("/containingBook", c.containingBook)
		r.With(httpin.NewInput(&createCollectionInput{})).Post("/", c.createCollection)
		r.With(httpin.NewInput(&addToCollectionInput{})).Post("/addBook", c.addToCollection)
		r.Delete("/removeBook/{collectionID}/{bookID}", c.removeFromCollection)
		r.Post("/removeBook/{collectionID}/{bookID}", c.removeFromCollection)
	})
}

type createCollectionInput struct {
	Body struct {
		Name string `json:"name"`
	} `in:"body=json"`
}

func (c *apiCollectionController) createCollection(w http.ResponseWriter, r *http.Request) {
	s := auth.RequireSession(r.Context())
	input := r.Context().Value(httpin.Input).(*createCollectionInput)

	collectionId, err := c.collectionService.CreateCollection(r.Context(), app.CreateCollectionCommand{
		UserID: s.UserID,
		Name:   input.Body.Name,
	})

	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}

	olhttp.NewAPIResponse(app.Int64String(collectionId)).Write(w)
}

type addToCollectionInput struct {
	Body []app.Int64String `in:"body=json"`
}

func (c *apiCollectionController) addToCollection(w http.ResponseWriter, r *http.Request) {
	s := auth.RequireSession(r.Context())
	input := r.Context().Value(httpin.Input).(*addToCollectionInput)

	bookID, err := olhttp.URLQueryParamInt64(r, "bookId")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}

	err = c.collectionService.AddToCollections(r.Context(), app.AddToCollectionsCommand{
		ActorUserID: s.UserID,
		BookID:      bookID,
		CollectionID: app.MapSlice(input.Body, func(v app.Int64String) int64 {
			return int64(v)
		}),
	})

	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}

	olhttp.NewAPIResponseOK().Write(w)
}

func (c *apiCollectionController) removeFromCollection(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	collectionID, err := olhttp.URLParamInt64(r, "collectionID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}

	session := auth.RequireSession(r.Context())

	err = c.collectionService.RemoveFromCollection(r.Context(), app.RemoveFromCollectionCommand{
		BookID:       bookID,
		CollectionID: collectionID,
		UserID:       session.UserID,
	})

	if r.Method == http.MethodPost {
		// received a form
		l := i18n.GetLocalizer(r.Context())
		flash.Add(r, flash.Text(l.T("collection.bookRemovedFromCollection")))
		http.Redirect(w, r, fmt.Sprintf("/col/%d", collectionID), http.StatusFound)
	} else {
		apiWriteOK(w)
	}
}

type recentCollectionDto struct {
	ID   int64  `json:"id,string"`
	Name string `json:"name"`
}

func collectionDtoToAPI(c app.CollectionDto) recentCollectionDto {
	return recentCollectionDto{
		ID:   c.ID,
		Name: c.Name,
	}
}

func (c *apiCollectionController) containingBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLQueryParamInt64(r, "bookId")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}

	s := auth.RequireSession(r.Context())

	collections, err := c.collectionService.GetBookCollections(r.Context(), app.GetBookCollectionsQuery{
		ActorUserID: s.UserID,
		BookID:      bookID,
	})

	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}

	olhttp.NewAPIResponse(app.MapSlice(collections, collectionDtoToAPI)).Write(w)
}

func (c *apiCollectionController) getRecent(w http.ResponseWriter, r *http.Request) {
	s := auth.RequireSession(r.Context())

	result, err := c.collectionService.GetRecentUserCollections(r.Context(), app.GetRecentCollectionsQuery{
		UserID: s.UserID,
		Limit:  20,
	})

	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}

	olhttp.NewAPIResponse(app.MapSlice(result, collectionDtoToAPI)).Write(w)
}
