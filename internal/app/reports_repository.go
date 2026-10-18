package app

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/app/dal"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/jackc/pgx/v5"
	"uuid"
)

type reportRepository struct{ db DB }

func (r *reportRepository) BookChapterExists(ctx context.Context, bookID, chapterID int64) (bool, error) {
	exists, err := store.New(moderationDB(ctx, r.db)).Report_BookChapterExists(ctx, store.Report_BookChapterExistsParams{BookID: bookID, ChapterID: chapterID})
	if err != nil {
		return false, apperror.WrapUnexpectedDBError(err)
	}
	return exists, nil
}

func (r *reportRepository) TargetExists(ctx context.Context, targetType ReportTargetType, targetID string) (bool, error) {
	q := store.New(moderationDB(ctx, r.db))
	var (
		exists bool
		err    error
	)
	switch targetType {
	case ReportTargetUser:
		id, parseErr := uuid.Parse(targetID)
		if parseErr != nil {
			return false, nil
		}
		exists, err = q.Report_UserExists(ctx, uuidDomainToDb(id))
	case ReportTargetBook, ReportTargetComment:
		id, parseErr := strconv.ParseInt(targetID, 10, 64)
		if parseErr != nil {
			return false, nil
		}
		if targetType == ReportTargetBook {
			exists, err = q.Report_BookExists(ctx, id)
		} else {
			exists, err = q.Report_CommentExists(ctx, id)
		}
	default:
		return false, nil
	}
	if err != nil {
		return false, apperror.WrapUnexpectedDBError(err)
	}
	return exists, nil
}

func (r *reportRepository) Create(ctx context.Context, report Report) (Report, error) {
	created, err := store.New(moderationDB(ctx, r.db)).Report_Create(ctx, store.Report_CreateParams{
		Time: timeToTimestamptz(report.Time), ReporterUserID: uuidDomainToDb(report.ReporterUserID),
		TargetType: string(report.TargetType), TargetID: report.TargetID, Reason: report.Reason, Description: report.Description,
		BookChapterID: int64NullableDomainToDb(report.BookChapterID), BookExcerpt: report.BookExcerpt,
	})
	if err != nil {
		return Report{}, apperror.WrapUnexpectedDBError(err)
	}
	report.ID = created.ID
	report.Number = created.Number
	return report, nil
}

func (r *reportRepository) GetByID(ctx context.Context, reportID int64) (Report, string, error) {
	row, err := store.New(moderationDB(ctx, r.db)).Report_GetByID(ctx, reportID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Report{}, "", ErrReportNotFound
		}
		return Report{}, "", apperror.WrapUnexpectedDBError(err)
	}
	chapterID := Null[int64]()
	if row.BookChapterID.Valid {
		chapterID = Value(row.BookChapterID.Int64)
	}
	return Report{
		ID:             row.ID,
		Number:         row.Number,
		Time:           timeDbToDomain(row.Time),
		ReporterUserID: uuidDbToDomain(row.ReporterUserID),
		TargetType:     ReportTargetType(row.TargetType),
		TargetID:       row.TargetID,
		Reason:         row.Reason,
		Description:    row.Description,
		BookChapterID:  chapterID,
		BookExcerpt:    row.BookExcerpt,
		Status:         row.Status,
		Priority:       row.Priority,
	}, row.ReporterUserName, nil

}

func (r *reportRepository) GetBookContext(ctx context.Context, reportID int64) (*ModerationReportBookContextData, error) {
	row, err := store.New(moderationDB(ctx, r.db)).Report_GetBookContext(ctx, reportID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, apperror.WrapUnexpectedDBError(err)
	}
	chapterID := Null[int64]()
	if row.BookChapterID.Valid {
		chapterID = Value(row.BookChapterID.Int64)
	}
	chapter := Null[string]()
	if row.Chapter.Valid {
		chapter = Value(row.Chapter.String)
	}
	return &ModerationReportBookContextData{
		BookID: row.BookID, ChapterID: chapterID, Title: row.Title, Author: row.Author,
		CoverID: row.Cover, Chapter: chapter, Excerpt: row.BookExcerpt, Rating: string(row.AgeRating), Warnings: row.Warnings,
		IsPermanentlyRemoved: row.IsPermRemoved, IsBanned: row.IsBanned, IsTrashed: row.IsTrashed, IsPubliclyVisible: row.IsPubliclyVisible,
		BookCreatedAt: timeDbToDomain(row.BookCreatedAt), ChapterCreatedAt: timeNullableDbToDomain(row.ChapterCreatedAt),
		ChapterUpdatedAt: timeNullableDbToDomain(row.ChapterUpdatedAt), ChapterContentUpdatedAt: timeNullableDbToDomain(row.ChapterContentUpdatedAt),
		RelatedReports: int(row.RelatedReports),
	}, nil
}

func (r *reportRepository) GetUserContext(ctx context.Context, reportID int64) (*ModerationReportUserContextData, error) {
	row, err := store.New(moderationDB(ctx, r.db)).Report_GetUserContext(ctx, reportID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, apperror.WrapUnexpectedDBError(err)
	}
	return &ModerationReportUserContextData{ID: uuidDbToDomain(row.ID), Name: row.Name, Email: row.Email, About: row.About, Role: string(row.Role), JoinedAt: timeDbToDomain(row.JoinedAt), IsBanned: row.IsBanned, RelatedReports: int(row.RelatedReports)}, nil
}

func (r *reportRepository) GetCommentContext(ctx context.Context, reportID int64) (*ModerationReportCommentContextData, error) {
	row, err := store.New(moderationDB(ctx, r.db)).Report_GetCommentContext(ctx, reportID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, apperror.WrapUnexpectedDBError(err)
	}
	return &ModerationReportCommentContextData{ID: row.ID, ChapterID: row.ChapterID, BookID: row.BookID, Content: row.Content, AuthorID: uuidDbToDomain(row.AuthorID), AuthorName: row.AuthorName, ChapterName: row.ChapterName, BookName: row.BookName, CreatedAt: timeDbToDomain(row.CreatedAt), UpdatedAt: timeNullableDbToDomain(row.UpdatedAt), DeletedAt: timeNullableDbToDomain(row.DeletedAt), RelatedReports: int(row.RelatedReports)}, nil
}

func (r *reportRepository) Search(ctx context.Context, filter SearchModerationReportsQuery, limit, offset int32) ([]ModerationReportListEntry, int64, error) {
	queries := store.New(moderationDB(ctx, r.db))
	total, err := queries.Report_CountSearch(ctx, store.Report_CountSearchParams{Search: filter.Search, TargetType: filter.TargetType, TargetID: filter.TargetID, ChapterID: filter.ChapterID, Status: filter.Status, Priority: filter.Priority, DateFrom: loginHistoryTime(filter.DateFrom), DateTo: loginHistoryTime(filter.DateTo)})
	if err != nil {
		return nil, 0, apperror.WrapUnexpectedDBError(err)
	}
	rows, err := queries.Report_Search(ctx, store.Report_SearchParams{Search: filter.Search, TargetType: filter.TargetType, TargetID: filter.TargetID, ChapterID: filter.ChapterID, Status: filter.Status, Priority: filter.Priority, DateFrom: loginHistoryTime(filter.DateFrom), DateTo: loginHistoryTime(filter.DateTo), Sort: filter.Sort, PageLimit: limit, PageOffset: offset})
	if err != nil {
		return nil, 0, apperror.WrapUnexpectedDBError(err)
	}
	entries := MapSlice(rows, func(row store.Report_SearchRow) ModerationReportListEntry {
		return ModerationReportListEntry{Report: Report{ID: row.ID, Number: row.Number, Time: timeDbToDomain(row.Time), ReporterUserID: uuidDbToDomain(row.ReporterUserID), TargetType: ReportTargetType(row.TargetType), TargetID: row.TargetID, Reason: row.Reason, Description: row.Description, Status: row.Status, Priority: row.Priority}, ReporterUserName: row.ReporterUserName}
	})
	return entries, total, nil
}

func (r *reportRepository) GetEvents(ctx context.Context, reportID int64) ([]ReportDecision, error) {
	rows, err := store.New(moderationDB(ctx, r.db)).Report_GetEvents(ctx, reportID)
	if err != nil {
		return nil, apperror.WrapUnexpectedDBError(err)
	}
	result := make([]ReportDecision, 0, len(rows))
	for _, row := range rows {
		var payload reportDecisionLogPayload
		if err = json.Unmarshal(row.Payload, &payload); err != nil {
			return nil, apperror.WrapUnexpectedAppError(err)
		}
		result = append(result, ReportDecision{OldPriority: payload.OldPriority, NewPriority: payload.NewPriority, Kind: row.Type, Time: timeDbToDomain(row.Time), ActorUserID: uuidDbToDomain(row.ActorUserID), ActorUserName: row.ActorUserName, Disposition: payload.Disposition, Action: payload.Action, PolicyReason: row.Reason, InternalNote: payload.InternalNote, NotifyTarget: payload.NotifyTarget, Payload: payload.ActionPayload})
	}
	return result, nil
}

type reportDecisionLogPayload struct {
	OperationID   string          `json:"operationId,omitempty"`
	RequestHash   string          `json:"requestHash,omitempty"`
	OldPriority   string          `json:"oldPriority,omitempty"`
	NewPriority   string          `json:"newPriority,omitempty"`
	Version       int             `json:"version"`
	Disposition   string          `json:"disposition"`
	Action        string          `json:"action,omitempty"`
	InternalNote  string          `json:"internalNote,omitempty"`
	NotifyTarget  bool            `json:"notifyTarget"`
	ActionPayload json.RawMessage `json:"actionPayload"`
}

func (r *reportRepository) AddDecision(ctx context.Context, reportID int64, status string, decision ReportDecision) error {
	tx, err := moderationDB(ctx, r.db).Begin(ctx)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	q := store.New(moderationDB(ctx, r.db)).WithTx(tx)
	rows, err := q.Report_SetStatus(ctx, store.Report_SetStatusParams{ReportID: reportID, NewStatus: status})
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return apperror.WrapUnexpectedDBError(err)
	}
	if rows != 1 {
		dal.RollbackTx(ctx, tx)
		return ErrReportAlreadyResolved
	}
	actionPayload := decision.Payload
	if len(actionPayload) == 0 {
		actionPayload = json.RawMessage(`{}`)
	}
	payload, err := json.Marshal(reportDecisionLogPayload{OperationID: decision.OperationID, RequestHash: decision.RequestHash, Version: 1, Disposition: decision.Disposition, Action: decision.Action, InternalNote: decision.InternalNote, NotifyTarget: decision.NotifyTarget, ActionPayload: actionPayload})
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return apperror.WrapUnexpectedAppError(err)
	}
	err = q.Moderation_AddLog(ctx, store.Moderation_AddLogParams{ID: GenID(), Time: timeToTimestamptz(decision.Time), Type: "report_decision", TargetType: "report", TargetID: strconv.FormatInt(reportID, 10), Payload: payload, ActorUserID: uuidDomainToDb(decision.ActorUserID), Reason: decision.PolicyReason})
	if err != nil {
		dal.RollbackTx(ctx, tx)
		return apperror.WrapUnexpectedDBError(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	return nil
}

func NewReportRepository(db DB) ReportRepository { return &reportRepository{db: db} }

func (r *reportRepository) WithLockedReport(ctx context.Context, id int64, run func(context.Context) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	defer tx.Rollback(ctx)
	if _, err = store.New(tx).Report_Lock(ctx, id); err != nil {
		if err == pgx.ErrNoRows {
			return ErrReportNotFound
		}
		return apperror.WrapUnexpectedDBError(err)
	}
	if err = run(context.WithValue(ctx, moderationTransactionKey{}, tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *reportRepository) AddNote(ctx context.Context, id int64, event ReportDecision, priority string) error {
	q := store.New(moderationDB(ctx, r.db))
	payload := map[string]any{"version": 1, "internalNote": event.InternalNote}
	kind := "report_note"
	if priority != "" {
		report, _, err := r.GetByID(ctx, id)
		if err != nil {
			return err
		}
		payload["oldPriority"], payload["newPriority"] = report.Priority, priority
		kind = "report_priority_changed"
		if err := q.Report_SetPriority(ctx, store.Report_SetPriorityParams{ID: id, Priority: priority}); err != nil {
			return apperror.WrapUnexpectedDBError(err)
		}
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return q.Moderation_AddLog(ctx, store.Moderation_AddLogParams{ID: GenID(), Time: timeToTimestamptz(event.Time), Type: kind, TargetType: "report", TargetID: strconv.FormatInt(id, 10), Payload: encoded, ActorUserID: uuidDomainToDb(event.ActorUserID), Reason: event.PolicyReason})
}

func (r *reportRepository) FindOperation(ctx context.Context, id int64, key string) (uuid.UUID, string, bool, error) {
	row, err := store.New(moderationDB(ctx, r.db)).Report_FindOperation(ctx, store.Report_FindOperationParams{ReportID: id, OperationID: key})
	if err == pgx.ErrNoRows {
		return uuid.Nil(), "", false, nil
	}
	if err != nil {
		return uuid.Nil(), "", false, apperror.WrapUnexpectedDBError(err)
	}
	return uuidDbToDomain(row.ActorUserID), row.RequestHash, true, nil
}
