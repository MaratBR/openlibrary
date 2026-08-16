package app

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/app/dal"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type reportRepository struct{ db DB }

func (r *reportRepository) BookChapterExists(ctx context.Context, bookID, chapterID int64) (bool, error) {
	exists, err := store.New(r.db).Report_BookChapterExists(ctx, store.Report_BookChapterExistsParams{BookID: bookID, ChapterID: chapterID})
	if err != nil {
		return false, apperror.WrapUnexpectedDBError(err)
	}
	return exists, nil
}

func (r *reportRepository) TargetExists(ctx context.Context, targetType ReportTargetType, targetID string) (bool, error) {
	q := store.New(r.db)
	var (
		exists bool
		err    error
	)
	switch targetType {
	case ReportTargetUser:
		id, parseErr := uuid.FromString(targetID)
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
	created, err := store.New(r.db).Report_Create(ctx, store.Report_CreateParams{
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
	row, err := store.New(r.db).Report_GetByID(ctx, reportID)
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
	row, err := store.New(r.db).Report_GetBookContext(ctx, reportID)
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
	row, err := store.New(r.db).Report_GetUserContext(ctx, reportID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, apperror.WrapUnexpectedDBError(err)
	}
	return &ModerationReportUserContextData{ID: uuidDbToDomain(row.ID), Name: row.Name, Email: row.Email, About: row.About, Role: string(row.Role), JoinedAt: timeDbToDomain(row.JoinedAt), IsBanned: row.IsBanned, RelatedReports: int(row.RelatedReports)}, nil
}

func (r *reportRepository) GetCommentContext(ctx context.Context, reportID int64) (*ModerationReportCommentContextData, error) {
	row, err := store.New(r.db).Report_GetCommentContext(ctx, reportID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, apperror.WrapUnexpectedDBError(err)
	}
	return &ModerationReportCommentContextData{ID: row.ID, ChapterID: row.ChapterID, BookID: row.BookID, Content: row.Content, AuthorID: uuidDbToDomain(row.AuthorID), AuthorName: row.AuthorName, ChapterName: row.ChapterName, BookName: row.BookName, CreatedAt: timeDbToDomain(row.CreatedAt), UpdatedAt: timeNullableDbToDomain(row.UpdatedAt), DeletedAt: timeNullableDbToDomain(row.DeletedAt), RelatedReports: int(row.RelatedReports)}, nil
}

func (r *reportRepository) Search(ctx context.Context, search, targetType string, limit, offset int32) ([]ModerationReportListEntry, int64, error) {
	queries := store.New(r.db)
	total, err := queries.Report_CountSearch(ctx, store.Report_CountSearchParams{Search: search, TargetType: targetType})
	if err != nil {
		return nil, 0, apperror.WrapUnexpectedDBError(err)
	}
	rows, err := queries.Report_Search(ctx, store.Report_SearchParams{Search: search, TargetType: targetType, PageLimit: limit, PageOffset: offset})
	if err != nil {
		return nil, 0, apperror.WrapUnexpectedDBError(err)
	}
	entries := MapSlice(rows, func(row store.Report_SearchRow) ModerationReportListEntry {
		return ModerationReportListEntry{Report: Report{ID: row.ID, Number: row.Number, Time: timeDbToDomain(row.Time), ReporterUserID: uuidDbToDomain(row.ReporterUserID), TargetType: ReportTargetType(row.TargetType), TargetID: row.TargetID, Reason: row.Reason, Description: row.Description}, ReporterUserName: row.ReporterUserName}
	})
	return entries, total, nil
}

func (r *reportRepository) GetEvents(ctx context.Context, reportID int64) ([]ReportDecision, error) {
	rows, err := store.New(r.db).Report_GetEvents(ctx, reportID)
	if err != nil {
		return nil, apperror.WrapUnexpectedDBError(err)
	}
	result := make([]ReportDecision, 0, len(rows))
	for _, row := range rows {
		var payload reportDecisionLogPayload
		if err = json.Unmarshal(row.Payload, &payload); err != nil {
			return nil, apperror.WrapUnexpectedAppError(err)
		}
		result = append(result, ReportDecision{Time: timeDbToDomain(row.Time), ActorUserID: uuidDbToDomain(row.ActorUserID), ActorUserName: row.ActorUserName, Disposition: payload.Disposition, Action: payload.Action, PolicyReason: row.Reason, InternalNote: payload.InternalNote, NotifyTarget: payload.NotifyTarget, Payload: payload.ActionPayload})
	}
	return result, nil
}

type reportDecisionLogPayload struct {
	Version       int             `json:"version"`
	Disposition   string          `json:"disposition"`
	Action        string          `json:"action,omitempty"`
	InternalNote  string          `json:"internalNote,omitempty"`
	NotifyTarget  bool            `json:"notifyTarget"`
	ActionPayload json.RawMessage `json:"actionPayload"`
}

func (r *reportRepository) AddDecision(ctx context.Context, reportID int64, status string, decision ReportDecision) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	q := store.New(r.db).WithTx(tx)
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
	payload, err := json.Marshal(reportDecisionLogPayload{Version: 1, Disposition: decision.Disposition, Action: decision.Action, InternalNote: decision.InternalNote, NotifyTarget: decision.NotifyTarget, ActionPayload: actionPayload})
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
