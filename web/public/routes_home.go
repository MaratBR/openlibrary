package public

import (
	"cmp"
	"context"
	"net/http"
	"slices"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/analytics"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

type homeController struct {
	metricService analytics.MetricService
	bookService   app.BookService
}

func newHomeController(metricService analytics.MetricService, bookService app.BookService) *homeController {
	return &homeController{metricService: metricService, bookService: bookService}
}

func (c *homeController) Register(r chi.Router) {
	r.Get("/", c.homePage)

	r.Get("/cats", func(w http.ResponseWriter, r *http.Request) {
		olhttp.WriteTemplate(w, r.Context(), templates.Cats())
	})

	r.Get("/ui-demo", func(w http.ResponseWriter, r *http.Request) {
		olhttp.WriteTemplate(w, r.Context(), templates.UIDemo())
	})
}

func (c *homeController) homePage(w http.ResponseWriter, r *http.Request) {
	olhttp.WriteTemplate(w, r.Context(), templates.Home([]templ.Component{
		templates.Home_MainHeroSection(),
		c.renderMostViewedWidget(r.Context()),
	}))
}

func (c *homeController) renderMostViewedWidget(ctx context.Context) templ.Component {
	books, views, err := c.getMostViewedBooks(ctx)
	if err != nil {
		return templates.Home_WidgetError(err)
	}

	return templates.Home_MostViewedBooks(books, views)
}

func (c *homeController) getMostViewedBooks(ctx context.Context) ([]app.BookListDto, map[int64]int64, error) {
	bookViewData, err := c.metricService.GetTopBooks(ctx, analytics.GetTopBooksByMetricQuery{
		Metric: analytics.MetricViews,
		Period: store.OlAnalyticsBucketPeriodTypeAll,
		Limit:  10,
	})
	if err != nil {
		return nil, nil, err
	}

	views := make(map[int64]int64)
	bookIds := make([]int64, 0, len(bookViewData))

	for _, entry := range bookViewData {
		bookIds = append(bookIds, entry.BookID)
		views[entry.BookID] = entry.Value.Samples
	}

	books, err := c.bookService.GetBooksById(ctx, bookIds)
	if err != nil {
		return nil, nil, err
	}

	slices.SortStableFunc(books, func(a app.BookListDto, b app.BookListDto) int {
		aViews, _ := views[a.ID]
		bViews, _ := views[b.ID]
		return cmp.Compare(bViews, aViews)
	})

	return books, views, nil
}
