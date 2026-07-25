package public

import (
	"context"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/analytics"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/a-h/templ"
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
	period := analytics.ANALYTICS_PERIOD_TOTAL
	books, views, err := c.getMostViewedBooks(ctx, period)
	if err != nil {
		return templates.Home_WidgetError(err)
	}

	return templates.Home_MostViewedBooks(analytics.ANALYTICS_PERIOD_TOTAL, books, views)
}

func (c *homeController) getMostViewedBooks(ctx context.Context, period analytics.AnalyticsPeriod) ([]app.BookListDto, map[int64]int64, error) {
	bookViewData, err := c.viewsService.GetMostViewedBooks(ctx, period)
	if err != nil {
		return nil, nil, err
	}

	views := make(map[int64]int64)
	bookIds := make([]int64, 0, len(bookViewData))

	for _, entry := range bookViewData {
		bookIds = append(bookIds, entry.BookID)
		views[entry.BookID] = entry.Views
	}

	books, err := c.bookService.GetBooksById(ctx, bookIds)
	if err != nil {
		return nil, nil, err
	}

	return books, views, nil
}
