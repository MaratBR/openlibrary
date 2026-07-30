package app

import (
	"context"
	"encoding/json"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/app/dal"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type contentModerationRepository struct {
	db DB
}

func (r *contentModerationRepository) transaction(ctx context.Context, run func(*store.Queries) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	if err = run(store.New(r.db).WithTx(tx)); err != nil {
		dal.RollbackTx(ctx, tx)
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	return nil
}

func (r *contentModerationRepository) SetChapterVisibilityAndLog(ctx context.Context, chapterID int64, visible bool, log ModerationAuditEntry) error {
	return r.transaction(ctx, func(q *store.Queries) error {
		chapter, err := q.Moderation_GetChapter(ctx, chapterID)
		if err != nil {
			if err == pgx.ErrNoRows {
				return ErrTypeChapterDoesNotExist.New("chapter not found")
			}
			return apperror.WrapUnexpectedDBError(err)
		}
		if err = q.Moderation_SetChapterVisibility(ctx, store.Moderation_SetChapterVisibilityParams{ID: chapterID, IsPubliclyVisible: visible}); err != nil {
			return apperror.WrapUnexpectedDBError(err)
		}
		payload, err := json.Marshal(map[string]any{"action": log.Action, "chapterId": chapterID, "visible": visible})
		if err != nil {
			return apperror.WrapUnexpectedAppError(err)
		}
		return addPersistedBookLog(ctx, q, chapter.BookID, log, payload)
	})
}

func addPersistedBookLog(ctx context.Context, q *store.Queries, bookID int64, log ModerationAuditEntry, payload []byte) error {
	err := q.ModAddBookLog(ctx, store.ModAddBookLogParams{
		ID: log.ID, Time: timeToTimestamptz(log.Time), BookID: bookID,
		ActionType: store.BookActionTypeSignificantUpdate, Payload: payload,
		ActorUserID: uuidDomainToDb(log.ActorUserID), Reason: log.Reason,
	})
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	return nil
}

func (r *contentModerationRepository) SetCommentRemovedAndLog(ctx context.Context, commentID int64, removed bool, log ModerationAuditEntry) error {
	return r.transaction(ctx, func(q *store.Queries) error {
		if _, err := q.Comment_GetByID(ctx, commentID); err != nil {
			if err == pgx.ErrNoRows {
				return ErrTypeCommentNotFound.New("comment not found")
			}
			return apperror.WrapUnexpectedDBError(err)
		}
		if err := q.Moderation_SetCommentRemoved(ctx, store.Moderation_SetCommentRemovedParams{ID: commentID, Removed: removed}); err != nil {
			return apperror.WrapUnexpectedDBError(err)
		}
		err := q.Moderation_AddCommentLog(ctx, store.Moderation_AddCommentLogParams{
			ID: log.ID, Time: timeToTimestamptz(log.Time), CommentID: commentID,
			ActionType: string(log.Action), Payload: log.Payload,
			ActorUserID: uuidDomainToDb(log.ActorUserID), Reason: log.Reason,
		})
		if err != nil {
			return apperror.WrapUnexpectedDBError(err)
		}
		return nil
	})
}

func (r *contentModerationRepository) SetUserBanAndLog(ctx context.Context, userID uuid.UUID, expiresAt time.Time, banned bool, log ModerationAuditEntry) error {
	return r.transaction(ctx, func(q *store.Queries) error {
		if _, err := q.User_Get(ctx, uuidDomainToDb(userID)); err != nil {
			if err == pgx.ErrNoRows {
				return ErrUserNotFound
			}
			return apperror.WrapUnexpectedDBError(err)
		}
		if banned {
			if err := q.Moderation_AddUserBan(ctx, store.Moderation_AddUserBanParams{
				ID: GenID(), UserID: uuidDomainToDb(userID), CreatedAt: timeToTimestamptz(log.Time),
				BannedByUserID: uuidDomainToDb(log.ActorUserID), Note: log.Reason, ExpiresAt: timeToTimestamptz(expiresAt),
			}); err != nil {
				return apperror.WrapUnexpectedDBError(err)
			}
		}
		if err := q.Moderation_SetUserBanned(ctx, store.Moderation_SetUserBannedParams{ID: uuidDomainToDb(userID), IsBanned: banned}); err != nil {
			return apperror.WrapUnexpectedDBError(err)
		}
		if banned {
			if err := q.Session_TerminateAllByUserID(ctx, uuidDomainToDb(userID)); err != nil {
				return apperror.WrapUnexpectedDBError(err)
			}
		}
		return addPersistedUserLog(ctx, q, userID, log, map[string]any{"expiresAt": expiresAt, "banned": banned})
	})
}

func (r *contentModerationRepository) RenameUserAndLog(ctx context.Context, userID uuid.UUID, name string, log ModerationAuditEntry) error {
	return r.transaction(ctx, func(q *store.Queries) error {
		user, err := q.User_Get(ctx, uuidDomainToDb(userID))
		if err != nil {
			if err == pgx.ErrNoRows {
				return ErrUserNotFound
			}
			return apperror.WrapUnexpectedDBError(err)
		}
		if err = q.Moderation_RenameUser(ctx, store.Moderation_RenameUserParams{ID: uuidDomainToDb(userID), Name: name}); err != nil {
			return apperror.WrapUnexpectedDBError(err)
		}
		return addPersistedUserLog(ctx, q, userID, log, map[string]any{"oldName": user.Name, "newName": name})
	})
}

func (r *contentModerationRepository) ChangeUserAboutAndLog(ctx context.Context, userID uuid.UUID, about string, log ModerationAuditEntry) error {
	return r.transaction(ctx, func(q *store.Queries) error {
		user, err := q.User_Get(ctx, uuidDomainToDb(userID))
		if err != nil {
			if err == pgx.ErrNoRows {
				return ErrUserNotFound
			}
			return apperror.WrapUnexpectedDBError(err)
		}
		if err = q.Moderation_ChangeUserAbout(ctx, store.Moderation_ChangeUserAboutParams{ID: uuidDomainToDb(userID), About: about}); err != nil {
			return apperror.WrapUnexpectedDBError(err)
		}
		return addPersistedUserLog(ctx, q, userID, log, map[string]any{"oldAbout": user.About, "newAbout": about})
	})
}

func addPersistedUserLog(ctx context.Context, q *store.Queries, userID uuid.UUID, log ModerationAuditEntry, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return apperror.WrapUnexpectedAppError(err)
	}
	err = q.Moderation_AddUserLog(ctx, store.Moderation_AddUserLogParams{
		ID: log.ID, UserID: uuidDomainToDb(userID), ActorUserID: uuidDomainToDb(log.ActorUserID),
		ActionType: string(log.Action), Payload: payload, Time: timeToTimestamptz(log.Time), Reason: log.Reason,
	})
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	return nil
}

func NewContentModerationRepository(db DB) ContentModerationRepository {
	return &contentModerationRepository{db: db}
}
