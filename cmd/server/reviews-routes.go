package main

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
)

type reviewsController struct {
	service app.ReviewsService
}

func newReviewsController(service app.ReviewsService) *reviewsController {
	return &reviewsController{service: service}
}

func (c *reviewsController) GetMyReview(w http.ResponseWriter, r *http.Request) {
	bookID, err := urlParamInt64(r, "bookID")
	if err != nil {
		writeUnprocessableEntity(w, err.Error())
		return
	}

	session, ok := getSession(r)
	if !ok {
		writeUnauthorizedError(w)
		return
	}

	review, err := c.service.GetReview(r.Context(), app.GetReviewQuery{
		BookID: bookID,
		UserID: session.UserID,
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, review)
}

func (c *reviewsController) GetReviewsDistribution(w http.ResponseWriter, r *http.Request) {
	bookID, err := urlParamInt64(r, "bookID")
	if err != nil {
		writeUnprocessableEntity(w, err.Error())
		return
	}

	distribution, err := c.service.GetBookReviewsDistribution(r.Context(), bookID)
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, distribution.Distribution)
}

type reviewsResponse struct {
	Reviews    []app.ReviewDto       `json:"reviews"`
	Pagination app.PaginationOptions `json:"pagination"`
}

func (c *reviewsController) GetReviews(w http.ResponseWriter, r *http.Request) {
	bookID, err := urlParamInt64(r, "bookID")
	if err != nil {
		writeUnprocessableEntity(w, err.Error())
		return
	}

	reviews, err := c.service.GetBookReviews(r.Context(), app.GetBookReviewsQuery{
		BookID:   bookID,
		Page:     1,
		PageSize: 5,
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, reviewsResponse{
		Reviews: reviews.Reviews,
	})
}

type createReviewRequest struct {
	Content string          `json:"content"`
	Rate    app.RatingValue `json:"rating"`
}

func (c *reviewsController) UpdateOrCreateReview(w http.ResponseWriter, r *http.Request) {
	request := createReviewRequest{}
	if err := readJSON(r, &request); err != nil {
		writeUnprocessableEntity(w, err.Error())
		return
	}
	bookID, err := urlParamInt64(r, "bookID")
	if err != nil {
		writeUnprocessableEntity(w, err.Error())
		return
	}
	review, err := c.service.UpdateReview(r.Context(), app.UpdateReviewCommand{
		BookID:  bookID,
		UserID:  requireSession(r).UserID,
		Content: request.Content,
		Rating:  request.Rate,
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeJSON(w, review)
}

func (c *reviewsController) DeleteReview(w http.ResponseWriter, r *http.Request) {
	bookID, err := urlParamInt64(r, "bookID")
	if err != nil {
		writeUnprocessableEntity(w, err.Error())
		return
	}
	err = c.service.DeleteReview(r.Context(), app.DeleteReviewCommand{
		BookID: bookID,
		UserID: requireSession(r).UserID,
	})
	if err != nil {
		writeApplicationError(w, err)
		return
	}
	writeOK(w)
}
