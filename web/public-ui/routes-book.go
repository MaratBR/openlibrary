package publicui

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/commonutil"
	"github.com/MaratBR/openlibrary/web/public-ui/templates"
)

type bookController struct {
	service            app.BookService
	reviewService      app.ReviewsService
	readingListService app.ReadingListService
}

func newBookController(service app.BookService, reviewService app.ReviewsService, readingListService app.ReadingListService) *bookController {
	return &bookController{
		service:            service,
		reviewService:      reviewService,
		readingListService: readingListService,
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
		reviews           []app.ReviewDto
		ratingAndReview   app.RatingAndReview
		readingListStatus app.Nullable[app.BookReadingListDto]
	)

	if userID.Valid {
		ratingAndReview, err = b.reviewService.GetReview(r.Context(), app.GetReviewQuery{
			BookID: bookID,
			UserID: userID.UUID,
		})
		if err != nil {
			write500(w, r, err)
			return
		}

		readingListStatus, err = b.readingListService.GetStatus(r.Context(), userID.UUID, bookID)
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

	templates.BookPage(
		r.Context(),
		book,
		ratingAndReview,
		readingListStatus,
		reviews,
	).Render(r.Context(), w)
}

func (b *bookController) GetBookTOC(w http.ResponseWriter, r *http.Request) {
	bookID, err := commonutil.URLParamInt64(r, "bookID")
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}

	chapters, err := b.service.GetBookChapters(r.Context(), app.GetBookChaptersQuery{
		ID: bookID,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	templates.BookTOC(r.Context(), bookID, chapters).Render(r.Context(), w)
}
