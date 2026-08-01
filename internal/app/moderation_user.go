package app

import (
	"context"
	"encoding/json"
	"math"
	"math/rand/v2"
	"time"

	"github.com/gofrs/uuid"
)

type ModerationUserInfo struct {
	ID              uuid.UUID
	Name            string
	Email           string
	About           string
	Avatar          string
	JoinedAt        time.Time
	Role            UserRole
	IsBanned        bool
	IsEmailVerified bool
	BooksTotal      int64
	CommentsTotal   int64
	FollowersTotal  int64
}

type GetModerationUserQuery struct {
	ActorUserID uuid.UUID
	UserID      uuid.UUID
}

type ModerationUserPageQuery struct {
	ActorUserID uuid.UUID
	UserID      uuid.UUID
	Page        uint32
	PageSize    uint32
}

type ModerationPage[T any] struct {
	Entries    []T
	Page       uint32
	PageSize   uint32
	Total      int64
	TotalPages uint32
}

type ModerationUserBook struct {
	ID                int64
	Name              string
	CreatedAt         time.Time
	IsPubliclyVisible bool
	IsBanned          bool
	IsTrashed         bool
}

type ModerationUserComment struct {
	ID, ChapterID, BookID          int64
	Content, ChapterName, BookName string
	CreatedAt                      time.Time
	Deleted                        bool
}

type ModerationUserHistoryEntry struct {
	ID            int64
	Time          time.Time
	Type          string
	Reason        string
	Payload       json.RawMessage
	ActorUserID   uuid.UUID
	ActorUserName string
}

type ModerationUserReport struct {
	Number, TargetID, Reason, Description, ReporterUserName string
	ID                                                      int64
	Time                                                    time.Time
	TargetType                                              ReportTargetType
	ReporterUserID                                          uuid.UUID
}

type ModerationUserRepository interface {
	GetUserInfo(context.Context, uuid.UUID) (ModerationUserInfo, error)
	GetBooks(context.Context, uuid.UUID, int32, int32) ([]ModerationUserBook, error)
	GetComments(context.Context, uuid.UUID, int32, int32) ([]ModerationUserComment, error)
	GetHistory(context.Context, uuid.UUID, int32, int32) ([]ModerationUserHistoryEntry, int64, error)
	GetReports(context.Context, uuid.UUID, int32, int32) ([]ModerationUserReport, int64, error)
	CreateTemporaryRandomReports(context.Context, uuid.UUID, int) error
}

type ModerationUserService interface {
	GetUserInfo(context.Context, GetModerationUserQuery) (ModerationUserInfo, error)
	GetBooks(context.Context, ModerationUserPageQuery) (ModerationPage[ModerationUserBook], error)
	GetComments(context.Context, ModerationUserPageQuery) (ModerationPage[ModerationUserComment], error)
	GetHistory(context.Context, ModerationUserPageQuery) (ModerationPage[ModerationUserHistoryEntry], error)
	GetReports(context.Context, ModerationUserPageQuery) (ModerationPage[ModerationUserReport], error)
}

func normalizeModerationPage(page, pageSize uint32) (uint32, uint32, int32, int32) {
	page = max(1, page)
	pageSize = min(max(1, pageSize), 100)
	return page, pageSize, int32(pageSize), int32((page - 1) * pageSize)
}

func moderationPage[T any](entries []T, page, pageSize uint32, total int64) ModerationPage[T] {
	return ModerationPage[T]{Entries: entries, Page: page, PageSize: pageSize, Total: total, TotalPages: uint32(math.Ceil(float64(total) / float64(pageSize)))}
}

func (s *moderationUserService) GetBooks(ctx context.Context, query ModerationUserPageQuery) (ModerationPage[ModerationUserBook], error) {
	if err := s.auth.AuthorizeModerator(ctx, query.ActorUserID); err != nil {
		return ModerationPage[ModerationUserBook]{}, err
	}
	info, err := s.repo.GetUserInfo(ctx, query.UserID)
	if err != nil {
		return ModerationPage[ModerationUserBook]{}, err
	}
	page, size, limit, offset := normalizeModerationPage(query.Page, query.PageSize)
	entries, err := s.repo.GetBooks(ctx, query.UserID, limit, offset)
	if err != nil {
		return ModerationPage[ModerationUserBook]{}, err
	}
	return moderationPage(entries, page, size, info.BooksTotal), nil
}

func (s *moderationUserService) GetComments(ctx context.Context, query ModerationUserPageQuery) (ModerationPage[ModerationUserComment], error) {
	if err := s.auth.AuthorizeModerator(ctx, query.ActorUserID); err != nil {
		return ModerationPage[ModerationUserComment]{}, err
	}
	info, err := s.repo.GetUserInfo(ctx, query.UserID)
	if err != nil {
		return ModerationPage[ModerationUserComment]{}, err
	}
	page, size, limit, offset := normalizeModerationPage(query.Page, query.PageSize)
	entries, err := s.repo.GetComments(ctx, query.UserID, limit, offset)
	if err != nil {
		return ModerationPage[ModerationUserComment]{}, err
	}
	return moderationPage(entries, page, size, info.CommentsTotal), nil
}

func (s *moderationUserService) GetHistory(ctx context.Context, query ModerationUserPageQuery) (ModerationPage[ModerationUserHistoryEntry], error) {
	if err := s.auth.AuthorizeModerator(ctx, query.ActorUserID); err != nil {
		return ModerationPage[ModerationUserHistoryEntry]{}, err
	}
	if _, err := s.repo.GetUserInfo(ctx, query.UserID); err != nil {
		return ModerationPage[ModerationUserHistoryEntry]{}, err
	}
	page, size, limit, offset := normalizeModerationPage(query.Page, query.PageSize)
	entries, total, err := s.repo.GetHistory(ctx, query.UserID, limit, offset)
	if err != nil {
		return ModerationPage[ModerationUserHistoryEntry]{}, err
	}
	return moderationPage(entries, page, size, total), nil
}

func (s *moderationUserService) GetReports(ctx context.Context, query ModerationUserPageQuery) (ModerationPage[ModerationUserReport], error) {
	if err := s.auth.AuthorizeModerator(ctx, query.ActorUserID); err != nil {
		return ModerationPage[ModerationUserReport]{}, err
	}
	if _, err := s.repo.GetUserInfo(ctx, query.UserID); err != nil {
		return ModerationPage[ModerationUserReport]{}, err
	}
	// TEMPORARY: Seed demo reports until real report creation is connected to controllers.
	_, total, err := s.repo.GetReports(ctx, query.UserID, 1, 0)
	if err != nil {
		return ModerationPage[ModerationUserReport]{}, err
	}
	if total < 50 {
		if err = s.repo.CreateTemporaryRandomReports(ctx, query.UserID, 4); err != nil {
			return ModerationPage[ModerationUserReport]{}, err
		}
	}
	page, size, limit, offset := normalizeModerationPage(query.Page, query.PageSize)
	entries, total, err := s.repo.GetReports(ctx, query.UserID, limit, offset)
	if err != nil {
		return ModerationPage[ModerationUserReport]{}, err
	}
	return moderationPage(entries, page, size, total), nil
}

var temporaryReportReasons = []string{"Spam", "Harassment", "Inappropriate content", "Impersonation", "Other"}
var temporaryReportDescriptions = []string{"Automatically generated temporary report.", "Flagged for moderator review.", "Multiple users reported this activity.", "Potential policy violation."}

func temporaryRandomReportValues() (string, string) {
	return temporaryReportReasons[rand.IntN(len(temporaryReportReasons))], temporaryReportDescriptions[rand.IntN(len(temporaryReportDescriptions))]
}

type moderationUserService struct {
	auth ModerationAuthorizer
	repo ModerationUserRepository
}

func (s *moderationUserService) GetUserInfo(ctx context.Context, query GetModerationUserQuery) (ModerationUserInfo, error) {
	if err := s.auth.AuthorizeModerator(ctx, query.ActorUserID); err != nil {
		return ModerationUserInfo{}, err
	}
	return s.repo.GetUserInfo(ctx, query.UserID)
}

func NewModerationUserService(auth ModerationAuthorizer, repo ModerationUserRepository) ModerationUserService {
	return &moderationUserService{auth: auth, repo: repo}
}
