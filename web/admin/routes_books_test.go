package admin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/session"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAdminBookSearchID(t *testing.T) {
	for _, query := range []string{"6704893821110461031", " 6704893821110461031 ", "/book/title-6704893821110461031", "https://example.org/book/title-6704893821110461031?admin.override=1", "/book/6704893821110461031/"} {
		id, ok := adminBookSearchID(query)
		require.True(t, ok, query)
		require.Equal(t, int64(6704893821110461031), id)
	}
	for _, query := range []string{"", "0", "-1", "9223372036854775808", "title-123", "https://example.org/users/123", "/book/123/chapters/456", "javascript:/book/123", "/book/title", "/book/-1", "https://example.org/book/title-9223372036854775808"} {
		_, ok := adminBookSearchID(query)
		require.False(t, ok, query)
	}
}

type adminBooksStub struct {
	app.BookService
	query app.GetBookQuery
	err   error
}

func (s *adminBooksStub) GetBookDetails(_ context.Context, q app.GetBookQuery) (app.BookDetailsDto, error) {
	s.query = q
	return app.BookDetailsDto{ID: q.ID, Name: "Restricted book", IsBanned: true}, s.err
}

type adminModerationBooksStub struct {
	app.ModerationBookService
	query app.SearchModerationBooksQuery
}

func (s *adminModerationBooksStub) SearchBooks(_ context.Context, q app.SearchModerationBooksQuery) (app.ModerationPage[app.ModerationBookListEntry], error) {
	s.query = q
	return app.ModerationPage[app.ModerationBookListEntry]{Page: q.Page}, nil
}
func booksTestRouter(role app.UserRole, books *adminBooksStub, moderation *adminModerationBooksStub) http.Handler {
	r := chi.NewRouter()
	log := zap.NewNop().Sugar()
	r.Use(session.Middleware(testStore{&testSession{map[string]string{}}}, log))
	r.Use(auth.NewAuthorizationMiddleware(testSessions{info: &app.SessionInfo{UserRole: app.RoleAdmin}}, &testUsers{role: role}, auth.MiddlewareOptions{}, log))
	r.Route("/admin", func(r chi.Router) {
		r.Use(requireAdmin)
		setupSPA(r, nil, nil, nil, newBooksController(books, moderation))
	})
	return r
}
func TestAdminBooksAPI(t *testing.T) {
	books := &adminBooksStub{}
	moderation := &adminModerationBooksStub{}
	h := booksTestRouter(app.RoleAdmin, books, moderation)
	out := httptest.NewRecorder()
	h.ServeHTTP(out, adminRequest("GET", "/admin/api/books?q=https%3A%2F%2Fexample.org%2Fbook%2Ftitle-6704893821110461031", "", true))
	require.Equal(t, 200, out.Code)
	require.JSONEq(t, `{"redirect":"/books/6704893821110461031"}`, out.Body.String())
	require.True(t, books.query.AdminOverride)
	require.True(t, books.query.ActorUserID.Valid)
	out = httptest.NewRecorder()
	h.ServeHTTP(out, adminRequest("GET", "/admin/api/books?q=title&p=2", "", true))
	require.Equal(t, 200, out.Code)
	require.Equal(t, "title", moderation.query.Search)
	require.True(t, moderation.query.IncludeBanned)
	require.True(t, moderation.query.IncludeDeleted)
	require.Equal(t, uint32(2), moderation.query.Page)
	require.Equal(t, uint32(20), moderation.query.PageSize)
	out = httptest.NewRecorder()
	h.ServeHTTP(out, adminRequest("GET", "/admin/api/books/6704893821110461031", "", true))
	require.Equal(t, 200, out.Code)
	require.Contains(t, out.Body.String(), `"id":"6704893821110461031"`)
	books.err = app.ErrTypeBookNotFound.New("missing")
	out = httptest.NewRecorder()
	h.ServeHTTP(out, adminRequest("GET", "/admin/api/books/1", "", true))
	require.Equal(t, 404, out.Code)
	out = httptest.NewRecorder()
	h.ServeHTTP(out, adminRequest("GET", "/admin/api/books?q=1", "", true))
	require.Equal(t, 200, out.Code)
	require.NotContains(t, out.Body.String(), "redirect")
}
func TestAdminBooksRejectNonAdmins(t *testing.T) {
	for _, role := range []app.UserRole{app.RoleModerator, app.RoleUser} {
		h := booksTestRouter(role, &adminBooksStub{}, &adminModerationBooksStub{})
		for _, path := range []string{"/admin/api/books", "/admin/api/books/1"} {
			out := httptest.NewRecorder()
			h.ServeHTTP(out, adminRequest("GET", path, "", true))
			require.Equal(t, 403, out.Code)
		}
	}
	out := httptest.NewRecorder()
	booksTestRouter(app.RoleAdmin, &adminBooksStub{}, &adminModerationBooksStub{}).ServeHTTP(out, adminRequest("GET", "/admin/api/books", "", false))
	require.Equal(t, 401, out.Code)
}
