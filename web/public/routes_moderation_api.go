package public

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/auth"
	"github.com/MaratBR/openlibrary/internal/olhttp"
	"github.com/go-chi/chi/v5"
)

type apiControllerModeration struct {
	books        app.ModerationBookService
	content      app.ContentModerationService
	loginHistory app.LoginHistoryService
}

func newAPIModerationController(books app.ModerationBookService, content app.ContentModerationService, loginHistory app.LoginHistoryService) *apiControllerModeration {
	return &apiControllerModeration{books: books, content: content, loginHistory: loginHistory}
}

func (c *apiControllerModeration) Register(r chi.Router) {
	r.Route("/moderation", func(r chi.Router) {
		r.Use(apiRequiresAuthorizationMiddleware)
		r.Get("/books/{bookID}", c.getBook)
		r.Get("/books/{bookID}/log", c.getBookLog)
		r.Post("/books/{bookID}/actions/{action}", c.bookAction)
		r.Post("/chapters/{chapterID}/actions/{action}", c.chapterAction)
		r.Post("/comments/{commentID}/actions/{action}", c.commentAction)
		r.Post("/users/{userID}/actions/{action}", c.userAction)
		r.Get("/users/{userID}/login-history", c.getLoginHistory)
	})
}

type moderationReasonRequest struct {
	Reason string `json:"reason"`
}
type moderationBanRequest struct {
	Reason string    `json:"reason"`
	Until  time.Time `json:"until"`
}
type moderationValueRequest struct {
	Reason string `json:"reason"`
	Value  string `json:"value"`
}

// go2tsdef:generate
type ModerationBookResponse struct {
	ID                   int64  `json:"id,string"`
	Name                 string `json:"name"`
	Summary              string `json:"summary"`
	IsBanned             bool   `json:"isBanned"`
	IsShadowBanned       bool   `json:"isShadowBanned"`
	IsPermanentlyRemoved bool   `json:"isPermanentlyRemoved"`
}

// go2tsdef:generate
type BookModerationLogEntryResponse struct {
	Time          time.Time       `json:"time"`
	Action        string          `json:"action"`
	Payload       json.RawMessage `json:"payload" go2tsdef:"unknown"`
	Reason        string          `json:"reason"`
	ActorUserID   string          `json:"actorUserId"`
	ActorUserName string          `json:"actorUserName"`
}

// go2tsdef:generate
type BookModerationLogResponse struct {
	Entries         []BookModerationLogEntryResponse `json:"entries"`
	Page            int32                            `json:"page"`
	PageSize        int32                            `json:"pageSize"`
	HasNextPage     bool                             `json:"hasNextPage"`
	HasPreviousPage bool                             `json:"hasPreviousPage"`
	TotalPages      uint32                           `json:"totalPages"`
}

// go2tsdef:generate
type LoginHistoryEntryResponse struct {
	IPAddress  string    `json:"ipAddress"`
	UserAgent  string    `json:"userAgent"`
	LoggedInAt time.Time `json:"loggedInAt"`
}

// go2tsdef:generate
type UserLoginHistoryResponse []LoginHistoryEntryResponse

func decodeModerationJSON(w http.ResponseWriter, r *http.Request, value any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64*1024))
	decoder.DisallowUnknownFields()
	return decoder.Decode(value)
}

func (c *apiControllerModeration) getBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	result, err := c.books.GetBookInfo(r.Context(), app.GetBookInfoQuery{ActorUserID: auth.RequireSession(r.Context()).UserID, BookID: bookID})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	olhttp.NewAPIResponse(ModerationBookResponse{
		ID: result.ID, Name: result.Name, Summary: result.Summary,
		IsBanned: result.IsBanned, IsShadowBanned: result.IsShadowBanned,
		IsPermanentlyRemoved: result.IsPermDeleted,
	}).Write(w)
}

func (c *apiControllerModeration) getBookLog(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	result, err := c.books.GetBookLog(r.Context(), app.GetBookLogQuery{ActorUserID: auth.RequireSession(r.Context()).UserID, BookID: bookID, Page: olhttp.GetPage(r.URL.Query(), "page"), PageSize: olhttp.GetPageSize(r.URL.Query(), "pageSize", 1, 100, 25)})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	entries := app.MapSlice(result.Entries, func(entry app.BookModerationLog) BookModerationLogEntryResponse {
		return BookModerationLogEntryResponse{
			Time: entry.Time, Action: string(entry.Action), Payload: entry.Payload,
			Reason: entry.Reason, ActorUserID: entry.ActorUserID.String(), ActorUserName: entry.ActorUserName,
		}
	})
	olhttp.NewAPIResponse(BookModerationLogResponse{
		Entries: entries, Page: result.Page, PageSize: result.PageSize,
		HasNextPage: result.HasNextPage, HasPreviousPage: result.HasPreviousPage,
		TotalPages: result.TotalPages,
	}).Write(w)
}

func (c *apiControllerModeration) bookAction(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	var input moderationReasonRequest
	if err = decodeModerationJSON(w, r, &input); err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	cmd := app.ModerationPerformBookActionCommand{ActorUserID: auth.RequireSession(r.Context()).UserID, BookID: bookID, Reason: input.Reason}
	switch chi.URLParam(r, "action") {
	case "ban":
		err = c.books.BanBook(r.Context(), cmd)
	case "unban":
		err = c.books.UnBanBook(r.Context(), cmd)
	case "shadow-ban":
		err = c.books.ShadowBanBook(r.Context(), cmd)
	case "unshadow-ban":
		err = c.books.UnShadowBanBook(r.Context(), cmd)
	case "permanent-remove":
		err = c.books.PermanentlyRemoveBook(r.Context(), cmd)
	default:
		apiWriteBadRequest(w, errors.New("unknown book moderation action"))
		return
	}
	writeModerationActionResult(w, err)
}

func (c *apiControllerModeration) chapterAction(w http.ResponseWriter, r *http.Request) {
	id, err := olhttp.URLParamInt64(r, "chapterID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	var input moderationReasonRequest
	if err = decodeModerationJSON(w, r, &input); err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	cmd := app.ModerateChapterCommand{ActorUserID: auth.RequireSession(r.Context()).UserID, ChapterID: id, Reason: input.Reason}
	switch chi.URLParam(r, "action") {
	case "hide":
		err = c.content.HideChapter(r.Context(), cmd)
	case "restore":
		err = c.content.RestoreChapter(r.Context(), cmd)
	default:
		apiWriteBadRequest(w, errors.New("unknown chapter moderation action"))
		return
	}
	writeModerationActionResult(w, err)
}

func (c *apiControllerModeration) commentAction(w http.ResponseWriter, r *http.Request) {
	id, err := olhttp.URLParamInt64(r, "commentID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	var input moderationReasonRequest
	if err = decodeModerationJSON(w, r, &input); err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	cmd := app.ModerateCommentCommand{ActorUserID: auth.RequireSession(r.Context()).UserID, CommentID: id, Reason: input.Reason}
	switch chi.URLParam(r, "action") {
	case "remove":
		err = c.content.RemoveComment(r.Context(), cmd)
	case "restore":
		err = c.content.RestoreComment(r.Context(), cmd)
	default:
		apiWriteBadRequest(w, errors.New("unknown comment moderation action"))
		return
	}
	writeModerationActionResult(w, err)
}

func (c *apiControllerModeration) userAction(w http.ResponseWriter, r *http.Request) {
	userID, err := olhttp.URLParamUUID(r, "userID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	actor := auth.RequireSession(r.Context()).UserID
	switch chi.URLParam(r, "action") {
	case "ban":
		var input moderationBanRequest
		if err = decodeModerationJSON(w, r, &input); err != nil {
			apiWriteBadRequest(w, err)
			return
		}
		cmd := app.BanUserCommand{ActorUserID: actor, UserID: userID, Until: input.Until, Reason: input.Reason}
		err = c.content.BanUser(r.Context(), cmd)
	case "permanent-ban", "unban":
		var input moderationReasonRequest
		if err = decodeModerationJSON(w, r, &input); err != nil {
			apiWriteBadRequest(w, err)
			return
		}
		cmd := app.BanUserCommand{ActorUserID: actor, UserID: userID, Reason: input.Reason}
		if chi.URLParam(r, "action") == "permanent-ban" {
			err = c.content.PermanentlyBanUser(r.Context(), cmd)
		} else {
			err = c.content.UnbanUser(r.Context(), cmd)
		}
	case "rename", "change-about":
		var input moderationValueRequest
		if err = decodeModerationJSON(w, r, &input); err != nil {
			apiWriteBadRequest(w, err)
			return
		}
		cmd := app.ModerateUserProfileCommand{ActorUserID: actor, UserID: userID, Value: input.Value, Reason: input.Reason}
		if chi.URLParam(r, "action") == "rename" {
			err = c.content.RenameUser(r.Context(), cmd)
		} else {
			err = c.content.ChangeUserAbout(r.Context(), cmd)
		}
	default:
		apiWriteBadRequest(w, errors.New("unknown user moderation action"))
		return
	}
	writeModerationActionResult(w, err)
}

func (c *apiControllerModeration) getLoginHistory(w http.ResponseWriter, r *http.Request) {
	userID, err := olhttp.URLParamUUID(r, "userID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	result, err := c.loginHistory.GetUserLoginHistory(r.Context(), app.GetLoginHistoryQuery{ActorUserID: auth.RequireSession(r.Context()).UserID, UserID: userID})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	olhttp.NewAPIResponse(UserLoginHistoryResponse(app.MapSlice(result, func(entry app.LoginHistoryEntry) LoginHistoryEntryResponse {
		return LoginHistoryEntryResponse{IPAddress: entry.IPAddress, UserAgent: entry.UserAgent, LoggedInAt: entry.LoggedInAt}
	}))).Write(w)
}

func writeModerationActionResult(w http.ResponseWriter, err error) {
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	apiWriteOK(w)
}
