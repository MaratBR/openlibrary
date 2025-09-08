package app

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type moderationBookService struct {
	db DB
}

// GetBookInfo implements ModerationBookService.
func (m *moderationBookService) GetBookInfo(ctx context.Context, query GetBookInfoQuery) (BookModerationInfo, error) {
	queries := store.New(m.db)

	// TODO: authorization

	row, err := queries.ModGetBookInfo(ctx, query.BookID)
	if err != nil {
		if err == store.ErrNoRows {
			return BookModerationInfo{}, ModerationBookNotFoundError.New(fmt.Sprintf("book with ID %d could not be found", query.BookID))
		} else {
			return BookModerationInfo{}, wrapUnexpectedDBError(err)
		}
	}

	return BookModerationInfo{
		IsBanned:       row.IsBanned,
		IsShadowBanned: row.IsShadowBanned,
		IsPermDeleted:  row.IsPermRemoved,
		Name:           row.Name,
		Summary:        row.Summary,
		ID:             query.BookID,
	}, nil
}

// GetBookLog implements ModerationBookService.
func (m *moderationBookService) GetBookLog(ctx context.Context, query GetBookLogQuery) (BookLogResult, error) {
	queries := store.New(m.db)

	var (
		page     int32
		pageSize int32
	)

	pageSize = int32(query.PageSize)
	page = int32(query.Page)

	if pageSize < 1 {
		pageSize = 1
	} else if pageSize > 1000 {
		pageSize = 1000
	}

	if page < 1 {
		page = 1
	} else if page > 1000 {
		page = 1000
	}

	rows, err := queries.ModGetBookLogFiltered(ctx, store.ModGetBookLogFilteredParams{
		BookID:      query.BookID,
		ActionTypes: query.OfTypes,
		Limit:       pageSize + 1,
		Offset:      (page - 1) * pageSize,
	})
	if err != nil {
		return BookLogResult{}, wrapUnexpectedDBError(err)
	}

	count, err := queries.ModCountBookLogFiltered(ctx, store.ModCountBookLogFilteredParams{
		BookID:      query.BookID,
		ActionTypes: query.OfTypes,
	})
	if err != nil {
		return BookLogResult{}, wrapUnexpectedDBError(err)
	}

	var hasNextPage bool

	if len(rows) == int(pageSize)+1 {
		rows = rows[:len(rows)-1]
		hasNextPage = true
	}

	result := BookLogResult{
		Page:            page,
		PageSize:        pageSize,
		HasPreviousPage: page > 1,
		HasNextPage:     hasNextPage,
		TotalPages:      uint32(math.Ceil(float64(count) / float64(pageSize))),
	}

	result.Entries = make([]BookModerationLog, 0, len(rows))

	for _, row := range rows {
		result.Entries = append(result.Entries, BookModerationLog{
			Time:          row.Time.Time,
			Action:        row.ActionType,
			Payload:       row.Payload,
			Reason:        row.Reason,
			ActorUserID:   uuidDbToDomain(row.ActorUserID),
			ActorUserName: row.ActorUserName,
		})
	}

	return result, nil
}

// BanBook implements ModerationBookService.
func (m *moderationBookService) BanBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	tx, err := m.db.Begin(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	queries := store.New(m.db).WithTx(tx)

	err = queries.ModSetBookBanned(ctx, store.ModSetBookBannedParams{
		ID:       cmd.BookID,
		IsBanned: true,
	})
	if err != nil {
		rollbackTx(ctx, tx)
		return wrapUnexpectedDBError(err)
	}

	err = m.addBookLog(ctx, queries, cmd.BookID, store.BookActionTypeBan, cmd.Reason, cmd.ActorUserID)
	if err != nil {
		rollbackTx(ctx, tx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	return nil
}

func (m *moderationBookService) addBookLog(ctx context.Context, queries *store.Queries, bookID int64, actionType store.BookActionType, reason string, actorUserID uuid.UUID) error {
	err := queries.ModAddBookLog(ctx, store.ModAddBookLogParams{
		ID:          GenID(),
		BookID:      bookID,
		Reason:      reason,
		Payload:     nil,
		Time:        pgtype.Timestamptz{Valid: true, Time: time.Now()},
		ActionType:  actionType,
		ActorUserID: uuidDomainToDb(actorUserID),
	})
	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	return nil
}

// PermanentlyRemoveBook implements ModerationBookService.
func (m *moderationBookService) PermanentlyRemoveBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	queries := store.New(m.db).WithTx(tx)

	err = queries.ModSetBookBanned(ctx, store.ModSetBookBannedParams{
		ID:       cmd.BookID,
		IsBanned: true,
	})
	if err != nil {
		rollbackTx(ctx, tx)
		return wrapUnexpectedDBError(err)
	}

	err = m.addBookLog(ctx, queries, cmd.BookID, store.BookActionTypeBan, cmd.Reason, cmd.ActorUserID)
	if err != nil {
		rollbackTx(ctx, tx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	return nil

}

// ShadowBanBook implements ModerationBookService.
func (m *moderationBookService) ShadowBanBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	tx, err := m.db.Begin(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	queries := store.New(m.db).WithTx(tx)

	err = queries.ModSetBookShadowBanned(ctx, store.ModSetBookShadowBannedParams{
		ID:             cmd.BookID,
		IsShadowBanned: true,
	})
	if err != nil {
		rollbackTx(ctx, tx)
		return wrapUnexpectedDBError(err)
	}

	err = m.addBookLog(ctx, queries, cmd.BookID, store.BookActionTypeShadowBan, cmd.Reason, cmd.ActorUserID)
	if err != nil {
		rollbackTx(ctx, tx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	return nil

}

// UnBanBook implements ModerationBookService.
func (m *moderationBookService) UnBanBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	tx, err := m.db.Begin(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	queries := store.New(m.db).WithTx(tx)

	err = queries.ModSetBookBanned(ctx, store.ModSetBookBannedParams{
		ID:       cmd.BookID,
		IsBanned: false,
	})
	if err != nil {
		rollbackTx(ctx, tx)
		return wrapUnexpectedDBError(err)
	}

	err = m.addBookLog(ctx, queries, cmd.BookID, store.BookActionTypeUnBan, cmd.Reason, cmd.ActorUserID)
	if err != nil {
		rollbackTx(ctx, tx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	return nil
}

// UnShadowBanBook implements ModerationBookService.
func (m *moderationBookService) UnShadowBanBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	tx, err := m.db.Begin(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	queries := store.New(m.db).WithTx(tx)

	err = queries.ModSetBookShadowBanned(ctx, store.ModSetBookShadowBannedParams{
		ID:             cmd.BookID,
		IsShadowBanned: false,
	})
	if err != nil {
		rollbackTx(ctx, tx)
		return wrapUnexpectedDBError(err)
	}

	err = m.addBookLog(ctx, queries, cmd.BookID, store.BookActionTypeUnShadowBan, cmd.Reason, cmd.ActorUserID)
	if err != nil {
		rollbackTx(ctx, tx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	return nil
}

func NewModerationBookService(db DB) ModerationBookService {
	return &moderationBookService{db: db}
}
