package public

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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
	users        app.ModerationUserService
}

func newAPIModerationController(books app.ModerationBookService, content app.ContentModerationService, loginHistory app.LoginHistoryService, users app.ModerationUserService) *apiControllerModeration {
	return &apiControllerModeration{books: books, content: content, loginHistory: loginHistory, users: users}
}

func (c *apiControllerModeration) Register(r chi.Router) {
	r.Route("/moderation", func(r chi.Router) {
		r.Use(apiRequiresAuthorizationMiddleware)
		r.Get("/books/{bookID}", c.getBook)
		r.Get("/books/{bookID}/log", c.getBookLog)
		r.Post("/books/{bookID}/actions/{action}", c.bookAction)
		r.Post("/chapters/{chapterID}/actions/{action}", c.chapterAction)
		r.Post("/comments/{commentID}/actions/{action}", c.commentAction)
		r.Get("/users/{userID}", c.getUser)
		r.Get("/users/{userID}/books", c.getUserBooks)
		r.Get("/users/{userID}/comments", c.getUserComments)
		r.Get("/users/{userID}/history", c.getUserHistory)
		r.Get("/users/{userID}/reports", c.getUserReports)
		r.Post("/users/{userID}/actions/{action}", c.userAction)
		r.Get("/users/{userID}/login-history", c.getLoginHistory)
	})
}

// go2tsdef:generate
type ModerationReasonRequest struct {
	Reason string `json:"reason"`
}

// go2tsdef:generate
type ModerationBanRequest struct {
	Reason string    `json:"reason"`
	Until  time.Time `json:"until"`
}

// go2tsdef:generate
type ModerationValueRequest struct {
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
type UserLoginHistoryResponse struct {
	Entries    []LoginHistoryEntryResponse `json:"entries"`
	Page       uint32                      `json:"page"`
	PageSize   uint32                      `json:"pageSize"`
	Total      int                         `json:"total"`
	TotalPages uint32                      `json:"totalPages"`
}

// go2tsdef:generate
type ModerationUserResponse struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	About           string    `json:"about"`
	Avatar          string    `json:"avatar"`
	JoinedAt        time.Time `json:"joinedAt"`
	Role            string    `json:"role"`
	IsBanned        bool      `json:"isBanned"`
	IsEmailVerified bool      `json:"isEmailVerified"`
	BooksTotal      int64     `json:"booksTotal"`
	CommentsTotal   int64     `json:"commentsTotal"`
	FollowersTotal  int64     `json:"followersTotal"`
}

// go2tsdef:generate
type ModerationUserBookResponse struct {
	ID                int64     `json:"id,string"`
	Name              string    `json:"name"`
	CreatedAt         time.Time `json:"createdAt"`
	IsPubliclyVisible bool      `json:"isPubliclyVisible"`
	IsBanned          bool      `json:"isBanned"`
	IsTrashed         bool      `json:"isTrashed"`
}

// go2tsdef:generate
type ModerationUserCommentResponse struct {
	ID          string    `json:"id"`
	ChapterID   string    `json:"chapterId"`
	BookID      string    `json:"bookId"`
	Content     string    `json:"content"`
	ChapterName string    `json:"chapterName"`
	BookName    string    `json:"bookName"`
	CreatedAt   time.Time `json:"createdAt"`
	Deleted     bool      `json:"deleted"`
}

// go2tsdef:generate
type ModerationUserHistoryResponse struct {
	ID            string    `json:"id"`
	Time          time.Time `json:"time"`
	Type          string    `json:"type"`
	Reason        string    `json:"reason"`
	ActorUserID   string    `json:"actorUserId"`
	ActorUserName string    `json:"actorUserName"`
}

// go2tsdef:generate
type ModerationUserReportResponse struct {
	ID               string    `json:"id"`
	Number           string    `json:"number"`
	Time             time.Time `json:"time"`
	TargetType       string    `json:"targetType"`
	TargetID         string    `json:"targetId"`
	Reason           string    `json:"reason"`
	Description      string    `json:"description"`
	ReporterUserID   string    `json:"reporterUserId"`
	ReporterUserName string    `json:"reporterUserName"`
}

// go2tsdef:generate
type ModerationUserBooksPageResponse struct {
	Entries    []ModerationUserBookResponse `json:"entries"`
	Page       uint32                       `json:"page"`
	PageSize   uint32                       `json:"pageSize"`
	TotalPages uint32                       `json:"totalPages"`
	Total      int64                        `json:"total"`
}

// go2tsdef:generate
type ModerationUserCommentsPageResponse struct {
	Entries    []ModerationUserCommentResponse `json:"entries"`
	Page       uint32                          `json:"page"`
	PageSize   uint32                          `json:"pageSize"`
	TotalPages uint32                          `json:"totalPages"`
	Total      int64                           `json:"total"`
}

// go2tsdef:generate
type ModerationUserHistoryPageResponse struct {
	Entries    []ModerationUserHistoryResponse `json:"entries"`
	Page       uint32                          `json:"page"`
	PageSize   uint32                          `json:"pageSize"`
	TotalPages uint32                          `json:"totalPages"`
	Total      int64                           `json:"total"`
}

// go2tsdef:generate
type ModerationUserReportsPageResponse struct {
	Entries    []ModerationUserReportResponse `json:"entries"`
	Page       uint32                         `json:"page"`
	PageSize   uint32                         `json:"pageSize"`
	TotalPages uint32                         `json:"totalPages"`
	Total      int64                          `json:"total"`
}

func decodeModerationJSON(w http.ResponseWriter, r *http.Request, value any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64*1024))
	decoder.DisallowUnknownFields()
	return decoder.Decode(value)
}

func (c *apiControllerModeration) getUser(w http.ResponseWriter, r *http.Request) {
	userID, err := olhttp.URLParamUUID(r, "userID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	result, err := c.users.GetUserInfo(r.Context(), app.GetModerationUserQuery{
		ActorUserID: auth.RequireSession(r.Context()).UserID,
		UserID:      userID,
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	olhttp.NewAPIResponse(ModerationUserResponse{
		ID: result.ID.String(), Name: result.Name, Email: result.Email, About: result.About, Avatar: result.Avatar,
		JoinedAt: result.JoinedAt, Role: string(result.Role), IsBanned: result.IsBanned,
		IsEmailVerified: result.IsEmailVerified, BooksTotal: result.BooksTotal,
		CommentsTotal: result.CommentsTotal, FollowersTotal: result.FollowersTotal,
	}).Write(w)
}

func moderationUserPageQuery(r *http.Request) (app.ModerationUserPageQuery, error) {
	userID, err := olhttp.URLParamUUID(r, "userID")
	if err != nil {
		return app.ModerationUserPageQuery{}, err
	}
	return app.ModerationUserPageQuery{ActorUserID: auth.RequireSession(r.Context()).UserID, UserID: userID, Page: olhttp.GetPage(r.URL.Query(), "page"), PageSize: olhttp.GetPageSize(r.URL.Query(), "pageSize", 1, 100, 20)}, nil
}

func (c *apiControllerModeration) getUserBooks(w http.ResponseWriter, r *http.Request) {
	query, err := moderationUserPageQuery(r)
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	result, err := c.users.GetBooks(r.Context(), query)
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	entries := app.MapSlice(result.Entries, func(entry app.ModerationUserBook) ModerationUserBookResponse {
		return ModerationUserBookResponse{ID: entry.ID, Name: entry.Name, CreatedAt: entry.CreatedAt, IsPubliclyVisible: entry.IsPubliclyVisible, IsBanned: entry.IsBanned, IsTrashed: entry.IsTrashed}
	})
	olhttp.NewAPIResponse(ModerationUserBooksPageResponse{Entries: entries, Page: result.Page, PageSize: result.PageSize, Total: result.Total, TotalPages: result.TotalPages}).Write(w)
}

func (c *apiControllerModeration) getUserComments(w http.ResponseWriter, r *http.Request) {
	query, err := moderationUserPageQuery(r)
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	result, err := c.users.GetComments(r.Context(), query)
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	entries := app.MapSlice(result.Entries, func(entry app.ModerationUserComment) ModerationUserCommentResponse {
		return ModerationUserCommentResponse{ID: strconv.FormatInt(entry.ID, 10), ChapterID: strconv.FormatInt(entry.ChapterID, 10), BookID: strconv.FormatInt(entry.BookID, 10), Content: entry.Content, ChapterName: entry.ChapterName, BookName: entry.BookName, CreatedAt: entry.CreatedAt, Deleted: entry.Deleted}
	})
	olhttp.NewAPIResponse(ModerationUserCommentsPageResponse{Entries: entries, Page: result.Page, PageSize: result.PageSize, Total: result.Total, TotalPages: result.TotalPages}).Write(w)
}

func (c *apiControllerModeration) getUserHistory(w http.ResponseWriter, r *http.Request) {
	query, err := moderationUserPageQuery(r)
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	result, err := c.users.GetHistory(r.Context(), query)
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	entries := app.MapSlice(result.Entries, func(entry app.ModerationUserHistoryEntry) ModerationUserHistoryResponse {
		return ModerationUserHistoryResponse{ID: strconv.FormatInt(entry.ID, 10), Time: entry.Time, Type: entry.Type, Reason: entry.Reason, ActorUserID: entry.ActorUserID.String(), ActorUserName: entry.ActorUserName}
	})
	olhttp.NewAPIResponse(ModerationUserHistoryPageResponse{Entries: entries, Page: result.Page, PageSize: result.PageSize, Total: result.Total, TotalPages: result.TotalPages}).Write(w)
}

func (c *apiControllerModeration) getUserReports(w http.ResponseWriter, r *http.Request) {
	query, err := moderationUserPageQuery(r)
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	result, err := c.users.GetReports(r.Context(), query)
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	entries := app.MapSlice(result.Entries, func(entry app.ModerationUserReport) ModerationUserReportResponse {
		return ModerationUserReportResponse{ID: strconv.FormatInt(entry.ID, 10), Number: entry.Number, Time: entry.Time, TargetType: string(entry.TargetType), TargetID: entry.TargetID, Reason: entry.Reason, Description: entry.Description, ReporterUserID: entry.ReporterUserID.String(), ReporterUserName: entry.ReporterUserName}
	})
	olhttp.NewAPIResponse(ModerationUserReportsPageResponse{Entries: entries, Page: result.Page, PageSize: result.PageSize, Total: result.Total, TotalPages: result.TotalPages}).Write(w)
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
	var input ModerationReasonRequest
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
	var input ModerationReasonRequest
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
	var input ModerationReasonRequest
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
		var input ModerationBanRequest
		if err = decodeModerationJSON(w, r, &input); err != nil {
			apiWriteBadRequest(w, err)
			return
		}
		cmd := app.BanUserCommand{ActorUserID: actor, UserID: userID, Until: input.Until, Reason: input.Reason}
		err = c.content.BanUser(r.Context(), cmd)
	case "permanent-ban", "unban":
		var input ModerationReasonRequest
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
		var input ModerationValueRequest
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
	entries := app.MapSlice(result, func(entry app.LoginHistoryEntry) LoginHistoryEntryResponse {
		return LoginHistoryEntryResponse{IPAddress: entry.IPAddress, UserAgent: entry.UserAgent, LoggedInAt: entry.LoggedInAt}
	})
	page := olhttp.GetPage(r.URL.Query(), "page")
	pageSize := olhttp.GetPageSize(r.URL.Query(), "pageSize", 1, 100, 20)
	start := min(len(entries), int((page-1)*pageSize))
	end := min(len(entries), start+int(pageSize))
	totalPages := uint32(0)
	if len(entries) > 0 {
		totalPages = uint32((len(entries) + int(pageSize) - 1) / int(pageSize))
	}
	olhttp.NewAPIResponse(UserLoginHistoryResponse{Entries: entries[start:end], Page: page, PageSize: pageSize, Total: len(entries), TotalPages: totalPages}).Write(w)
}

func writeModerationActionResult(w http.ResponseWriter, err error) {
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	apiWriteOK(w)
}
