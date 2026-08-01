package app

import (
	"context"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type moderationUserRepository struct{ queries *store.Queries }

func (r *moderationUserRepository) GetUserInfo(ctx context.Context, userID uuid.UUID) (ModerationUserInfo, error) {
	row, err := r.queries.Moderation_GetUserInfo(ctx, uuidDomainToDb(userID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return ModerationUserInfo{}, ErrUserNotFound
		}
		return ModerationUserInfo{}, apperror.WrapUnexpectedDBError(err)
	}
	return ModerationUserInfo{
		ID: uuidDbToDomain(row.ID), Name: row.Name, Email: row.Email, About: row.About,
		Avatar: getUserAvatar(row.Name, 256), JoinedAt: timeDbToDomain(row.JoinedAt),
		Role: UserRole(row.Role), IsBanned: row.IsBanned, IsEmailVerified: row.EmailVerified,
		BooksTotal: row.BooksTotal, CommentsTotal: row.CommentsTotal, FollowersTotal: row.FollowersTotal,
	}, nil
}

func (r *moderationUserRepository) GetBooks(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]ModerationUserBook, error) {
	rows, err := r.queries.Moderation_GetUserBooks(ctx, store.Moderation_GetUserBooksParams{UserID: uuidDomainToDb(userID), PageLimit: limit, PageOffset: offset})
	if err != nil {
		return nil, apperror.WrapUnexpectedDBError(err)
	}
	return MapSlice(rows, func(row store.Moderation_GetUserBooksRow) ModerationUserBook {
		return ModerationUserBook{ID: row.ID, Name: row.Name, CreatedAt: timeDbToDomain(row.CreatedAt), IsPubliclyVisible: row.IsPubliclyVisible, IsBanned: row.IsBanned, IsTrashed: row.IsTrashed}
	}), nil
}

func (r *moderationUserRepository) GetComments(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]ModerationUserComment, error) {
	rows, err := r.queries.Moderation_GetUserComments(ctx, store.Moderation_GetUserCommentsParams{UserID: uuidDomainToDb(userID), PageLimit: limit, PageOffset: offset})
	if err != nil {
		return nil, apperror.WrapUnexpectedDBError(err)
	}
	return MapSlice(rows, func(row store.Moderation_GetUserCommentsRow) ModerationUserComment {
		return ModerationUserComment{ID: row.ID, Content: row.Content, CreatedAt: timeDbToDomain(row.CreatedAt), Deleted: row.DeletedAt.Valid, ChapterID: row.ChapterID, ChapterName: row.ChapterName, BookID: row.BookID, BookName: row.BookName}
	}), nil
}

func (r *moderationUserRepository) GetHistory(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]ModerationUserHistoryEntry, int64, error) {
	total, err := r.queries.Moderation_CountUserHistory(ctx, uuidDomainToDb(userID))
	if err != nil {
		return nil, 0, apperror.WrapUnexpectedDBError(err)
	}
	rows, err := r.queries.Moderation_GetUserHistory(ctx, store.Moderation_GetUserHistoryParams{UserID: uuidDomainToDb(userID), PageLimit: limit, PageOffset: offset})
	if err != nil {
		return nil, 0, apperror.WrapUnexpectedDBError(err)
	}
	return MapSlice(rows, func(row store.Moderation_GetUserHistoryRow) ModerationUserHistoryEntry {
		return ModerationUserHistoryEntry{ID: row.ID, Time: timeDbToDomain(row.Time), Type: row.Type, Reason: row.Reason, Payload: row.Payload, ActorUserID: uuidDbToDomain(row.ActorUserID), ActorUserName: row.ActorUserName.String}
	}), total, nil
}

func (r *moderationUserRepository) GetReports(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]ModerationUserReport, int64, error) {
	total, err := r.queries.Moderation_CountUserReports(ctx, uuidDomainToDb(userID))
	if err != nil {
		return nil, 0, apperror.WrapUnexpectedDBError(err)
	}
	rows, err := r.queries.Moderation_GetUserReports(ctx, store.Moderation_GetUserReportsParams{UserID: uuidDomainToDb(userID), PageLimit: limit, PageOffset: offset})
	if err != nil {
		return nil, 0, apperror.WrapUnexpectedDBError(err)
	}
	return MapSlice(rows, func(row store.Moderation_GetUserReportsRow) ModerationUserReport {
		return ModerationUserReport{ID: row.ID, Number: row.Number, Time: timeDbToDomain(row.Time), TargetType: ReportTargetType(row.TargetType), TargetID: row.TargetID, Reason: row.Reason, Description: row.Description, ReporterUserID: uuidDbToDomain(row.ReporterUserID), ReporterUserName: row.ReporterUserName.String}
	}), total, nil
}

func (r *moderationUserRepository) CreateTemporaryRandomReports(ctx context.Context, userID uuid.UUID, count int) error {
	for range count {
		reporterUserID, err := r.queries.Moderation_GetRandomUserID(ctx)
		if err != nil {
			return apperror.WrapUnexpectedDBError(err)
		}
		reason, description := temporaryRandomReportValues()
		if _, err = r.queries.Report_Create(ctx, store.Report_CreateParams{
			Time: timeToTimestamptz(time.Now()), ReporterUserID: reporterUserID,
			TargetType: string(ReportTargetUser), TargetID: userID.String(), Reason: reason, Description: description,
		}); err != nil {
			return apperror.WrapUnexpectedDBError(err)
		}
	}
	return nil
}

func NewModerationUserRepository(db DB) ModerationUserRepository {
	return &moderationUserRepository{queries: store.New(db)}
}
