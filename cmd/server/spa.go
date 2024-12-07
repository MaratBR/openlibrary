package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	uiserver "github.com/MaratBR/openlibrary/cmd/server/ui-server"
	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/commonutil"
	"github.com/gofrs/uuid"
	"github.com/knadh/koanf/v2"
)

type serverSessionData struct {
	ExpiresAt time.Time `json:"expiresAt"`
}

type spaHandler struct {
	devHandler uiserver.Handler
}

func newSPAHandler(
	config *koanf.Koanf,
	bookService *app.BookService,
	userService app.UserService,
	searchService app.SearchService,
	tagsService app.TagsService,
) spaHandler {
	h := spaHandler{}
	devHandler := uiserver.NewDevStaticHandler(config)

	devHandler.DefaultPreloadData(func(r *http.Request, data *uiserver.Data) {
		data.ServerMetadata["session"] = nil
		data.ServerMetadata["user"] = nil
		data.ServerMetadata["clientPreload"] = true
		data.ServerMetadata["serverPreload"] = true

		session, ok := getSession(r)

		if ok {
			user, err := userService.GetUserSelfData(r.Context(), session.UserID)
			if err == nil {
				data.ServerMetadata["session"] = serverSessionData{
					ExpiresAt: session.ExpiresAt,
				}
				data.ServerMetadata["user"] = user
			}
		}
	})

	devHandler.PreloadData("/search", func(r *http.Request, data *uiserver.Data) {
		search, err := getSearchRequest(r)
		if err != nil {
			return
		}
		result, err := performBookSearch(searchService, tagsService, r, search)
		if err != nil {
			return
		}

		{
			key := "/api/search?" + strings.ReplaceAll(r.URL.RawQuery, "%20", "+")
			data.AddPreloadedData(key, result)

			for _, book := range result.Books {
				if book.Cover != "" {
					data.AddPreloadURL("image", book.Cover)
				}
			}
			data.AddPreloadedData(key, result)
		}

		{
			bookExtremes, err := searchService.GetBookExtremes(r.Context())
			if err == nil {
				data.AddPreloadedData("/api/search/book-extremes", bookExtremes)
			}
		}

		if len(search.ExcludeTags)+len(search.IncludeTags) > 0 {
			tagIds := commonutil.MergeArrays(search.IncludeTags, search.ExcludeTags)
			tags, err := tagsService.GetTagsByIds(r.Context(), tagIds)
			if err == nil {
				key := "/api/tags/lookup?q=" + url.QueryEscape(i64Array(tagIds))
				data.AddPreloadedData(key, tags)
			}

		}

	})

	devHandler.PreloadData("/user/__profile", func(r *http.Request, data *uiserver.Data) {
		data.ServerMetadata["iframeAllowed"] = true

		userID, err := urlQueryParamUUID(r, "userId")
		if err == nil && userID != uuid.Nil {
			user, err := userService.GetUserDetails(r.Context(), app.GetUserQuery{
				ID:     userID,
				UserID: getNullableUserID(r),
			})
			if err == nil {
				data.AddPreloadedData(fmt.Sprintf("/api/users/%s", user.ID.String()), user)
			}
		}
	})

	// devHandler.PreloadData("/user/{userID}", func(r *http.Request, data *uiserver.Data) {
	// 	userID := chi.URLParam(r, "userID")
	// 	data.AddPreloadURL()
	// })

	devHandler.PreloadData("/book/{bookID}", func(r *http.Request, data *uiserver.Data) {
		bookID, err := urlParamInt64(r, "bookID")
		userID := getNullableUserID(r)
		if err == nil {
			book, err := bookService.GetBook(r.Context(), app.GetBookQuery{ID: bookID, ActorUserID: userID})
			if err == nil {
				data.AddPreloadedData(fmt.Sprintf("/api/books/%d", bookID), book)
				if book.Cover != "" {
					data.AddPreloadURL("image", book.Cover)
				}
			}
		}
	})

	h.devHandler = *devHandler

	return h
}

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.devHandler.ServeHTTP(w, r)
}
