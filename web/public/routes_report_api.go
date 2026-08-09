package public

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/go-chi/chi/v5"
)

type apiControllerReports struct{ reports app.ReportService }

func newAPIReportsController(reports app.ReportService) *apiControllerReports {
	return &apiControllerReports{reports: reports}
}

func (c *apiControllerReports) Register(r chi.Router) {
	r.Get("/reports/reasons", c.reasons)
	r.With(requiresAuthorizationMiddleware).Post("/reports", c.create)
}

func (c *apiControllerReports) reasons(w http.ResponseWriter, r *http.Request) {
	reasons, err := c.reports.GetReasons(app.ReportTargetType(r.URL.Query().Get("targetType")))
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	olhttp.NewAPIResponse(reasons).Write(w)
}

func (c *apiControllerReports) create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TargetType  app.ReportTargetType `json:"targetType"`
		TargetID    string               `json:"targetId"`
		Reason      string               `json:"reason"`
		Description string               `json:"description"`
		ChapterID   *int64               `json:"chapterId,string"`
		Excerpt     string               `json:"excerpt"`
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16*1024))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	command := app.CreateReportCommand{ReporterUserID: auth.RequireSession(r.Context()).UserID, TargetType: input.TargetType, TargetID: input.TargetID, Reason: input.Reason, Description: input.Description, BookExcerpt: input.Excerpt}
	if input.ChapterID != nil {
		command.BookChapterID = app.Value(*input.ChapterID)
	}
	report, err := c.reports.Create(r.Context(), command)
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	olhttp.NewAPIResponse(map[string]string{"id": strconv.FormatInt(report.ID, 10), "number": report.Number}).Write(w)
}
