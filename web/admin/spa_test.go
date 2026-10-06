package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/csrf"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/internal/session"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"uuid"
)

type testSession struct{ data map[string]string }

func (s *testSession) Get(k string) (string, bool) { v, ok := s.data[k]; return v, ok }
func (s *testSession) Put(k, v string)             { s.data[k] = v }
func (s *testSession) Save(context.Context) error  { return nil }
func (s *testSession) ID() string                  { return "admin-test" }

type testStore struct{ s *testSession }

func (s testStore) Get(context.Context, string) (session.Session, error) { return s.s, nil }

type testSessions struct {
	app.SessionService
	info *app.SessionInfo
}

func (s testSessions) GetBySID(context.Context, string) (*app.SessionInfo, error) { return s.info, nil }

type testUsers struct {
	app.UserService
	role    app.UserRole
	updated *app.UpdateUserCommand
	query   *app.UsersQuery
	readErr error
}

func (s *testUsers) GetUserSelfData(context.Context, uuid.UUID) (*app.SelfUserDto, error) {
	return &app.SelfUserDto{ID: uuid.MustParse("416f17a0-1740-4e61-99e0-e9e5309c8030"), Role: s.role}, nil
}
func (s *testUsers) GetUserDetails(context.Context, app.GetUserQuery) (*app.UserDetailsDto, error) {
	return &app.UserDetailsDto{}, s.readErr
}
func (s *testUsers) UpdateUser(_ context.Context, c app.UpdateUserCommand) error {
	s.updated = &c
	return nil
}
func (s *testUsers) ListUsers(_ context.Context, q app.UsersQuery) (app.UserListResponse, error) {
	s.query = &q
	return app.UserListResponse{Users: []app.UserSearchItem{}, Page: q.Page}, nil
}

type testTags struct {
	app.TagsService
	updated *app.UpdateTagCommand
	query   *app.ListTagsQuery
	missing bool
}

func (s *testTags) GetTag(_ context.Context, id int64) (app.TagDetailsItemDto, error) {
	if s.missing {
		return app.TagDetailsItemDto{}, app.ErrTagNotFound
	}
	return app.TagDetailsItemDto{DefinedTagDto: app.DefinedTagDto{ID: id}}, nil
}
func (s *testTags) UpdateTag(_ context.Context, c app.UpdateTagCommand) error {
	s.updated = &c
	return nil
}
func (s *testTags) List(_ context.Context, q app.ListTagsQuery) (app.ListTagsResult, error) {
	s.query = &q
	return app.ListTagsResult{Tags: []app.TagDetailsItemDto{}, Page: q.Page}, nil
}

func testAdminRouter(users *testUsers, tags *testTags, debug *debugController) http.Handler {
	r := chi.NewRouter()
	log := zap.NewNop().Sugar()
	r.Use(session.Middleware(testStore{&testSession{map[string]string{}}}, log))
	r.Use(auth.NewAuthorizationMiddleware(testSessions{info: &app.SessionInfo{UserRole: app.RoleAdmin}}, users, auth.MiddlewareOptions{}, log))
	r.Route("/admin", func(r chi.Router) {
		r.Use(requireAdmin)
		setupSPA(r, newUsersController(users), &tagsController{service: tags}, debug)
	})
	return r
}
func adminRequest(method, path, body string, authenticated bool) *http.Request {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if authenticated {
		r.AddCookie(&http.Cookie{Name: "sid", Value: "admin-test"})
	}
	return r
}
func TestAdminAPIRequiresCurrentAdminRole(t *testing.T) {
	for _, tc := range []struct {
		name   string
		role   app.UserRole
		auth   bool
		status int
	}{
		{"anonymous", app.RoleAdmin, false, 401}, {"user despite cached admin", app.RoleUser, true, 403}, {"moderator", app.RoleModerator, true, 403}, {"admin", app.RoleAdmin, true, 200}, {"system", app.RoleSystem, true, 200},
	} {
		t.Run(tc.name, func(t *testing.T) {
			users := &testUsers{role: tc.role}
			handler := testAdminRouter(users, &testTags{}, &debugController{})
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, adminRequest("GET", "/admin/api/users", "", tc.auth))
			require.Equal(t, tc.status, w.Code)
			require.Contains(t, w.Header().Get("Content-Type"), "application/json")
			if tc.status != 200 {
				require.Nil(t, users.query)
			}
		})
	}
}
func TestAdminLegacyRedirectRetainsQueryAndTagID(t *testing.T) {
	h := testAdminRouter(&testUsers{role: app.RoleAdmin}, &testTags{}, &debugController{})
	for _, tc := range []struct{ path, to string }{{"/admin/users?q=alice&p=2", "/admin#/users?q=alice&p=2"}, {"/admin/tags/tag-details/6704893821110461031/edit", "/admin#/tags/6704893821110461031/edit"}, {"/admin/books/", "/admin#/books"}} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, adminRequest("GET", tc.path, "", true))
		require.Equal(t, 302, w.Code)
		require.Equal(t, tc.to, w.Header().Get("Location"))
	}
}
func TestAdminUserUpdatePreservesGenderAndOptionalPassword(t *testing.T) {
	users := &testUsers{role: app.RoleAdmin}
	h := testAdminRouter(users, &testTags{}, &debugController{})
	for _, gender := range []string{"", "nonbinary", "female"} {
		w := httptest.NewRecorder()
		body := url.Values{"role": {"user"}, "about": {""}, "gender": {gender}}.Encode()
		h.ServeHTTP(w, adminRequest("POST", "/admin/api/users/416f17a0-1740-4e61-99e0-e9e5309c8030", body, true))
		require.Equal(t, 200, w.Code)
		require.Equal(t, gender, users.updated.Gender)
		require.Empty(t, users.updated.Password)
		require.True(t, users.updated.ActorUserID.Valid)
	}
	users.updated = nil
	w := httptest.NewRecorder()
	h.ServeHTTP(w, adminRequest("POST", "/admin/api/users/416f17a0-1740-4e61-99e0-e9e5309c8030", "role=invalid&gender=&about=", true))
	require.Equal(t, 400, w.Code)
	require.Nil(t, users.updated)
	users.readErr = app.ErrUserNotFound
	w = httptest.NewRecorder()
	h.ServeHTTP(w, adminRequest("POST", "/admin/api/users/416f17a0-1740-4e61-99e0-e9e5309c8030", "role=user&gender=&about=", true))
	require.Equal(t, 404, w.Code)
	require.Nil(t, users.updated)
}
func TestAdminTagUpdateAllowsUncheckedFlagsAndEmptyDescription(t *testing.T) {
	tags := &testTags{}
	h := testAdminRouter(&testUsers{role: app.RoleAdmin}, tags, &debugController{})
	w := httptest.NewRecorder()
	h.ServeHTTP(w, adminRequest("POST", "/admin/api/tags/6704893821110461031", "name=Action&type=genre&description=&synonymOf=", true))
	require.Equal(t, 200, w.Code)
	require.NotNil(t, tags.updated)
	require.False(t, tags.updated.IsAdult)
	require.False(t, tags.updated.IsSpoiler)
	require.Empty(t, tags.updated.Description)
	require.False(t, tags.updated.SynonymOfTagID.Valid)
	require.Equal(t, int64(6704893821110461031), tags.updated.ID)
	for _, body := range []string{"name=Action&type=invalid&description=", "name=Action&type=genre&description=&synonymOf=6704893821110461031", "name=Action&type=genre&description=&synonymOf=oops", "name=Action&type=genre&description=&adult=oops"} {
		tags.updated = nil
		w = httptest.NewRecorder()
		h.ServeHTTP(w, adminRequest("POST", "/admin/api/tags/6704893821110461031", body, true))
		require.Equal(t, 400, w.Code)
		require.Nil(t, tags.updated)
	}
	tags.missing = true
	w = httptest.NewRecorder()
	h.ServeHTTP(w, adminRequest("POST", "/admin/api/tags/1", "name=Action&type=genre&description=", true))
	require.Equal(t, 404, w.Code)
	require.Nil(t, tags.updated)
}
func TestAdminListsForwardFiltersAndReturnEmptyArrays(t *testing.T) {
	users := &testUsers{role: app.RoleAdmin}
	tags := &testTags{}
	h := testAdminRouter(users, tags, &debugController{})
	w := httptest.NewRecorder()
	h.ServeHTTP(w, adminRequest("GET", "/admin/api/users?q=alice&p=2&usersFilter.role=moderator", "", true))
	require.Equal(t, 200, w.Code)
	require.Equal(t, "alice", users.query.Query)
	require.Equal(t, uint32(2), users.query.Page)
	require.Equal(t, []app.UserRole{app.RoleModerator}, users.query.Role)
	require.Contains(t, w.Body.String(), `"users":[]`)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, adminRequest("GET", "/admin/api/tags?q=action&p=2&onlyAdultTags=true&onlyParentTags=true", "", true))
	require.Equal(t, 200, w.Code)
	require.Equal(t, "action", tags.query.SearchQuery)
	require.Equal(t, uint32(2), tags.query.Page)
	require.True(t, tags.query.OnlyAdultTags)
	require.True(t, tags.query.OnlyParentTags)
	require.Contains(t, w.Body.String(), `"tags":[]`)
}
func TestAdminIDsInvalidAndPrecisionSafe(t *testing.T) {
	h := testAdminRouter(&testUsers{role: app.RoleAdmin}, &testTags{}, &debugController{})
	for _, path := range []string{"/admin/api/tags/0", "/admin/api/tags/-1", "/admin/api/tags/9223372036854775808", "/admin/api/users/not-a-uuid"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, adminRequest("GET", path, "", true))
		require.Equal(t, 400, w.Code)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, adminRequest("GET", "/admin/api/tags/6704893821110461031", "", true))
	require.Equal(t, 200, w.Code)
	var tag map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &tag))
	require.Equal(t, "6704893821110461031", tag["id"])
}
func TestAdminDebugRequiresExplicitValidAction(t *testing.T) {
	calls := 0
	debug := &debugController{actions: map[string]debugActionDescriptor{"books:elastic:reindex": {Func: func() string { calls++; return "" }}}}
	h := testAdminRouter(&testUsers{role: app.RoleAdmin}, &testTags{}, debug)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, adminRequest("GET", "/admin/api/debug", "", true))
	require.Equal(t, 200, w.Code)
	require.Zero(t, calls)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, adminRequest("POST", "/admin/api/debug", "act=unknown", true))
	require.Equal(t, 400, w.Code)
	require.Zero(t, calls)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, adminRequest("POST", "/admin/api/debug", "act=books%3Aelastic%3Areindex", true))
	require.Equal(t, 200, w.Code)
	require.Equal(t, 1, calls)
}
func TestAdminMutationsRetainCSRFProtection(t *testing.T) {
	users := &testUsers{role: app.RoleAdmin}
	handler := olhttp.ReqCtxMiddleware(csrf.NewHandler("test-secret").Middleware(testAdminRouter(users, &testTags{}, &debugController{})))
	// Root middleware rejects missing tokens before any service is called.
	for _, path := range []string{"/admin/api/users/416f17a0-1740-4e61-99e0-e9e5309c8030", "/admin/api/tags/1", "/admin/api/debug"} {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, adminRequest("POST", path, "role=user&gender=&about=", true))
		require.Equal(t, 403, w.Code)
		require.Nil(t, users.updated)
	}

	token := csrf.GenerateHMACCsrfToken("admin-test", "test-secret")
	request := adminRequest("POST", "/admin/api/users/416f17a0-1740-4e61-99e0-e9e5309c8030", "role=user&gender=&about=", true)
	request.AddCookie(&http.Cookie{Name: "csrf", Value: token})
	request.Header.Set("x-csrf-token", token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request)
	require.Equal(t, 200, w.Code)
	require.NotNil(t, users.updated)
}
