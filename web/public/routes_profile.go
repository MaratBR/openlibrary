package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
)

type profileController struct {
	service           app.UserService
	collectionService app.CollectionService
	bookService       app.BookService
	searchService     app.SearchService
}

func newProfileController(service app.UserService, bookService app.BookService, searchService app.SearchService, collectionService app.CollectionService) *profileController {
	return &profileController{service: service, bookService: bookService, searchService: searchService, collectionService: collectionService}
}

func (c *profileController) Register(r chi.Router) {
	r.Get("/users/{id}", c.GetProfile)
	r.Get("/users/{id}/books", c.GetBooks)
	r.Get("/users/{id}/collections", c.GetCollections)
}

func (c *profileController) GetBooks(w http.ResponseWriter, r *http.Request) {
	userID, err := olhttp.URLParamUUID(r, "id")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}

	user, err := c.service.GetUserDetails(r.Context(), app.GetUserQuery{
		ID: userID,
	})
	if err != nil {
		if err == app.ErrUserNotFound {
			olhttp.WriteTemplate(w, r.Context(), templates.UserNotFound())
			return
		}
		writeApplicationError(w, r, err)
		return
	}

	search := getBooksSearchRequest(r, &bookSearchRequestParsingOptions{
		DisableUsers: true,
	})

	searchQuery := app.BookSearchQuery{
		UserID: auth.GetNullableUserID(r.Context()),

		Query: search.Query,

		IncludeUsers:    []uuid.UUID{userID},
		ExcludeUsers:    search.ExcludeUsers,
		IncludeTags:     search.IncludeTags,
		ExcludeTags:     search.ExcludeTags,
		Words:           search.Words,
		Chapters:        search.Chapters,
		WordsPerChapter: search.WordsPerChapter,

		IncludeBanned: false,
		IncludeHidden: false,
		IncludeEmpty:  false,
		Page:          search.Page,
		PageSize:      search.PageSize,
	}
	result, err := c.searchService.SearchBooks(r.Context(), searchQuery)
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	explainedQuery, err := c.searchService.ExplainSearchQuery(r.Context(), searchQuery)
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	olhttp.WriteTemplate(w, r.Context(), templates.ProfileBooks(user, result, explainedQuery))
}

func (c *profileController) GetCollections(w http.ResponseWriter, r *http.Request) {
	userID, err := olhttp.URLParamUUID(r, "id")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}

	user, err := c.service.GetUserDetails(r.Context(), app.GetUserQuery{
		ID: userID,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	page := olhttp.GetPage(r.URL.Query(), "page")

	result, err := c.collectionService.GetUserCollections(r.Context(), app.GetUserCollectionsQuery{
		UserID:   userID,
		Page:     int32(page),
		PageSize: 20,
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

	olhttp.WriteTemplate(w, r.Context(), templates.ProfileCollections(user, result, booksMap))
}

func (c *profileController) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := olhttp.URLParamUUID(r, "id")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}

	user, err := c.service.GetUserDetails(r.Context(), app.GetUserQuery{
		ID: userID,
	})
	if err != nil {
		if err == app.ErrUserNotFound {
			olhttp.WriteTemplate(w, r.Context(), templates.UserNotFound())
			return
		}
		writeApplicationError(w, r, err)
		return
	}

	pinnedBooks, err := c.bookService.GetPinnedBooks(r.Context(), app.GetPinnedUserBooksQuery{
		UserID: user.ID,
		Limit:  6,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	templates.User(user, pinnedBooks).Render(r.Context(), w)
}
