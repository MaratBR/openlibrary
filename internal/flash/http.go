package flash

import (
	"context"
	"net/http"
)

type flashCollection struct {
	arr []Message
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "flash:collection", &flashCollection{})))
	})
}

func Add(r *http.Request, message Message) {
	collection := r.Context().Value("flash:collection").(*flashCollection)
	collection.arr = append(collection.arr, message)
}

func GetFlashes(ctx context.Context) []Message {
	collectionAny := ctx.Value("flash:collection")
	if collectionAny == nil {
		return nil
	}
	collection := collectionAny.(*flashCollection)
	return collection.arr
}
