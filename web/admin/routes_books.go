package admin

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/go-chi/chi/v5"
	"github.com/joomcode/errorx"
)

type booksController struct {
	books      app.BookService
	moderation app.ModerationBookService
}

func newBooksController(books app.BookService, moderation app.ModerationBookService) *booksController {
	return &booksController{books: books, moderation: moderation}
}

// Parse IDs without floating-point conversion, including canonical slug-ID URLs.
func adminBookSearchID(query string) (int64, bool) {
	query = strings.TrimSpace(query)
	if id, err := strconv.ParseInt(query, 10, 64); err == nil && id > 0 {
		return id, true
	}
	u, err := url.Parse(query)
	if err != nil || (u.Scheme != "" && u.Scheme != "https" && u.Scheme != "http") {
		return 0, false
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) != 2 || parts[0] != "book" || strings.HasPrefix(parts[1], "-") {
		return 0, false
	}
	id, _ := olhttp.ParseInt64Slug(parts[1])
	return id, id > 0
}

type adminBookEntry struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Author    string   `json:"author"`
	CreatedAt string   `json:"createdAt"`
	Words     int32    `json:"words"`
	Chapters  int32    `json:"chapters"`
	Reasons   []string `json:"visibilityReasons"`
}

func (c *booksController) listJSON(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	search := strings.TrimSpace(q.Get("q"))
	if id, ok := adminBookSearchID(search); ok {
		_, err := c.books.GetBookDetails(r.Context(), app.GetBookQuery{ID: id, ActorUserID: auth.GetNullableUserID(r.Context()), AdminOverride: true})
		if err == nil {
			adminJSON(w, http.StatusOK, map[string]string{"redirect": "/books/" + strconv.FormatInt(id, 10)})
			return
		}
		if !errorx.IsOfType(err, app.ErrTypeBookNotFound) {
			adminReadError(w, err)
			return
		}
		search = strconv.FormatInt(id, 10)
	}
	result, err := c.moderation.SearchBooks(r.Context(), app.SearchModerationBooksQuery{
		ActorUserID: auth.RequireUser(r.Context()).ID, Search: search, IncludeBanned: true, IncludeDeleted: true,
		Page: olhttp.GetPage(q, "p"), PageSize: 20,
	})
	if err != nil {
		adminReadError(w, err)
		return
	}
	items := make([]adminBookEntry, 0, len(result.Entries))
	for _, b := range result.Entries {
		reasons := []string{}
		for _, flag := range []struct {
			value  bool
			reason string
		}{{!b.IsPubliclyVisible, "private"}, {b.IsBanned, "banned"}, {b.IsShadowBanned, "shadowBanned"}, {b.IsTrashed, "deleted"}, {b.IsPermanentlyRemoved, "removed"}} {
			if flag.value {
				reasons = append(reasons, flag.reason)
			}
		}
		items = append(items, adminBookEntry{ID: strconv.FormatInt(b.ID, 10), Name: b.Name, Author: b.AuthorUserName, CreatedAt: b.CreatedAt.Format("2006-01-02"), Words: b.Words, Chapters: b.Chapters, Reasons: reasons})
	}
	adminJSON(w, http.StatusOK, struct {
		Books      []adminBookEntry `json:"books"`
		Page       uint32           `json:"page"`
		TotalPages uint32           `json:"totalPages"`
		Total      int64            `json:"total"`
	}{items, result.Page, result.TotalPages, result.Total})
}

func (c *booksController) bookJSON(w http.ResponseWriter, r *http.Request) {
	id, err := positiveTagID(chi.URLParam(r, "id"))
	if err != nil {
		adminAPIError(w, http.StatusBadRequest)
		return
	}
	book, err := c.books.GetBookDetails(r.Context(), app.GetBookQuery{ID: id, ActorUserID: auth.GetNullableUserID(r.Context()), AdminOverride: true})
	if err != nil {
		if errorx.IsOfType(err, app.ErrTypeBookNotFound) {
			adminAPIError(w, http.StatusNotFound)
		} else {
			adminReadError(w, err)
		}
		return
	}
	adminJSON(w, http.StatusOK, book)
}
