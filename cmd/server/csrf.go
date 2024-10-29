package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MaratBR/openlibrary/internal/commonutil"
)

func csrfMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			_, err := r.Cookie("csrf")
			if err == http.ErrNoCookie {
				w.Header().Add("Set-Cookie", fmt.Sprintf("csrf=%s; Path=/; Max-Age=%f", generateCsrfToken(), (time.Hour*24).Seconds()))
			}
		} else {
			csrfCookie, err := r.Cookie("csrf")
			if err != nil {
				writeCsrfError(w)
				return
			}

			csrfHeader := r.Header.Get("x-csrf-token")
			if csrfHeader == "" || csrfHeader != csrfCookie.Value {
				writeCsrfError(w)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func writeCsrfError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("csrf token is missing"))
}

func generateCsrfToken() string {
	token, err := commonutil.GenerateRandomStringURLSafe(32)
	if err != nil {
		panic(err)
	}
	return token
}
