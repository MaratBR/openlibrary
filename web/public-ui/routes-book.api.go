package publicui

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	httputil "github.com/MaratBR/openlibrary/internal/http-util"
)

type apiBookController struct {
	service            app.BookService
	reviewService      app.ReviewsService
	readingListService app.ReadingListService
}

func newAPIBookController(service app.BookService, reviewService app.ReviewsService, readingListService app.ReadingListService) *apiBookController {
	return &apiBookController{
		service:            service,
		reviewService:      reviewService,
		readingListService: readingListService,
	}
}

func (c *apiBookController) RateBook(w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())

	bookID, err := httputil.URLQueryParamInt64(r, "bookId")
	if err != nil {
		apiWriteRequestError(w, err)
		return
	}

	rating, err := httputil.URLQueryParamInt64(r, "rating")
	if err != nil {
		apiWriteRequestError(w, err)
		return
	}

	err = c.reviewService.UpdateRating(r.Context(), app.UpdateRatingCommand{
		BookID: bookID,
		UserID: session.UserID,
		Rating: app.CreateRatingValue(int16(rating)),
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}

	apiWriteOK(w)
}
