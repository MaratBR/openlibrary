package app

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/gofrs/uuid"
)

var (
	ModerationBookErrors                = apperror.AppErrors.NewSubNamespace("mod_books")
	InvalidModerationActionError        = ModerationBookErrors.NewType("invalid_mod_action")
	ErrInvalidModerationAction_NoReason = InvalidModerationActionError.New("no reason provided")
	ModerationBookNotFoundError         = ModerationBookErrors.NewType("book_404", apperror.ErrTraitEntityNotFound)
)

type ModerationPerformBookActionCommand struct {
	Reason      string
	Value       string
	ActorUserID uuid.UUID
	BookID      int64
}

func (c ModerationPerformBookActionCommand) Validate() error {
	if c.Reason == "" || strings.Trim(c.Reason, " \n\t") == "" {
		return ErrInvalidModerationAction_NoReason
	}

	return nil
}

type GetBookInfoQuery struct {
	ActorUserID uuid.UUID
	BookID      int64
}

type SearchModerationBooksQuery struct {
	ActorUserID                              uuid.UUID
	Search                                   string
	ExactName, IncludeBanned, IncludeDeleted bool
	Page, PageSize                           uint32
}

type ModerationBookListEntry struct {
	ID                                                                           int64
	Name                                                                         string
	CreatedAt                                                                    time.Time
	IsBanned, IsShadowBanned, IsTrashed, IsPermanentlyRemoved, IsPubliclyVisible bool
	Words, Chapters                                                              int32
	AuthorUserID                                                                 uuid.UUID
	AuthorUserName                                                               string
	ReportsCount                                                                 int64
}

type BookModerationInfo struct {
	ID                  int64
	Name                string
	Summary             string
	IsBanned            bool
	IsShadowBanned      bool
	IsPermDeleted       bool
	AuthorUserID        uuid.UUID
	AuthorUserName      string
	CreatedAt           time.Time
	AgeRating           string
	IsPubliclyVisible   bool
	Words               int32
	Chapters            int32
	ReportsCount        int64
	LatestPendingReport *BookPendingReport
	BanReason           string
}

type BookPendingReport struct {
	ID             int64
	Number, Reason string
	Time           time.Time
}
type BookModerationChapter struct {
	ID                                   int64
	Name                                 string
	CreatedAt                            time.Time
	UpdatedAt                            *time.Time
	Words                                int32
	IsPubliclyVisible, HasPendingReports bool
}

type BookModerationLog struct {
	Time          time.Time
	Action        BookActionType
	Payload       json.RawMessage
	Reason        string
	ActorUserID   uuid.UUID
	ActorUserName string
}

type GetBookLogQuery struct {
	Page        uint32
	PageSize    uint32
	OfTypes     []BookActionType
	BookID      int64
	ActorUserID uuid.UUID
}

type BookActionType string

const (
	BookActionTypeBan         BookActionType = "ban"
	BookActionTypeShadowBan   BookActionType = "shadow_ban"
	BookActionTypePermRemoval BookActionType = "perm_removal"
	BookActionTypeUnBan       BookActionType = "un_ban"
	BookActionTypeUnShadowBan BookActionType = "un_shadow_ban"
)

type BookLogResult struct {
	Entries         []BookModerationLog
	Page            int32
	PageSize        int32
	HasNextPage     bool
	HasPreviousPage bool
	TotalPages      uint32
}

type ModerationBookService interface {
	SearchBooks(ctx context.Context, query SearchModerationBooksQuery) (ModerationPage[ModerationBookListEntry], error)
	GetBookInfo(ctx context.Context, query GetBookInfoQuery) (BookModerationInfo, error)
	GetBookLog(ctx context.Context, query GetBookLogQuery) (BookLogResult, error)
	GetBookChapters(ctx context.Context, query GetBookInfoQuery) ([]BookModerationChapter, error)
	ChangeAgeRating(ctx context.Context, cmd ModerationPerformBookActionCommand) error
	ChangeSummary(ctx context.Context, cmd ModerationPerformBookActionCommand) error
	BanBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error
	ShadowBanBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error
	PermanentlyRemoveBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error

	UnBanBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error
	UnShadowBanBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error
}
