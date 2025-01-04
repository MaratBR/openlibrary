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
		PageSize: 20,
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
