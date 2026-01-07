package olhttp

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
)

func ReqCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newContext := r.Context()

		newContext = context.WithValue(newContext, contextKeyRequest, r)
		newContext = context.WithValue(newContext, contextKeyIDCounter, new(int))

		r = r.WithContext(newContext)
		next.ServeHTTP(w, r)
	})
}

type olhttpContextKey byte

var (
	contextKeyRequest   olhttpContextKey = 1
	contextKeyIDCounter olhttpContextKey = 2
)

func GetRequest(ctx context.Context) *http.Request {
	v := ctx.Value(contextKeyRequest).(*http.Request)
	if v == nil {
		panic("request not found in context")
	}
	return v
}

func GetID(ctx context.Context) string {
	counter := ctx.Value(contextKeyIDCounter)

	if counter == nil {
		return fmt.Sprintf("IDRNG%d", rand.Int())
	} else {
		intPtr := counter.(*int)
		*intPtr++
		return fmt.Sprintf("IDC%d", *intPtr)
	}
}
