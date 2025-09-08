package app

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
)

var (
	ModerationBookErrors                = AppErrors.NewSubNamespace("mod_books")
	InvalidModerationActionError        = ModerationBookErrors.NewType("invalid_mod_action")
	ErrInvalidModerationAction_NoReason = InvalidModerationActionError.New("no reason provided")
	ModerationBookNotFoundError         = ModerationBookErrors.NewType("book_404", ErrTraitEntityNotFound)
)

type ModerationPerformBookActionCommand struct {
	Reason      string
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

type BookModerationInfo struct {
	ID             int64
	Name           string
	Summary        string
	IsBanned       bool
	IsShadowBanned bool
	IsPermDeleted  bool
}

type BookModerationLog struct {
	Time          time.Time
	Action        store.BookActionType
	Payload       json.RawMessage
	Reason        string
	ActorUserID   uuid.UUID
	ActorUserName string
}

type GetBookLogQuery struct {
	Page        uint32
	PageSize    uint32
	OfTypes     []store.BookActionType
	BookID      int64
	ActorUserID uuid.UUID
}

type BookLogResult struct {
	Entries         []BookModerationLog
	Page            int32
	PageSize        int32
	HasNextPage     bool
	HasPreviousPage bool
	TotalPages      uint32
}

type ModerationBookService interface {
	GetBookInfo(ctx context.Context, query GetBookInfoQuery) (BookModerationInfo, error)
	GetBookLog(ctx context.Context, query GetBookLogQuery) (BookLogResult, error)
	BanBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error
	ShadowBanBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error
	PermanentlyRemoveBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error

	UnBanBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error
	UnShadowBanBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error
}
