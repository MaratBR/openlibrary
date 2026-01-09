package app

import (
	"context"
	"fmt"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/store"
)

type commentsService struct {
	db store.DBTX
}

// AddComment implements CommentsService.
func (c *commentsService) AddComment(ctx context.Context, command AddCommentCommand) (AddCommentResult, error) {
	queries := store.New(c.db)
	id := GenID()
	err := queries.Comment_Insert(ctx, store.Comment_InsertParams{
		ID:        id,
		ChapterID: command.ChapterID,
		ParentID:  int64NullableDomainToDb(command.ParentCommentID),
		UserID:    uuidDomainToDb(command.UserID),
		Content:   command.Content,
	})
	if err != nil {
		return AddCommentResult{}, apperror.WrapUnexpectedDBError(err)
	}

	comment, err := c.getByID(ctx, id)
	if err != nil {
		return AddCommentResult{}, apperror.UnexpectedError.New(err.Error())
	}

	return AddCommentResult{Comment: comment}, err
}

// GetList implements CommentsService.
func (c *commentsService) GetList(ctx context.Context, query GetCommentsQuery) (result GetCommentsResult, err error) {
	queries := store.New(c.db)

	result.Cursor = query.Cursor

	if query.Cursor == 0 {
		var rows []store.Comment_GetByChapterRow
		rows, err = queries.Comment_GetByChapter(ctx, store.Comment_GetByChapterParams{
			ChapterID: query.ChapterID,
			Limit:     query.Limit,
		})
		if err != nil {
			err = apperror.WrapUnexpectedDBError(err)
			return
		}
		result.Comments = MapSlice(rows, func(r store.Comment_GetByChapterRow) CommentDto {
			return CommentDto{
				ID:        r.ID,
				Content:   r.Content,
				User:      CommentUserDto{ID: uuidDbToDomain(r.UserID), Name: r.UserName, Avatar: getUserAvatar(r.UserName, 84)},
				CreatedAt: timeDbToDomain(r.CreatedAt),
				UpdatedAt: timeNullableDbToDomain(r.UpdatedAt),
			}
		})
	} else {
		ts := time.Unix(int64(query.Cursor), 0)

		var rows []store.Comment_GetByChapterAfterRow
		rows, err = queries.Comment_GetByChapterAfter(ctx, store.Comment_GetByChapterAfterParams{
			ChapterID: query.ChapterID,
			Limit:     query.Limit,
			CreatedAt: timeToTimestamptz(ts),
		})
		if err != nil {
			err = apperror.WrapUnexpectedDBError(err)
			return
		}
		result.Comments = MapSlice(rows, func(r store.Comment_GetByChapterAfterRow) CommentDto {
			return CommentDto{
				ID:        r.ID,
				Content:   r.Content,
				User:      CommentUserDto{ID: uuidDbToDomain(r.UserID), Name: r.UserName, Avatar: getUserAvatar(r.UserName, 84)},
				CreatedAt: timeDbToDomain(r.CreatedAt),
				UpdatedAt: timeNullableDbToDomain(r.UpdatedAt),
			}
		})
	}

	if len(result.Comments) == 0 {
		result.NextCursor = 0
	} else {
		unixTs := result.Comments[len(result.Comments)-1].CreatedAt.Unix()
		result.NextCursor = uint32(unixTs)

		if query.ActorUserID.Valid {
			commentIds := make([]int64, len(result.Comments))
			for i := range result.Comments {
				commentIds[i] = result.Comments[i].ID
			}

			var likedComments []store.Comment_GetLikedCommentsRow
			likedComments, err = queries.Comment_GetLikedComments(ctx, store.Comment_GetLikedCommentsParams{
				UserID: uuidDomainToDb(query.ActorUserID.UUID),
				Ids:    commentIds,
			})
			if err != nil {
				err = apperror.WrapUnexpectedDBError(err)
				return
			}
			likedCommentsMapping := map[int64]time.Time{}

			for _, likedComment := range likedComments {
				likedCommentsMapping[likedComment.CommentID] = timeDbToDomain(likedComment.LikedAt)
			}

			for i := range result.Comments {
				comment := result.Comments[i]
				if likedAt, ok := likedCommentsMapping[comment.ID]; ok {
					comment.LikedAt = Value(likedAt)
					result.Comments[i] = comment
				}
			}
		}
	}

	return

}

// UpdateComment implements CommentsService.
func (c *commentsService) UpdateComment(ctx context.Context, command UpdateCommentCommand) (UpdateCommentResult, error) {
	queries := store.New(c.db)

	result, err := queries.Comment_Update(ctx, store.Comment_UpdateParams{
		ID:      command.ID,
		Content: command.Content,
	})
	if err != nil {
		return UpdateCommentResult{}, apperror.WrapUnexpectedDBError(err)
	}

	if result.RowsAffected() == 0 {
		return UpdateCommentResult{}, ErrTypeCommentNotFound.New(fmt.Sprintf("comment with id %d not found", command.ID))
	}

	comment, err := c.getByID(ctx, command.ID)
	if err != nil {
		return UpdateCommentResult{}, apperror.UnexpectedError.New(err.Error())
	}

	return UpdateCommentResult{Comment: comment}, err
}

func (c *commentsService) getByID(ctx context.Context, id int64) (CommentDto, error) {
	row, err := store.New(c.db).Comment_GetWithUserByID(ctx, id)
	if err != nil {
		if err == store.ErrNoRows {
			return CommentDto{}, ErrTypeCommentNotFound.New(fmt.Sprintf("comment with id %d not found", id))
		}
		return CommentDto{}, apperror.WrapUnexpectedDBError(err)
	}

	return CommentDto{
		ID:        row.ID,
		Content:   row.Content,
		User:      CommentUserDto{ID: uuidDbToDomain(row.UserID), Name: row.UserName, Avatar: getUserAvatar(row.UserName, 84)},
		CreatedAt: timeDbToDomain(row.CreatedAt),
		UpdatedAt: timeNullableDbToDomain(row.UpdatedAt),
	}, nil
}

func (c *commentsService) LikeComment(ctx context.Context, command LikeCommentCommand) (bool, error) {
	queries := store.New(c.db)
	userID := uuidDomainToDb(command.UserID)
	var err error

	if command.Like {
		err = queries.Comment_Like(ctx, store.Comment_LikeParams{
			CommentID: command.CommentID,
			UserID:    userID,
		})
	} else {
		err = queries.Comment_UnLike(ctx, store.Comment_UnLikeParams{
			CommentID: command.CommentID,
			UserID:    userID,
		})
	}

	if err != nil {
		return false, apperror.WrapUnexpectedDBError(err)
	}

	return command.Like, nil
}

func NewCommentsService(db store.DBTX) CommentsService {
	return &commentsService{db: db}
}
