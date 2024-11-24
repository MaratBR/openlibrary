package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/gofrs/uuid"
)

func getSession(r *http.Request) (*app.SessionInfo, bool) {
	value := r.Context().Value(sessionInfoKey)
	if value == nil {
		return nil, false
	}

	sessionInfo := value.(*app.SessionInfo)
	return sessionInfo, true
}

func requireSession(r *http.Request) *app.SessionInfo {
	sessionInfo, ok := getSession(r)
	if !ok {
		panic("no session")
	}
	return sessionInfo
}

func getNullableUserID(r *http.Request) uuid.NullUUID {
	session, ok := getSession(r)
	if !ok {
		return uuid.NullUUID{}
	}
	return uuid.NullUUID{Valid: true, UUID: session.UserID}
}

func writeSidCookie(w http.ResponseWriter, cookieName, sid string, expiration time.Duration, secure bool) {
	httpSecure := "; Secure"
	if !secure {
		httpSecure = ""
	}
	w.Header().Add("Set-Cookie", fmt.Sprintf("%s=%s; Path=/; Max-Age=%d; HttpOnly%s", cookieName, sid, int(expiration.Seconds()), httpSecure))
}

func removeSidCookie(w http.ResponseWriter, cookieName string, secure bool) {
	writeSidCookie(w, cookieName, "", 0, secure)
}
