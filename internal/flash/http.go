package flash

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/session"
)

type flashCollection struct {
	arr     []Message
	touched bool
}

func (c *flashCollection) Add(message Message) {
	c.arr = append(c.arr, message)
	c.touched = true
}

func (c *flashCollection) PullAll() []Message {
	arr := c.arr
	c.arr = nil
	c.touched = true
	return arr
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, ok := session.Get(r)
		if !ok {
			slog.Error("could not load flash messages to session: session is not attached to request context")
			next.ServeHTTP(w, r)
			return
		}

		col := &flashCollection{touched: false}
		r = r.WithContext(context.WithValue(r.Context(), "flash:collection", col))

		// get list of flash messages
		value, ok := s.Get("flash:collection")

		if ok && value != "" {
			err := json.Unmarshal([]byte(value), &col.arr)
			if err != nil {
				slog.Error("failed to unmarshal flash messages from session", "err", err)
			}
		}

		// run the handler
		next.ServeHTTP(w, r)

		if col.touched {
			// save flash messages to session
			b, err := json.Marshal(col.arr)
			if err != nil {
				slog.Error("failed to serialized flashes list", "err", err)
				return
			}
			s.Put("flash:collection", string(b))
		}
	})
}

func Add(r *http.Request, message Message) {
	collectionAny := r.Context().Value("flash:collection")
	if collectionAny == nil {
		panic("cannot find flash messages collection")
	}
	collection := collectionAny.(*flashCollection)
	collection.Add(message)
}

func PullFlashes(ctx context.Context) []Message {
	collectionAny := ctx.Value("flash:collection")
	if collectionAny == nil {
		return nil
	}
	collection := collectionAny.(*flashCollection)
	return collection.PullAll()
}
