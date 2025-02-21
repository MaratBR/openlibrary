package csrf

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MaratBR/openlibrary/internal/commonutil"
)

type csrfTokenKeyType struct{}

var (
	csrfTokenKey csrfTokenKeyType
)

type Handler struct {
	secret       string
	cookie       string
	header       string
	paramName    string
	sidCookie    string
	anonymousSid string
}

func (h Handler) SIDCookie() string {
	return h.sidCookie
}

func (h *Handler) getCSRF(r *http.Request) string {
	cookie, err := r.Cookie(h.cookie)
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	return ""
}

func (h *Handler) getSubmittedCSRF(r *http.Request) string {
	var dsValue string // double submitted value

	dsValue = r.Header.Get(h.header)
	if dsValue != "" {
		return dsValue
	}

	err := r.ParseForm()
	if err == nil {
		dsValue = r.Form.Get(h.paramName)
		if dsValue != "" {
			return dsValue
		}
	}

	return ""
}

func (h *Handler) Verify(r *http.Request) bool {
	var (
		sid            string
		csrfFromCookie string
	)

	{
		cookie, err := r.Cookie(h.cookie)
		if err != nil {
			return false
		}
		csrfFromCookie = cookie.Value
	}

	{
		sidCookie, err := r.Cookie(h.sidCookie)
		if err != nil {
			sid = h.anonymousSid
		} else {
			sid = sidCookie.Value
		}
	}

	dsValue := h.getSubmittedCSRF(r)
	if dsValue == "" || dsValue != csrfFromCookie {
		return false
	}

	isValidToken := VerifyHMACCsrfToken(csrfFromCookie, sid, h.secret)
	return isValidToken
}

func (h *Handler) WriteCSRFFromSid(r *http.Request, w http.ResponseWriter) string {
	return h.WriteCSRFToken(w, h.getSID(r))
}

func (h *Handler) WriteCSRFToken(
	w http.ResponseWriter, sid string,
) string {
	token := GenerateHMACCsrfToken(sid, h.secret)
	w.Header().Add("Set-Cookie", fmt.Sprintf("%s=%s; Path=/; Max-Age=%f", h.cookie, token, (time.Hour*24).Seconds()))
	return token
}

func (h *Handler) WriteAnonymousCSRFToken(
	w http.ResponseWriter,
) {
	h.WriteCSRFToken(w, h.anonymousSid)
}

func (h *Handler) getSID(r *http.Request) string {
	var sid string
	sidCookie, err := r.Cookie("sid")
	if err != nil || sidCookie.Value == "" {
		sid = h.anonymousSid
	} else {
		sid = sidCookie.Value
	}
	return sid
}

func (h *Handler) CheckEndpoint(w http.ResponseWriter, r *http.Request) {
	if !h.Verify(r) {
		writeCsrfError(w)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func (h *Handler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csrfToken := h.getCSRF(r)

		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			if r.Method == http.MethodGet {
				_, err := r.Cookie(h.cookie)
				if err == http.ErrNoCookie {
					csrfToken = h.WriteCSRFFromSid(r, w)
				}
			}
		} else {
			if !h.Verify(r) {
				writeCsrfError(w)
				return
			}
		}

		newContext := context.WithValue(r.Context(), csrfTokenKey, csrfToken)
		r = r.WithContext(newContext)

		next.ServeHTTP(w, r)
	})
}

func NewHandler(secret string) *Handler {
	return &Handler{
		secret:    secret,
		cookie:    "csrf",
		header:    "x-csrf-token",
		paramName: "__csrf",
		sidCookie: "sid",
	}
}

func writeCsrfError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("csrf token is missing"))
}

func getCsrfHMACPayload(sid, secret, randomValue string) string {
	var sb strings.Builder
	sb.WriteString(strconv.FormatInt(int64(len(randomValue)), 10))
	sb.WriteRune('!')
	sb.WriteString(randomValue)
	sb.WriteRune('!')
	sb.WriteString(strconv.FormatInt(int64(len(sid)), 10))
	sb.WriteRune('!')
	sb.WriteString(sid)
	return sb.String()
}

func GenerateHMACCsrfToken(
	sid, secret string,
) string {
	randomValue, err := commonutil.GenerateRandomStringURLSafe(32)
	if err != nil {
		panic(err)
	}
	payload := getCsrfHMACPayload(sid, secret, randomValue)
	sha256hmac := hmac.New(sha256.New, []byte(secret))
	sha256hmac.Write([]byte(payload))
	csrfToken := base64.URLEncoding.EncodeToString(sha256hmac.Sum(nil)) + "." + randomValue
	return csrfToken
}

func VerifyHMACCsrfToken(
	csrfToken, sid, secret string,
) bool {
	parts := strings.Split(csrfToken, ".")
	if len(parts) != 2 {
		return false
	}
	sha256hmac := hmac.New(sha256.New, []byte(secret))
	payload := getCsrfHMACPayload(sid, secret, parts[1])
	sha256hmac.Write([]byte(payload))
	expectedHash := base64.URLEncoding.EncodeToString(sha256hmac.Sum(nil))
	return parts[0] == expectedHash
}
