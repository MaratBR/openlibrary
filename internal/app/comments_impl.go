package app

import (
	"context"
	"fmt"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
)

type commentsService struct {
	db store.DBTX
}

// AddComment implements CommentsService.
func (c *commentsService) AddComment(ctx context.Context, command AddCommentCommand) (AddCommentResult, error) {
	queries := store.New(c.db)
	id := GenID()
	err := queries.InsertComment(ctx, store.InsertCommentParams{
		ID:        id,
		ChapterID: command.ChapterID,
		ParentID:  int64NullableDomainToDb(command.ParentCommentID),
		UserID:    uuidDomainToDb(command.UserID),
		Content:   command.Content,
	})
	if err != nil {
		return AddCommentResult{}, wrapUnexpectedDBError(err)
	}

	comment, err := c.getByID(ctx, id)
	if err != nil {
		return AddCommentResult{}, UnexpectedError.New(err.Error())
	}

	return AddCommentResult{Comment: comment}, err
}

// GetList implements CommentsService.
func (c *commentsService) GetList(ctx context.Context, query GetCommentsQuery) (result GetCommentsResult, err error) {
	queries := store.New(c.db)

	result.Cursor = query.Cursor

	if query.Cursor == 0 {
		var rows []store.GetChapterCommentsRow
		rows, err = queries.GetChapterComments(ctx, store.GetChapterCommentsParams{
			ChapterID: query.ChapterID,
			Limit:     query.Limit,
		})
		if err != nil {
			err = wrapUnexpectedDBError(err)
			return
		}
		result.Comments = MapSlice(rows, func(r store.GetChapterCommentsRow) CommentDto {
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

		var rows []store.GetChapterCommentsAfterRow
		rows, err = queries.GetChapterCommentsAfter(ctx, store.GetChapterCommentsAfterParams{
			ChapterID: query.ChapterID,
			Limit:     query.Limit,
			CreatedAt: timeToTimestamptz(ts),
		})
		if err != nil {
			err = wrapUnexpectedDBError(err)
			return
		}
		result.Comments = MapSlice(rows, func(r store.GetChapterCommentsAfterRow) CommentDto {
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
	}

	return

}

// UpdateComment implements CommentsService.
func (c *commentsService) UpdateComment(ctx context.Context, command UpdateCommentCommand) (UpdateCommentResult, error) {
	queries := store.New(c.db)

	result, err := queries.UpdateComment(ctx, store.UpdateCommentParams{
		ID:      command.ID,
		Content: command.Content,
	})
	if err != nil {
		return UpdateCommentResult{}, wrapUnexpectedDBError(err)
	}

	if result.RowsAffected() == 0 {
		return UpdateCommentResult{}, ErrTypeCommentNotFound.New(fmt.Sprintf("comment with id %d not found", command.ID))
	}

	comment, err := c.getByID(ctx, command.ID)
	if err != nil {
		return UpdateCommentResult{}, UnexpectedError.New(err.Error())
	}

	return UpdateCommentResult{Comment: comment}, err
}

func (c *commentsService) getByID(ctx context.Context, id int64) (CommentDto, error) {
	row, err := store.New(c.db).GetCommentWithUserByID(ctx, id)
	if err != nil {
		if err == store.ErrNoRows {
			return CommentDto{}, ErrTypeCommentNotFound.New(fmt.Sprintf("comment with id %d not found", id))
		}
		return CommentDto{}, wrapUnexpectedDBError(err)
	}

	return CommentDto{
		ID:        row.ID,
		Content:   row.Content,
		User:      CommentUserDto{ID: uuidDbToDomain(row.UserID), Name: row.UserName, Avatar: getUserAvatar(row.UserName, 84)},
		CreatedAt: timeDbToDomain(row.CreatedAt),
		UpdatedAt: timeNullableDbToDomain(row.UpdatedAt),
	}, nil
}

func NewCommentsService(db store.DBTX) CommentsService {
	return &commentsService{db: db}
}
