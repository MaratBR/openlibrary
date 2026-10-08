package public

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/internal/session"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"uuid"
)

type preferenceUsers struct {
	app.UserService
	role  app.UserRole
	loads int
	err   error
}

func (s *preferenceUsers) GetUserSelfData(context.Context, uuid.UUID) (*app.SelfUserDto, error) {
	return &app.SelfUserDto{Role: s.role}, nil
}
func (s *preferenceUsers) GetUserModerationSettings(context.Context, uuid.UUID) (*app.UserModerationSettings, error) {
	s.loads++
	return &app.UserModerationSettings{CensoredTagsMode: "none"}, s.err
}

func TestContentPreferencesAdminScope(t *testing.T) {
	for _, tc := range []struct {
		role   app.UserRole
		path   string
		bypass bool
	}{
		{app.RoleAdmin, "/search?admin.link=1", true},
		{app.RoleSystem, "/book/1?admin.override=1", true},
		{app.RoleAdmin, "/search", false},
		{app.RoleUser, "/search?admin.link=1", false},
		{app.RoleModerator, "/search?admin.override=1", false},
		{app.RoleAdmin, "/admin/api/books", true},
		{app.RoleUser, "/books-manager/", true},
		{app.RoleUser, "/_api/books-manager/tags", true},
		{app.RoleModerator, "/_api/moderation/books", true},
		{app.RoleUser, "/account/settings/moderation", true},
		{app.RoleUser, "/administrator", false},
	} {
		t.Run(string(tc.role)+tc.path, func(t *testing.T) {
			users := &preferenceUsers{role: tc.role}
			log := zap.NewNop().Sugar()
			router := chi.NewRouter()
			router.Use(olhttp.ReqCtxMiddleware, session.Middleware(overridePromptStore{}, log), auth.NewAuthorizationMiddleware(overridePromptSessions{}, users, auth.MiddlewareOptions{}, log))
			// No DB is needed for empty preferences. Calling it would panic.
			router.Use(contentPreferencesMiddleware(users, nil))
			router.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, tc.bypass, bypassContentPreferences(r))
				w.WriteHeader(204)
			}))
			r := httptest.NewRequest("GET", tc.path, nil)
			r.AddCookie(&http.Cookie{Name: session.CookieName, Value: "fixture"})
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
			require.Equal(t, 204, w.Code)
			if tc.bypass {
				require.Zero(t, users.loads)
			} else {
				require.Equal(t, 1, users.loads)
				require.Equal(t, "private, no-store", w.Header().Get("Cache-Control"))
			}
		})
	}
}

func TestPreferenceLoadFailureDoesNotFallThrough(t *testing.T) {
	users := &preferenceUsers{role: app.RoleUser, err: errors.New("preferences unavailable")}
	log := zap.NewNop().Sugar()
	router := chi.NewRouter()
	router.Use(olhttp.ReqCtxMiddleware, session.Middleware(overridePromptStore{}, log), auth.NewAuthorizationMiddleware(overridePromptSessions{}, users, auth.MiddlewareOptions{}, log), contentPreferencesMiddleware(users, nil))
	called := false
	router.Get("/search", func(http.ResponseWriter, *http.Request) { called = true })
	r := httptest.NewRequest("GET", "/search", nil)
	r.AddCookie(&http.Cookie{Name: session.CookieName, Value: "fixture"})
	r.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	require.Equal(t, 500, w.Code)
	require.False(t, called)
}

type preferenceTags struct{ app.TagsService }

func (preferenceTags) SearchTags(context.Context, string) ([]app.DefinedTagDto, error) {
	return []app.DefinedTagDto{{ID: 1, Name: "Hidden"}, {ID: 2, Name: "Visible"}}, nil
}

func TestTagSelectorExcludesBannedTags(t *testing.T) {
	for _, mode := range []app.CensorMode{"none", "hide", "censor"} {
		r := httptest.NewRequest("GET", "/_api/tags?q=hi", nil)
		r = r.WithContext(app.WithContentPreferences(r.Context(), &app.ContentPreferences{Mode: mode, TagIDs: []int64{1}}))
		w := httptest.NewRecorder()
		(&apiControllerTags{service: preferenceTags{}}).Tags(w, r)
		require.Equal(t, 200, w.Code)
		if mode == "none" {
			require.Contains(t, w.Body.String(), "Hidden")
		} else {
			require.NotContains(t, w.Body.String(), "Hidden")
		}
		require.Contains(t, w.Body.String(), "Visible")
	}
}
