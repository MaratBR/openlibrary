package public

import (
	"fmt"
	"net/http"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/flash"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/MaratBR/openlibrary/web/public/templates"
	"github.com/go-chi/chi/v5"
)

type modController struct {
	bookService    app.BookService
	modBookService app.ModerationBookService
}

func newModController(
	booksService app.BookService,
	modBookService app.ModerationBookService,
) *modController {
	return &modController{
		bookService:    booksService,
		modBookService: modBookService,
	}
}

func (c *modController) Register(r chi.Router) {
	r.Get("/mod/book/{bookID}", c.book)
	r.Get("/mod/book/{bookID}/log", c.bookLogs)
	r.Get("/mod/book/{bookID}/mod-action", c.bookPerformAction)
	r.Post("/mod/book/{bookID}/mod-action", c.bookPerformAction)
}

func (c *modController) bookLogs(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}

	s := auth.RequireSession(r.Context())

	query := r.URL.Query()
	page := olhttp.GetPage(query, "p")
	defaultPageSize := uint32(200)
	pageSize := olhttp.GetPageSize(query, "pageSize", 5, 1000, defaultPageSize)

	logs, err := c.modBookService.GetBookLog(r.Context(), app.GetBookLogQuery{
		PageSize:    pageSize,
		Page:        page,
		BookID:      bookID,
		ActorUserID: s.UserID,
	})

	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	book, err := c.modBookService.GetBookInfo(r.Context(), app.GetBookInfoQuery{
		ActorUserID: s.UserID,
		BookID:      bookID,
	})

	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	olhttp.WriteTemplate(w, r.Context(), templates.ModBookLog(book, logs))
}

func (c *modController) book(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("%s/mod-action", r.URL.Path), http.StatusFound)
}

func (c *modController) bookPerformAction(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		writeBadRequest(w, r, err)
		return
	}

	s := auth.RequireSession(r.Context())

	book, err := c.modBookService.GetBookInfo(r.Context(), app.GetBookInfoQuery{
		ActorUserID: s.UserID,
		BookID:      bookID,
	})

	if err != nil {
		writeApplicationError(w, r, err)
		return
	}

	if r.Method == http.MethodPost {
		sentResponse := c.handleBookAction(w, r, book)
		if !sentResponse {
			http.Redirect(w, r, r.URL.Path, http.StatusFound)
			return
		}
	} else {
		logs, err := c.modBookService.GetBookLog(r.Context(), app.GetBookLogQuery{
			ActorUserID: s.UserID,
			PageSize:    5,
			BookID:      book.ID,
		})
		if err != nil {
			writeApplicationError(w, r, err)
			return
		}

		olhttp.WriteTemplate(w, r.Context(), templates.ModBook(book, logs))
	}
}

func (c *modController) handleBookAction(w http.ResponseWriter, r *http.Request, book app.BookModerationInfo) bool {
	err := r.ParseForm()
	if err != nil {
		writeBadRequest(w, r, err)
		return true
	}

	action := r.Form.Get("act")
	reason := r.Form.Get("reason")

	s := auth.RequireSession(r.Context())

	err = nil

	cmd := app.ModerationPerformBookActionCommand{
		ActorUserID: s.UserID,
		Reason:      reason,
		BookID:      book.ID,
	}

	switch action {
	case "ban":
		err = c.modBookService.BanBook(r.Context(), cmd)
	case "unban":
		err = c.modBookService.UnBanBook(r.Context(), cmd)
	case "shadow_ban":
		err = c.modBookService.ShadowBanBook(r.Context(), cmd)
	case "shadow_unban":
		err = c.modBookService.UnShadowBanBook(r.Context(), cmd)
	case "perm_delete":
		err = c.modBookService.PermanentlyRemoveBook(r.Context(), cmd)
	default:
		// TODO localized error message?
		flash.Add(r, flash.Text("Unknown value for act: "+action))
		return false
	}

	if err != nil {
		writeApplicationError(w, r, err)
		return true
	}

	// TODO use actual wording
	flash.Add(r, flash.Text(action+": success"))
	return false
}
