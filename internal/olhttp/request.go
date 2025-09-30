package olhttp

import (
	"context"
	"net/http"
)

func ReqCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newContext := context.WithValue(r.Context(), "http.req", r)
		r = r.WithContext(newContext)
		next.ServeHTTP(w, r)
	})
}

func GetRequest(ctx context.Context) *http.Request {
	v := ctx.Value("http.req").(*http.Request)
	if v == nil {
		panic("request not found in context")
	}
	return v
}
