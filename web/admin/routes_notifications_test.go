package admin

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/notifications"
	"github.com/MaratBR/openlibrary/internal/csrf"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/stretchr/testify/require"
)

type recordingNotificationService struct {
	notifications.Service
	user, role, title, content string
	options                    notifications.Options
	err                        error
	called                     bool
}

func (s *recordingNotificationService) Send(_ context.Context, user, role, title, content string, options ...notifications.Options) (int64, error) {
	s.called = true
	s.user, s.role, s.title, s.content = user, role, title, content
	s.options = options[0]
	return 2, s.err
}

func TestAdminNotificationPresentationAndCSRF(t *testing.T) {
	service := &recordingNotificationService{}
	handler := olhttp.ReqCtxMiddleware(csrf.NewHandler("test-secret").Middleware(testAdminRouter(&testUsers{role: app.RoleAdmin}, &testTags{}, &debugController{notifications: service})))
	form := url.Values{"role": {"moderator"}, "title": {"Review"}, "content": {"<b>Important</b>"}, "type": {"suggestions"}, "category": {"moderation"}, "urgent": {"true"}, "iconType": {"emoji"}, "icon": {"📚"}}
	request := adminRequest("POST", "/admin/api/debug/notifications", form.Encode(), true)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request)
	require.Equal(t, http.StatusForbidden, w.Code)
	require.False(t, service.called)
	token := csrf.GenerateHMACCsrfToken("admin-test", "test-secret")
	request = adminRequest("POST", "/admin/api/debug/notifications", form.Encode(), true)
	request.AddCookie(&http.Cookie{Name: "csrf", Value: token})
	request.Header.Set("x-csrf-token", token)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, request)
	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"count":2}`, w.Body.String())
	require.Equal(t, "moderator", service.role)
	require.Equal(t, "<b>Important</b>", service.content)
	require.Equal(t, "suggestions", service.options.Type)
	require.Equal(t, notifications.CategoryModeration, service.options.Category)
	require.True(t, service.options.Urgent)
	require.Equal(t, &notifications.Icon{Type: notifications.IconEmoji, Payload: "📚"}, service.options.Icon)
	for _, tc := range []struct {
		err    error
		status int
	}{
		{notifications.ErrInvalidSend, http.StatusBadRequest}, {errors.New("database failed"), http.StatusInternalServerError},
	} {
		service.err = tc.err
		request = adminRequest("POST", "/admin/api/debug/notifications", form.Encode(), true)
		request.AddCookie(&http.Cookie{Name: "csrf", Value: token})
		request.Header.Set("x-csrf-token", token)
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, request)
		require.Equal(t, tc.status, w.Code)
	}
}
