package auth

import (
	"context"
	"net/http"

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
