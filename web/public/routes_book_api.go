package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/commonutil"
	olhttp "github.com/MaratBR/openlibrary/internal/olhttp"
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

	bookID, err := olhttp.URLQueryParamInt64(r, "bookId")
	if err != nil {
		apiwriteBadRequest(w, err)
		return
	}

	rating, err := olhttp.URLQueryParamInt64(r, "rating")
	if err != nil {
		apiwriteBadRequest(w, err)
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

func (c *apiBookController) UpdateReview(w http.ResponseWriter, r *http.Request) {
	session := auth.RequireSession(r.Context())

	bookID, err := olhttp.URLQueryParamInt64(r, "bookId")
	if err != nil {
		apiwriteBadRequest(w, err)
		return
	}

	rating, err := olhttp.URLQueryParamInt64(r, "rating")
	if err != nil {
		apiwriteBadRequest(w, err)
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

type createReviewRequest struct {
	Content string          `json:"content"`
	Rate    app.RatingValue `json:"rating"`
}

func (c *apiBookController) UpdateOrCreateReview(w http.ResponseWriter, r *http.Request) {
	request := createReviewRequest{}
	if err := readJSON(r, &request); err != nil {
		apiWriteUnprocessableEntity(w, err.Error())
		return
	}
	bookID, err := commonutil.URLParamInt64(r, "bookID")
	if err != nil {
		apiWriteUnprocessableEntity(w, err.Error())
		return
	}
	review, err := c.reviewService.UpdateReview(r.Context(), app.UpdateReviewCommand{
		BookID:  bookID,
		UserID:  auth.RequireSession(r.Context()).UserID,
		Content: request.Content,
		Rating:  request.Rate,
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	apiWriteJSON(w, review)
}

func (c *apiBookController) DeleteReview(w http.ResponseWriter, r *http.Request) {
	bookID, err := commonutil.URLParamInt64(r, "bookID")
	if err != nil {
		apiWriteUnprocessableEntity(w, err.Error())
		return
	}
	err = c.reviewService.DeleteReview(r.Context(), app.DeleteReviewCommand{
		BookID: bookID,
		UserID: auth.RequireSession(r.Context()).UserID,
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	apiWriteOK(w)
}
