package session

import (
	"fmt"
	"net/http"
	"time"
)

type redactedString string

func (c redactedString) String() string {
	return "[REDACTED]"
}

func (c redactedString) MarshalJSON() ([]byte, error) {
	return []byte("\"[REDACTED]\""), nil
}

func WriteSIDCookie(w http.ResponseWriter, sid string, expiration time.Duration, secure bool) {
	httpSecure := "; Secure"
	if !secure {
		httpSecure = ""
	}
	w.Header().Add("Set-Cookie", fmt.Sprintf("%s=%s; Path=/; Max-Age=%d; HttpOnly%s", CookieName, sid, int(expiration.Seconds()), httpSecure))
}

func RemoveSIDCookie(w http.ResponseWriter, secure bool) {
	WriteSIDCookie(w, "", 0, secure)
}
