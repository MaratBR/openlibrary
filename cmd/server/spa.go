package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	uiserver "github.com/MaratBR/openlibrary/cmd/server/ui-server"
	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/gofrs/uuid"
	"github.com/knadh/koanf/v2"
	"google.golang.org/protobuf/proto"
)

type serverSessionData struct {
	ExpiresAt time.Time `json:"expiresAt"`
}

type spaHandler struct {
	devHandler uiserver.Handler
}

func newSPAHandler(
	config *koanf.Koanf,
	bookService app.BookService,
	reviewsService app.ReviewsService,
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

		session, ok := auth.GetSession(r.Context())

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
			for _, book := range result.Items {
				if book.Cover != "" {
					data.AddPreloadURL("image", book.Cover)
				}
			}
			b, err := proto.Marshal(result)
			if err == nil {
				data.ServerMetadata["search"] = base64.StdEncoding.EncodeToString(b)
			}
		}

		{
			// bookExtremes, err := searchService.GetBookExtremes(r.Context())
			// if err == nil {
			// 	data.AddPreloadedData("/api/search/book-extremes", bookExtremes)
			// }
		}

		// if len(search.ExcludeTags)+len(search.IncludeTags) > 0 {
		// 	tagIds := commonutil.MergeArrays(search.IncludeTags, search.ExcludeTags)
		// 	tags, err := tagsService.GetTagsByIds(r.Context(), tagIds)
		// 	if err == nil {
		// 		key := "/api/tags/lookup?q=" + url.QueryEscape(i64Array(tagIds))
		// 		data.AddPreloadedData(key, tags)
		// 	}
		// }

	})

	devHandler.PreloadData("/users/__profile", func(r *http.Request, data *uiserver.Data) {
		data.ServerMetadata["iframeAllowed"] = true

		userID, err := urlQueryParamUUID(r, "userId")
		if err == nil && userID != uuid.Nil {
			user, err := userService.GetUserDetails(r.Context(), app.GetUserQuery{
				ID:     userID,
				UserID: auth.GetNullableUserID(r.Context()),
			})
			if err == nil {
				data.AddPreloadedData(fmt.Sprintf("/api/users/%s", user.ID.String()), user)
			}
		}
	})

	// devHandler.PreloadData("/users/{userID}", func(r *http.Request, data *uiserver.Data) {
	// 	userID := chi.URLParam(r, "userID")
	// 	data.AddPreloadURL()
	// })

	devHandler.PreloadData("/book/{bookID}", func(r *http.Request, data *uiserver.Data) {
		bookID, err := urlParamInt64(r, "bookID")
		userID := auth.GetNullableUserID(r.Context())
		if err == nil {
			book, err := bookService.GetBook(r.Context(), app.GetBookQuery{ID: bookID, ActorUserID: userID})
			if err == nil {
				data.AddPreloadedData(fmt.Sprintf("/api/books/%d", bookID), book)
				if book.Cover != "" {
					data.AddPreloadURL("image", book.Cover)
				}
			}

			reviews, err := reviewsService.GetBookReviews(r.Context(), app.GetBookReviewsQuery{
				BookID:   bookID,
				Page:     1,
				PageSize: 5,
			})
			if err == nil {
				data.AddPreloadedData(fmt.Sprintf("/api/reviews/%d", bookID), reviewsResponse{
					Reviews:    reviews.Reviews,
					Pagination: reviews.Pagination,
				})
			}

			if userID.Valid {
				review, err := reviewsService.GetReview(r.Context(), app.GetReviewQuery{
					BookID: bookID,
					UserID: userID.UUID,
				})
				if err == nil {
					data.AddPreloadedData(fmt.Sprintf("/api/reviews/%d/my", bookID), review)
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
