package public

import (
	"log/slog"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/go-chi/chi/v5"
	"github.com/joomcode/errorx"
)

type bookController struct {
	service            app.BookService
	reviewService      app.ReviewsService
	readingListService app.ReadingListService
	analytics          app.AnalyticsViewsService
}

func newBookController(service app.BookService, reviewService app.ReviewsService, readingListService app.ReadingListService, analytics app.AnalyticsViewsService) *bookController {
	return &bookController{
		service:            service,
		reviewService:      reviewService,
		readingListService: readingListService,
		analytics:          analytics,
	}
}

func (b *bookController) Register(r chi.Router) {
	// book page and its fragments
	r.Get("/book/{bookID}", b.book)
	r.Get("/book/{bookID}/__fragment/preview-card", b.bookPreview)
	r.Get("/book/{bookID}/__fragment/toc", b.bookTOC)
	r.Get("/book/{bookID}/__fragment/review", b.bookReview)
}

func (b *bookController) book(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "bookID")
	bookID, slug := olhttp.ParseInt64Slug(param)
	if bookID == 0 {
		olhttp.WriteTemplate(w, r.Context(), templates.BookNotFoundPage())
		return
	}

	userID := auth.GetNullableUserID(r.Context())
	book, err := b.service.GetBookDetails(r.Context(), app.GetBookQuery{ID: bookID, ActorUserID: userID})
	if err != nil {
		if errorx.IsOfType(err, app.ErrTypeBookNotFound) || errorx.IsOfType(err, app.ErrTypeBookPrivated) {
			// send 404 page
			olhttp.WriteTemplate(w, r.Context(), templates.BookNotFoundPage())
		} else {
			writeApplicationError(w, r, err)
			// send generic application error
		}
		return
	}

	var replaceURLWithSlug bool
	if book.Slug != slug {
		replaceURLWithSlug = true
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

	ip := olhttp.GetIP(r)
	b.analytics.IncrBookView(r.Context(), bookID, userID, ip)

	views, err := b.analytics.GetBookViews(r.Context(), bookID)

	templates.BookPage(
		book,
		views,
		ratingAndReview,
		readingListStatus,
		reviews,
		replaceURLWithSlug,
	).Render(r.Context(), w)
}

func (b *bookController) bookTOC(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLParamInt64(r, "bookID")
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

	s, ok := auth.GetSession(r.Context())
	var (
		activeChapterID int64
	)

	if ok {
		status, err := b.readingListService.GetStatus(r.Context(), s.UserID, bookID)
		if err != nil {
			slog.Error("failed to get reading list status", "userID", s.UserID, "bookID", bookID, "err", err)
		} else if status.Valid && status.Value.ChapterID.Valid {
			activeChapterID = int64(status.Value.ChapterID.Value)
		}
	}

	templates.BookTOC(r.Context(), bookID, chapters, activeChapterID).Render(r.Context(), w)
}

func (b *bookController) bookPreview(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}

	book, err := b.service.GetBookDetails(r.Context(), app.GetBookQuery{
		ID:          bookID,
		ActorUserID: auth.GetNullableUserID(r.Context()),
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	templates.BookPreviewPartial(book).Render(r.Context(), w)
}

func (b *bookController) bookReview(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}

	session, ok := auth.GetSession(r.Context())
	if !ok {
		writeUnauthorizedError(w)
		return
	}

	review, err := b.reviewService.GetReview(r.Context(), app.GetReviewQuery{
		BookID: bookID,
		UserID: session.UserID,
	})
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	templates.BookMyReview(bookID, review).Render(r.Context(), w)
}
