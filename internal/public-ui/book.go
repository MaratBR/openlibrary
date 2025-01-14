package publicui

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/commonutil"
	i18nprovider "github.com/MaratBR/openlibrary/internal/i18n-provider"
	"github.com/MaratBR/openlibrary/internal/public-ui/templates"
)

type bookController struct {
	service       app.BookService
	reviewService app.ReviewsService
}

func newBookController(service app.BookService, reviewService app.ReviewsService) *bookController {
	return &bookController{
		service:       service,
		reviewService: reviewService,
	}
}

func (b *bookController) GetBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := commonutil.URLParamInt64(r, "bookID")
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}

	userID := auth.GetNullableUserID(r.Context())
	book, err := b.service.GetBook(r.Context(), app.GetBookQuery{ID: bookID, ActorUserID: userID})
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}

	var (
		reviews []app.ReviewDto
		review  app.Nullable[app.ReviewDto]
	)

	if userID.Valid {
		review, err = b.reviewService.GetReview(r.Context(), app.GetReviewQuery{
			BookID: bookID,
			UserID: userID.UUID,
		})
		if err != nil {
			write500(w, r, err)
			return
		}
	}

	{
		reviewsResult, err := b.reviewService.GetBookReviews(r.Context(), app.GetBookReviewsQuery{
			BookID:   bookID,
			Page:     1,
			PageSize: 5,
		})
		if err != nil {
			write500(w, r, err)
			return
		}

		reviews = reviewsResult.Reviews
	}

	l := i18nprovider.GetLocalizer(r.Context())
	templates.BookPage(r.Context(), l, book, review, reviews).Render(r.Context(), w)

}
