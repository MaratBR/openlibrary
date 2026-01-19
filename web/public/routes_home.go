package public

import (
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/analytics"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/go-chi/chi/v5"
)

type homeController struct {
	viewsService analytics.ViewsService
	bookService  app.BookService
}

func newHomeController(viewsService analytics.ViewsService, bookService app.BookService) *homeController {
	return &homeController{viewsService: viewsService, bookService: bookService}
}

func (c *homeController) Register(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		olhttp.WriteTemplate(w, r.Context(), templates.Home())
	})

	r.Get("/ui-demo", func(w http.ResponseWriter, r *http.Request) {
		olhttp.WriteTemplate(w, r.Context(), templates.UIDemo())
	})

	r.Get("/__fragment/most-viewed-books", c.getMostViewedBooksWidget)
}

func (c *homeController) getMostViewedBooksWidget(w http.ResponseWriter, r *http.Request) {
	period, _ := olhttp.URLQueryParamInt64(r, "period")
	_, _ = olhttp.URLQueryParamInt64(r, "count")

	bookViewData, err := c.viewsService.GetMostViewedBooks(r.Context(), analytics.AnalyticsPeriod(period))
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	views := make(map[int64]int64)
	bookIds := make([]int64, 0, len(bookViewData))

	for _, entry := range bookViewData {
		bookIds = append(bookIds, entry.BookID)
		views[entry.BookID] = entry.Views
	}

	books, err := c.bookService.GetBooksById(r.Context(), bookIds)
	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	olhttp.WriteTemplate(w, r.Context(), templates.HomeBooksListFragment(analytics.AnalyticsPeriod(period), books, views))
}
