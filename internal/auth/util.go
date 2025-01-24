package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/gofrs/uuid"
)

func GetSession(ctx context.Context) (*app.SessionInfo, bool) {
	value := ctx.Value(sessionInfoKey)
	if value == nil {
		return nil, false
	}

	sessionInfo := value.(*app.SessionInfo)
	return sessionInfo, true
}

func GetUser(ctx context.Context) (*app.SelfUserDto, bool) {
	value := ctx.Value(userKey)
	if value == nil {
		return nil, false
	}

	sessionInfo := value.(*app.SelfUserDto)
	return sessionInfo, true
}

func RequireUser(ctx context.Context) *app.SelfUserDto {
	user, ok := GetUser(ctx)
	if !ok {
		panic("no user")
	}
	return user
}

func RequireSession(ctx context.Context) *app.SessionInfo {
	sessionInfo, ok := GetSession(ctx)
	if !ok {
		panic("no session")
	}
	return sessionInfo
}

func GetNullableUserID(ctx context.Context) uuid.NullUUID {
	session, ok := GetSession(ctx)
	if !ok {
		return uuid.NullUUID{}
	}
	return uuid.NullUUID{Valid: true, UUID: session.UserID}
}

func WriteSIDCookie(w http.ResponseWriter, cookieName, sid string, expiration time.Duration, secure bool) {
	httpSecure := "; Secure"
	if !secure {
		httpSecure = ""
	}
	w.Header().Add("Set-Cookie", fmt.Sprintf("%s=%s; Path=/; Max-Age=%d; HttpOnly%s", cookieName, sid, int(expiration.Seconds()), httpSecure))
}

func RemoveSIDCookie(w http.ResponseWriter, cookieName string, secure bool) {
	WriteSIDCookie(w, cookieName, "", 0, secure)
}

func attachSessionInfo(r *http.Request, sessionInfo *app.SessionInfo, user *app.SelfUserDto) *http.Request {
	newContext := context.WithValue(r.Context(), sessionInfoKey, sessionInfo)
	newContext = context.WithValue(newContext, userKey, user)
	return r.WithContext(newContext)
}

type sessionInfoKeyType string

const (
	sessionInfoKey sessionInfoKeyType = "ol:session"
	userKey        sessionInfoKeyType = "ol:user"
)
