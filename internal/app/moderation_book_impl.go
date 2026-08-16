package app

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/app/dal"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type moderationBookService struct {
	db   DB
	auth ModerationAuthorizer
}

func (m *moderationBookService) authorize(ctx context.Context, actor uuid.UUID) error {
	return m.auth.AuthorizeModerator(ctx, actor)
}

func (m *moderationBookService) SearchBooks(ctx context.Context, query SearchModerationBooksQuery) (ModerationPage[ModerationBookListEntry], error) {
	if err := m.authorize(ctx, query.ActorUserID); err != nil {
		return ModerationPage[ModerationBookListEntry]{}, err
	}
	query.Search = strings.TrimSpace(query.Search)
	searchID := pgtype.Int8{}
	if id, err := strconv.ParseInt(query.Search, 10, 64); err == nil {
		searchID = pgtype.Int8{Int64: id, Valid: true}
	}
	page, size, limit, offset := normalizeModerationPage(query.Page, query.PageSize)
	q := store.New(m.db)
	params := store.ModSearchBooksParams{Search: query.Search, SearchID: searchID, ExactName: query.ExactName, IncludeBanned: query.IncludeBanned, IncludeDeleted: query.IncludeDeleted, PageLimit: limit, PageOffset: offset}
	rows, err := q.ModSearchBooks(ctx, params)
	if err != nil {
		return ModerationPage[ModerationBookListEntry]{}, apperror.WrapUnexpectedDBError(err)
	}
	total, err := q.ModCountBooks(ctx, store.ModCountBooksParams{Search: params.Search, SearchID: params.SearchID, ExactName: params.ExactName, IncludeBanned: params.IncludeBanned, IncludeDeleted: params.IncludeDeleted})
	if err != nil {
		return ModerationPage[ModerationBookListEntry]{}, apperror.WrapUnexpectedDBError(err)
	}
	entries := MapSlice(rows, func(row store.ModSearchBooksRow) ModerationBookListEntry {
		return ModerationBookListEntry{ID: row.ID, Name: row.Name, CreatedAt: timeDbToDomain(row.CreatedAt), IsBanned: row.IsBanned, IsShadowBanned: row.IsShadowBanned, IsTrashed: row.IsTrashed, IsPermanentlyRemoved: row.IsPermRemoved, IsPubliclyVisible: row.IsPubliclyVisible, Words: row.Words, Chapters: row.Chapters, AuthorUserID: uuidDbToDomain(row.AuthorUserID), AuthorUserName: row.AuthorUserName, ReportsCount: row.ReportsCount}
	})
	return moderationPage(entries, page, size, total), nil
}

// GetBookInfo implements ModerationBookService.
func (m *moderationBookService) GetBookInfo(ctx context.Context, query GetBookInfoQuery) (BookModerationInfo, error) {
	queries := store.New(m.db)

	if err := m.authorize(ctx, query.ActorUserID); err != nil {
		return BookModerationInfo{}, err
	}

	row, err := queries.ModGetBookInfo(ctx, query.BookID)
	if err != nil {
		if err == store.ErrNoRows {
			return BookModerationInfo{}, ModerationBookNotFoundError.New(fmt.Sprintf("book with ID %d could not be found", query.BookID))
		} else {
			return BookModerationInfo{}, apperror.WrapUnexpectedDBError(err)
		}
	}

	result := BookModerationInfo{
		IsBanned:       row.IsBanned,
		IsShadowBanned: row.IsShadowBanned,
		IsPermDeleted:  row.IsPermRemoved,
		Name:           row.Name,
		Summary:        row.Summary,
		ID:             query.BookID,
		AuthorUserID:   uuidDbToDomain(row.AuthorUserID), AuthorUserName: row.AuthorUserName,
		CreatedAt: row.CreatedAt.Time, AgeRating: string(row.AgeRating),
		IsPubliclyVisible: row.IsPubliclyVisible, Words: row.Words, Chapters: row.Chapters,
		ReportsCount: row.ReportsCount, BanReason: row.BanReason,
	}
	if row.LatestPendingReportTime.Valid {
		result.LatestPendingReport = &BookPendingReport{ID: row.LatestPendingReportID, Number: row.LatestPendingReportNumber, Reason: row.LatestPendingReportReason, Time: row.LatestPendingReportTime.Time}
	}
	return result, nil
}

func nullableTime(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	return &value.Time
}

func (m *moderationBookService) GetBookChapters(ctx context.Context, query GetBookInfoQuery) ([]BookModerationChapter, error) {
	if err := m.authorize(ctx, query.ActorUserID); err != nil {
		return nil, err
	}
	rows, err := store.New(m.db).ModGetBookChapters(ctx, query.BookID)
	if err != nil {
		return nil, apperror.WrapUnexpectedDBError(err)
	}
	result := make([]BookModerationChapter, 0, len(rows))
	for _, row := range rows {
		result = append(result, BookModerationChapter{ID: row.ID, Name: row.Name, CreatedAt: row.CreatedAt.Time, UpdatedAt: nullableTime(row.UpdatedAt), Words: row.Words, IsPubliclyVisible: row.IsPubliclyVisible, HasPendingReports: row.HasPendingReports})
	}
	return result, nil
}

func (m *moderationBookService) changeBookValue(ctx context.Context, cmd ModerationPerformBookActionCommand, action BookActionType, change func(*store.Queries) error) error {
	if err := cmd.Validate(); err != nil {
		return err
	}
	if err := m.authorize(ctx, cmd.ActorUserID); err != nil {
		return err
	}
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	queries := store.New(m.db).WithTx(tx)
	if err = change(queries); err == nil {
		err = m.addBookLog(ctx, queries, cmd.BookID, action, cmd.Reason, cmd.ActorUserID)
	}
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return apperror.WrapUnexpectedDBError(err)
	}
	return tx.Commit(ctx)
}

func (m *moderationBookService) ChangeAgeRating(ctx context.Context, cmd ModerationPerformBookActionCommand) error {
	valid := map[string]bool{"?": true, "G": true, "PG": true, "PG-13": true, "R": true, "NC-17": true}
	if !valid[cmd.Value] {
		return InvalidModerationActionError.New("invalid age rating")
	}
	return m.changeBookValue(ctx, cmd, "change_age_rating", func(q *store.Queries) error {
		return q.ModChangeBookAgeRating(ctx, store.ModChangeBookAgeRatingParams{ID: cmd.BookID, AgeRating: store.AgeRating(cmd.Value)})
	})
}
func (m *moderationBookService) ChangeSummary(ctx context.Context, cmd ModerationPerformBookActionCommand) error {
	return m.changeBookValue(ctx, cmd, "change_summary", func(q *store.Queries) error {
		return q.ModChangeBookSummary(ctx, store.ModChangeBookSummaryParams{ID: cmd.BookID, Summary: cmd.Value})
	})
}

// GetBookLog implements ModerationBookService.
func (m *moderationBookService) GetBookLog(ctx context.Context, query GetBookLogQuery) (BookLogResult, error) {
	if err := m.authorize(ctx, query.ActorUserID); err != nil {
		return BookLogResult{}, err
	}
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
		ActionTypes: bookActionTypesToStrings(query.OfTypes),
		Limit:       pageSize + 1,
		Offset:      (page - 1) * pageSize,
	})
	if err != nil {
		return BookLogResult{}, apperror.WrapUnexpectedDBError(err)
	}

	count, err := queries.ModCountBookLogFiltered(ctx, store.ModCountBookLogFilteredParams{
		BookID:      query.BookID,
		ActionTypes: bookActionTypesToStrings(query.OfTypes),
	})
	if err != nil {
		return BookLogResult{}, apperror.WrapUnexpectedDBError(err)
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
			Action:        BookActionType(row.Type),
			Payload:       row.Payload,
			Reason:        row.Reason,
			ActorUserID:   uuidDbToDomain(row.ActorUserID),
			ActorUserName: row.ActorUserName,
		})
	}

	return result, nil
}

func bookActionTypesToStrings(types []BookActionType) []string {
	if types == nil {
		return nil
	}
	result := make([]string, len(types))
	for i, actionType := range types {
		result[i] = string(actionType)
	}
	return result
}

// BanBook implements ModerationBookService.
func (m *moderationBookService) BanBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}
	if err := m.authorize(ctx, cmd.ActorUserID); err != nil {
		return err
	}

	tx, err := m.db.Begin(ctx)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	queries := store.New(m.db).WithTx(tx)

	err = queries.ModSetBookBanned(ctx, store.ModSetBookBannedParams{
		ID:       cmd.BookID,
		IsBanned: true,
	})
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return apperror.WrapUnexpectedDBError(err)
	}

	err = m.addBookLog(ctx, queries, cmd.BookID, BookActionTypeBan, cmd.Reason, cmd.ActorUserID)
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	return nil
}

func (m *moderationBookService) addBookLog(ctx context.Context, queries *store.Queries, bookID int64, actionType BookActionType, reason string, actorUserID uuid.UUID) error {
	err := queries.ModAddBookLog(ctx, store.ModAddBookLogParams{
		ID:          GenID(),
		TargetID:    bookID,
		Reason:      reason,
		Payload:     nil,
		Time:        pgtype.Timestamptz{Valid: true, Time: time.Now()},
		Type:        string(actionType),
		ActorUserID: uuidDomainToDb(actorUserID),
	})
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	return nil
}

// PermanentlyRemoveBook implements ModerationBookService.
func (m *moderationBookService) PermanentlyRemoveBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}
	if err := m.authorize(ctx, cmd.ActorUserID); err != nil {
		return err
	}
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	queries := store.New(m.db).WithTx(tx)

	// The system user owns anonymized books after permanent removal.
	systemUser, err := queries.User_FindByLogin(ctx, "system")
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return apperror.WrapUnexpectedDBError(err)
	}
	if UserRole(systemUser.Role) != RoleSystem {
		dal.RollbackTx(ctx, tx)
		return ErrModerationForbidden
	}
	err = queries.ModPermRemoveBook(ctx, store.ModPermRemoveBookParams{ID: cmd.BookID, AuthorUserID: systemUser.ID})
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return apperror.WrapUnexpectedDBError(err)
	}

	err = m.addBookLog(ctx, queries, cmd.BookID, BookActionTypePermRemoval, cmd.Reason, cmd.ActorUserID)
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	return nil

}

// ShadowBanBook implements ModerationBookService.
func (m *moderationBookService) ShadowBanBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}
	if err := m.authorize(ctx, cmd.ActorUserID); err != nil {
		return err
	}

	tx, err := m.db.Begin(ctx)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	queries := store.New(m.db).WithTx(tx)

	err = queries.ModSetBookShadowBanned(ctx, store.ModSetBookShadowBannedParams{
		ID:             cmd.BookID,
		IsShadowBanned: true,
	})
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return apperror.WrapUnexpectedDBError(err)
	}

	err = m.addBookLog(ctx, queries, cmd.BookID, BookActionTypeShadowBan, cmd.Reason, cmd.ActorUserID)
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	return nil

}

// UnBanBook implements ModerationBookService.
func (m *moderationBookService) UnBanBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}
	if err := m.authorize(ctx, cmd.ActorUserID); err != nil {
		return err
	}

	tx, err := m.db.Begin(ctx)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	queries := store.New(m.db).WithTx(tx)

	err = queries.ModSetBookBanned(ctx, store.ModSetBookBannedParams{
		ID:       cmd.BookID,
		IsBanned: false,
	})
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return apperror.WrapUnexpectedDBError(err)
	}

	err = m.addBookLog(ctx, queries, cmd.BookID, BookActionTypeUnBan, cmd.Reason, cmd.ActorUserID)
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	return nil
}

// UnShadowBanBook implements ModerationBookService.
func (m *moderationBookService) UnShadowBanBook(ctx context.Context, cmd ModerationPerformBookActionCommand) error {
	if err := cmd.Validate(); err != nil {
		return err
	}
	if err := m.authorize(ctx, cmd.ActorUserID); err != nil {
		return err
	}

	tx, err := m.db.Begin(ctx)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	queries := store.New(m.db).WithTx(tx)

	err = queries.ModSetBookShadowBanned(ctx, store.ModSetBookShadowBannedParams{
		ID:             cmd.BookID,
		IsShadowBanned: false,
	})
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return apperror.WrapUnexpectedDBError(err)
	}

	err = m.addBookLog(ctx, queries, cmd.BookID, BookActionTypeUnShadowBan, cmd.Reason, cmd.ActorUserID)
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	return nil
}

func NewModerationBookService(db DB, auth ModerationAuthorizer) ModerationBookService {
	return &moderationBookService{db: db, auth: auth}
}
