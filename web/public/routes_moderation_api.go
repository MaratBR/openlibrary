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
	reports      app.ModerationReportService
	audit        app.ModerationAuditLogService
}

func newAPIModerationController(books app.ModerationBookService, content app.ContentModerationService, loginHistory app.LoginHistoryService, users app.ModerationUserService, reports app.ModerationReportService, audit app.ModerationAuditLogService) *apiControllerModeration {
	return &apiControllerModeration{books: books, content: content, loginHistory: loginHistory, users: users, reports: reports, audit: audit}
}

func (c *apiControllerModeration) Register(r chi.Router) {
	r.Route("/moderation", func(r chi.Router) {
		r.Use(apiRequiresAuthorizationMiddleware)
		r.Get("/books/{bookID}", c.getBook)
		r.Get("/books/{bookID}/chapters", c.getBookChapters)
		r.Get("/books/{bookID}/log", c.getBookLog)
		r.Post("/books/{bookID}/actions/{action}", c.bookAction)
		r.Post("/chapters/{chapterID}/actions/{action}", c.chapterAction)
		r.Post("/comments/{commentID}/actions/{action}", c.commentAction)
		r.Get("/reports", c.searchReports)
		r.Get("/reports/{reportID}", c.getReport)
		r.Get("/audit-log", c.getAuditLog)
		r.Get("/users", c.searchUsers)
		r.Get("/users/{userID}", c.getUser)
		r.Get("/users/{userID}/books", c.getUserBooks)
		r.Get("/users/{userID}/comments", c.getUserComments)
		r.Get("/users/{userID}/history", c.getUserHistory)
		r.Get("/users/{userID}/reports", c.getUserReports)
		r.Post("/users/{userID}/actions/{action}", c.userAction)
		r.Get("/users/{userID}/login-history", c.getLoginHistory)
		r.Get("/users/{userID}/login-locations", c.getLoginLocations)
	})
}

// go2tsdef:generate
type ModerationReasonRequest struct {
	Reason string `json:"reason"`
}

// go2tsdef:generate
type ModerationReportActivityResponse struct {
	Time        time.Time `json:"time"`
	Actor       string    `json:"actor"`
	Description string    `json:"description"`
	Kind        string    `json:"kind"`
}

// go2tsdef:generate
type ModerationReportDetailResponse struct {
	ID               string                               `json:"id"`
	Number           string                               `json:"number"`
	Time             time.Time                            `json:"time"`
	ReporterUserID   string                               `json:"reporterUserId"`
	ReporterUserName string                               `json:"reporterUserName"`
	TargetType       string                               `json:"targetType"`
	TargetID         string                               `json:"targetId"`
	Reason           string                               `json:"reason"`
	Description      string                               `json:"description"`
	Status           string                               `json:"status"`
	Priority         string                               `json:"priority"`
	Activities       []ModerationReportActivityResponse   `json:"activities"`
	BookContext      *ModerationReportBookContextResponse `json:"bookContext" go2tsdef:"ModerationReportBookContextResponse | null"`
}

// go2tsdef:generate
type ModerationReportBookContextResponse struct {
	Scope              string    `json:"scope"`
	Title              string    `json:"title"`
	Author             string    `json:"author"`
	CoverURL           string    `json:"coverUrl"`
	Chapter            string    `json:"chapter"`
	Excerpt            string    `json:"excerpt"`
	Rating             string    `json:"rating"`
	Warnings           []string  `json:"warnings"`
	PublicationState   string    `json:"publicationState"`
	LastUpdated        time.Time `json:"lastUpdated"`
	RelatedReports     int       `json:"relatedReports"`
	EditedAfterMinutes int64     `json:"editedAfterMinutes"`
}

// go2tsdef:generate
type ModerationReportListEntryResponse struct {
	ID               string    `json:"id"`
	Number           string    `json:"number"`
	Time             time.Time `json:"time"`
	ReporterUserID   string    `json:"reporterUserId"`
	ReporterUserName string    `json:"reporterUserName"`
	TargetType       string    `json:"targetType"`
	TargetID         string    `json:"targetId"`
	Reason           string    `json:"reason"`
	Description      string    `json:"description"`
}

// go2tsdef:generate
type ModerationReportsSearchResponse struct {
	Entries    []ModerationReportListEntryResponse `json:"entries"`
	Page       uint32                              `json:"page"`
	PageSize   uint32                              `json:"pageSize"`
	TotalPages uint32                              `json:"totalPages"`
	Total      int64                               `json:"total"`
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
	ID                   int64                                `json:"id,string"`
	Name                 string                               `json:"name"`
	Summary              string                               `json:"summary"`
	IsBanned             bool                                 `json:"isBanned"`
	IsShadowBanned       bool                                 `json:"isShadowBanned"`
	IsPermanentlyRemoved bool                                 `json:"isPermanentlyRemoved"`
	AuthorUserID         string                               `json:"authorUserId"`
	AuthorUserName       string                               `json:"authorUserName"`
	CreatedAt            time.Time                            `json:"createdAt"`
	AgeRating            string                               `json:"ageRating"`
	IsPubliclyVisible    bool                                 `json:"isPubliclyVisible"`
	Words                int32                                `json:"words"`
	Chapters             int32                                `json:"chapters"`
	ReportsCount         int64                                `json:"reportsCount"`
	LatestPendingReport  *ModerationBookPendingReportResponse `json:"latestPendingReport" go2tsdef:"ModerationBookPendingReportResponse | null"`
	BanReason            string                               `json:"banReason"`
}

// go2tsdef:generate
type ModerationBookPendingReportResponse struct {
	ID     string    `json:"id"`
	Number string    `json:"number"`
	Reason string    `json:"reason"`
	Time   time.Time `json:"time"`
}

// go2tsdef:generate
type ModerationBookChapterResponse struct {
	ID                string                  `json:"id"`
	Name              string                  `json:"name"`
	CreatedAt         time.Time               `json:"createdAt"`
	UpdatedAt         app.Nullable[time.Time] `json:"updatedAt"`
	Words             int32                   `json:"words"`
	IsPubliclyVisible bool                    `json:"isPubliclyVisible"`
	HasPendingReports bool                    `json:"hasPendingReports"`
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
type LoginLocationResponse struct {
	Country    string    `json:"country"`
	Region     string    `json:"region"`
	City       string    `json:"city"`
	LastSeenAt time.Time `json:"lastSeenAt"`
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
type ModerationUserListEntryResponse struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Avatar      string                  `json:"avatar"`
	JoinedAt    time.Time               `json:"joinedAt"`
	Role        string                  `json:"role"`
	IsBanned    bool                    `json:"isBanned"`
	LastVisitAt app.Nullable[time.Time] `json:"lastVisitAt"`
	BannedAt    app.Nullable[time.Time] `json:"bannedAt"`
	BanReason   string                  `json:"banReason"`
}

// go2tsdef:generate
type ModerationUsersPageResponse struct {
	Entries    []ModerationUserListEntryResponse `json:"entries"`
	Page       uint32                            `json:"page"`
	PageSize   uint32                            `json:"pageSize"`
	TotalPages uint32                            `json:"totalPages"`
	Total      int64                             `json:"total"`
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
	ID            string          `json:"id"`
	Time          time.Time       `json:"time"`
	Type          string          `json:"type"`
	Reason        string          `json:"reason"`
	ActorUserID   string          `json:"actorUserId"`
	ActorUserName string          `json:"actorUserName"`
	TargetType    string          `json:"targetType"`
	TargetID      string          `json:"targetId"`
	Payload       json.RawMessage `json:"payload" go2tsdef:"unknown"`
}

// go2tsdef:generate
type ModerationAuditLogEntryResponse struct {
	ID            string          `json:"id"`
	Time          time.Time       `json:"time"`
	Action        string          `json:"action"`
	TargetType    string          `json:"targetType"`
	TargetID      string          `json:"targetId"`
	Reason        string          `json:"reason"`
	Payload       json.RawMessage `json:"payload" go2tsdef:"unknown"`
	ActorUserID   string          `json:"actorUserId"`
	ActorUserName string          `json:"actorUserName"`
}

// go2tsdef:generate
type ModerationAuditLogPageResponse struct {
	Entries    []ModerationAuditLogEntryResponse `json:"entries"`
	Page       uint32                            `json:"page"`
	PageSize   uint32                            `json:"pageSize"`
	TotalPages uint32                            `json:"totalPages"`
	Total      int64                             `json:"total"`
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

func (c *apiControllerModeration) getReport(w http.ResponseWriter, r *http.Request) {
	reportID, err := olhttp.URLParamInt64(r, "reportID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	result, err := c.reports.GetReport(r.Context(), app.GetModerationReportQuery{ActorUserID: auth.RequireSession(r.Context()).UserID, ReportID: reportID})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	activities := app.MapSlice(result.Activities, func(activity app.ModerationReportActivity) ModerationReportActivityResponse {
		return ModerationReportActivityResponse{Time: activity.Time, Actor: activity.Actor, Description: activity.Description, Kind: activity.Kind}
	})
	var bookContext *ModerationReportBookContextResponse
	if result.BookContext != nil {
		bookContext = &ModerationReportBookContextResponse{
			Scope: result.BookContext.Scope, Title: result.BookContext.Title, Author: result.BookContext.Author, CoverURL: result.BookContext.CoverURL,
			Chapter: result.BookContext.Chapter, Excerpt: result.BookContext.Excerpt, Rating: result.BookContext.Rating,
			Warnings: result.BookContext.Warnings, PublicationState: result.BookContext.PublicationState,
			LastUpdated: result.BookContext.LastUpdated, RelatedReports: result.BookContext.RelatedReports,
			EditedAfterMinutes: int64(result.BookContext.EditedAfter / time.Minute),
		}
	}
	olhttp.NewAPIResponse(ModerationReportDetailResponse{
		ID: strconv.FormatInt(result.ID, 10), Number: result.Number, Time: result.Time,
		ReporterUserID: result.ReporterUserID.String(), ReporterUserName: result.ReporterUserName,
		TargetType: string(result.TargetType), TargetID: result.TargetID, Reason: result.Reason, Description: result.Description,
		Status: result.Status, Priority: result.Priority, Activities: activities, BookContext: bookContext,
	}).Write(w)
}

func (c *apiControllerModeration) searchReports(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	result, err := c.reports.SearchReports(r.Context(), app.SearchModerationReportsQuery{
		ActorUserID: auth.RequireSession(r.Context()).UserID,
		Search:      params.Get("search"), TargetType: params.Get("targetType"),
		Page: olhttp.GetPage(params, "page"), PageSize: olhttp.GetPageSize(params, "pageSize", 1, 100, 20),
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	entries := app.MapSlice(result.Entries, func(entry app.ModerationReportListEntry) ModerationReportListEntryResponse {
		return ModerationReportListEntryResponse{ID: strconv.FormatInt(entry.ID, 10), Number: entry.Number, Time: entry.Time, ReporterUserID: entry.ReporterUserID.String(), ReporterUserName: entry.ReporterUserName, TargetType: string(entry.TargetType), TargetID: entry.TargetID, Reason: entry.Reason, Description: entry.Description}
	})
	olhttp.NewAPIResponse(ModerationReportsSearchResponse{Entries: entries, Page: result.Page, PageSize: result.PageSize, Total: result.Total, TotalPages: result.TotalPages}).Write(w)
}

func (c *apiControllerModeration) searchUsers(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	result, err := c.users.SearchUsers(r.Context(), app.ModerationUsersQuery{
		ActorUserID: auth.RequireSession(r.Context()).UserID,
		Search:      params.Get("search"), Banned: params.Get("banned"), Role: params.Get("role"),
		Page: olhttp.GetPage(params, "page"), PageSize: olhttp.GetPageSize(params, "pageSize", 1, 100, 20),
	})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	entries := app.MapSlice(result.Entries, func(entry app.ModerationUserListEntry) ModerationUserListEntryResponse {
		return ModerationUserListEntryResponse{ID: entry.ID.String(), Name: entry.Name, Avatar: entry.Avatar, JoinedAt: entry.JoinedAt, Role: string(entry.Role), IsBanned: entry.IsBanned, LastVisitAt: app.NullableFromPtr(entry.LastVisitAt), BannedAt: app.NullableFromPtr(entry.BannedAt), BanReason: entry.BanReason}
	})
	olhttp.NewAPIResponse(ModerationUsersPageResponse{Entries: entries, Page: result.Page, PageSize: result.PageSize, Total: result.Total, TotalPages: result.TotalPages}).Write(w)
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
		return ModerationUserHistoryResponse{ID: strconv.FormatInt(entry.ID, 10), Time: entry.Time, Type: entry.Type, TargetType: entry.TargetType, TargetID: entry.TargetID, Reason: entry.Reason, Payload: entry.Payload, ActorUserID: entry.ActorUserID.String(), ActorUserName: entry.ActorUserName}
	})
	olhttp.NewAPIResponse(ModerationUserHistoryPageResponse{Entries: entries, Page: result.Page, PageSize: result.PageSize, Total: result.Total, TotalPages: result.TotalPages}).Write(w)
}

func (c *apiControllerModeration) getAuditLog(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	result, err := c.audit.GetAuditLog(r.Context(), app.ModerationAuditLogQuery{ActorUserID: auth.RequireSession(r.Context()).UserID, TargetType: params.Get("targetType"), Page: olhttp.GetPage(params, "page"), PageSize: olhttp.GetPageSize(params, "pageSize", 1, 100, 25)})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	entries := app.MapSlice(result.Entries, func(entry app.ModerationAuditLogEntry) ModerationAuditLogEntryResponse {
		return ModerationAuditLogEntryResponse{ID: strconv.FormatInt(entry.ID, 10), Time: entry.Time, Action: entry.Action, TargetType: entry.TargetType, TargetID: entry.TargetID, Reason: entry.Reason, Payload: entry.Payload, ActorUserID: entry.ActorUserID.String(), ActorUserName: entry.ActorUserName}
	})
	olhttp.NewAPIResponse(ModerationAuditLogPageResponse{Entries: entries, Page: result.Page, PageSize: result.PageSize, TotalPages: result.TotalPages, Total: result.Total}).Write(w)
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
		AuthorUserID:         result.AuthorUserID.String(), AuthorUserName: result.AuthorUserName, CreatedAt: result.CreatedAt,
		AgeRating: result.AgeRating, IsPubliclyVisible: result.IsPubliclyVisible, Words: result.Words, Chapters: result.Chapters,
		ReportsCount: result.ReportsCount, BanReason: result.BanReason, LatestPendingReport: mapPendingBookReport(result.LatestPendingReport),
	}).Write(w)
}

func mapPendingBookReport(report *app.BookPendingReport) *ModerationBookPendingReportResponse {
	if report == nil {
		return nil
	}
	return &ModerationBookPendingReportResponse{ID: strconv.FormatInt(report.ID, 10), Number: report.Number, Reason: report.Reason, Time: report.Time}
}

func (c *apiControllerModeration) getBookChapters(w http.ResponseWriter, r *http.Request) {
	bookID, err := olhttp.URLParamInt64(r, "bookID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	rows, err := c.books.GetBookChapters(r.Context(), app.GetBookInfoQuery{ActorUserID: auth.RequireSession(r.Context()).UserID, BookID: bookID})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	response := app.MapSlice(rows, func(row app.BookModerationChapter) ModerationBookChapterResponse {
		return ModerationBookChapterResponse{ID: strconv.FormatInt(row.ID, 10), Name: row.Name, CreatedAt: row.CreatedAt, UpdatedAt: app.NullableFromPtr(row.UpdatedAt), Words: row.Words, IsPubliclyVisible: row.IsPubliclyVisible, HasPendingReports: row.HasPendingReports}
	})
	olhttp.NewAPIResponse(response).Write(w)
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
	var input ModerationValueRequest
	if err = decodeModerationJSON(w, r, &input); err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	cmd := app.ModerationPerformBookActionCommand{ActorUserID: auth.RequireSession(r.Context()).UserID, BookID: bookID, Reason: input.Reason}
	cmd.Value = input.Value
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
	case "change-age-rating":
		err = c.books.ChangeAgeRating(r.Context(), cmd)
	case "change-summary":
		err = c.books.ChangeSummary(r.Context(), cmd)
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

func (c *apiControllerModeration) getLoginLocations(w http.ResponseWriter, r *http.Request) {
	userID, err := olhttp.URLParamUUID(r, "userID")
	if err != nil {
		apiWriteBadRequest(w, err)
		return
	}
	result, err := c.loginHistory.GetRecentLoginLocations(r.Context(), app.GetLoginHistoryQuery{ActorUserID: auth.RequireSession(r.Context()).UserID, UserID: userID})
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	response := app.MapSlice(result, func(location app.LoginLocation) LoginLocationResponse {
		return LoginLocationResponse{Country: location.Country, Region: location.Region, City: location.City, LastSeenAt: location.LastSeenAt}
	})
	olhttp.NewAPIResponse(response).Write(w)
}

func writeModerationActionResult(w http.ResponseWriter, err error) {
	if err != nil {
		apiWriteApplicationError(w, err)
		return
	}
	apiWriteOK(w)
}
