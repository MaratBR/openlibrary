package public

import (
	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"net/http"
	"time"
)

func moderationDateRange(r *http.Request) (app.Nullable[time.Time], app.Nullable[time.Time], error) {
	from, err := moderationDateFilter(r.URL.Query().Get("dateFrom"), false)
	if err != nil {
		return from, app.Null[time.Time](), err
	}
	to, err := moderationDateFilter(r.URL.Query().Get("dateTo"), true)
	if err == nil && from.Valid && to.Valid && !from.Value.Before(to.Value) {
		err = app.ErrModerationFilter
	}
	return from, to, err
}
func (c *apiControllerModeration) overview(w http.ResponseWriter, r *http.Request) {
	result, err := c.investigation.Overview(r.Context(), auth.RequireSession(r.Context()).UserID)
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	olhttp.NewAPIResponse(result).Write(w)
}
func (c *apiControllerModeration) searchChapters(w http.ResponseWriter, r *http.Request) {
	c.searchContent(w, r, "chapter")
}
func (c *apiControllerModeration) searchComments(w http.ResponseWriter, r *http.Request) {
	c.searchContent(w, r, "comment")
}
func (c *apiControllerModeration) searchContent(w http.ResponseWriter, r *http.Request, kind string) {
	p := r.URL.Query()
	result, err := c.investigation.Search(r.Context(), app.ModerationContentQuery{ActorUserID: auth.RequireSession(r.Context()).UserID, Kind: kind, Search: p.Get("search"), BookID: p.Get("bookId"), AuthorID: p.Get("authorId"), ChapterID: p.Get("chapterId"), State: p.Get("state"), Page: olhttp.GetPage(p, "page"), PageSize: olhttp.GetPageSize(p, "pageSize", 1, 100, 25)})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	olhttp.NewAPIResponse(result).Write(w)
}
func (c *apiControllerModeration) getChapter(w http.ResponseWriter, r *http.Request) {
	c.getContent(w, r, "chapter")
}
func (c *apiControllerModeration) getComment(w http.ResponseWriter, r *http.Request) {
	c.getContent(w, r, "comment")
}
func (c *apiControllerModeration) getContent(w http.ResponseWriter, r *http.Request, kind string) {
	id, err := olhttp.URLParamInt64(r, kind+"ID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	result, err := c.investigation.Get(r.Context(), app.ModerationContentQuery{ActorUserID: auth.RequireSession(r.Context()).UserID, Kind: kind}, id)
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	olhttp.NewAPIResponse(result).Write(w)
}

// go2tsdef:generate
type ModerationReportNoteRequest struct {
	Note     string `json:"note"`
	Priority string `json:"priority"`
	Reason   string `json:"reason"`
}

func (c *apiControllerModeration) reportNote(w http.ResponseWriter, r *http.Request) {
	id, err := olhttp.URLParamInt64(r, "reportID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	var input ModerationReportNoteRequest
	if err = olhttp.ReadJSONBodyMax(w, r, &input, moderationJSONMaxBytes); err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	err = c.reports.AddNote(r.Context(), app.ReportNoteCommand{ActorUserID: auth.RequireSession(r.Context()).UserID, ReportID: id, Note: input.Note, Priority: input.Priority, Reason: input.Reason})
	writeModerationActionResult(w, err)
}
