package app

import (
	"context"
	"fmt"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
)

type commentsService struct {
	db store.DBTX
}

// AddComment implements CommentsService.
func (c *commentsService) AddComment(ctx context.Context, command AddCommentCommand) (AddCommentResult, error) {
	if err := command.Validate(); err != nil {
		return AddCommentResult{}, err
	}
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
	if command.ParentCommentID.Valid {
		if err := queries.Comment_RecalculateSubcomments(ctx, int64NullableDomainToDb(command.ParentCommentID)); err != nil {
			return AddCommentResult{}, apperror.WrapUnexpectedDBError(err)
		}
	}

	comment, err := c.getByID(ctx, id, command.UserID)
	if err != nil {
		return AddCommentResult{}, apperror.UnexpectedError.New(err.Error())
	}

	return AddCommentResult{Comment: comment}, err
}

// GetList implements CommentsService.
func (c *commentsService) GetList(ctx context.Context, query GetCommentsQuery) (result GetCommentsResult, err error) {
	queries := store.New(c.db)

	result.Cursor = query.Cursor
	result.Total, err = queries.Comment_CountByChapter(ctx, query.ChapterID)
	if err != nil {
		return result, apperror.WrapUnexpectedDBError(err)
	}

	rows, err := queries.Comment_GetByChapter(ctx, store.Comment_GetByChapterParams{
		ChapterID:  query.ChapterID,
		Sort:       string(query.Sort),
		PageLimit:  query.Limit + 1,
		PageOffset: int32(query.Cursor),
	})
	if err != nil {
		return result, apperror.WrapUnexpectedDBError(err)
	}
	hasMore := len(rows) > int(query.Limit)
	if hasMore {
		rows = rows[:query.Limit]
	}
	result.Comments = MapSlice(rows, func(r store.Comment_GetByChapterRow) CommentDto {
		return moderationSafeComment(CommentDto{
			ID: r.ID, Content: r.Content,
			User:      CommentUserDto{ID: uuidDbToDomain(r.UserID), Name: r.UserName, Avatar: getUserAvatar(r.UserName, 84)},
			CreatedAt: timeDbToDomain(r.CreatedAt), UpdatedAt: timeNullableDbToDomain(r.UpdatedAt),
			Subcomments: int(r.Subcomments), Likes: int64(r.Likes), LikesUpdatedAt: timeDbToDomain(r.LikesRecalculatedAt),
		}, r.DeletedAt.Valid)
	})

	if len(result.Comments) == 0 {
		result.NextCursor = 0
	} else {
		if hasMore {
			result.NextCursor = query.Cursor + uint32(len(result.Comments))
		}

		if query.ActorUserID.Valid {
			err = c.fillWithLikedAtData(ctx, queries, result.Comments, query.ActorUserID.UUID)
			if err != nil {
				return
			}
		}
	}

	return
}

func (c *commentsService) GetReplies(ctx context.Context, query GetCommentRepliesQuery) (result GetCommentRepliesResult, err error) {
	queries := store.New(c.db)

	result.Cursor = query.Cursor

	if query.Cursor == 0 {
		var rows []store.Comment_GetChildCommentsRow
		rows, err = queries.Comment_GetChildComments(ctx, store.Comment_GetChildCommentsParams{
			ParentID: query.CommentID,
			Limit:    query.Limit,
		})
		if err != nil {
			err = apperror.WrapUnexpectedDBError(err)
			return
		}
		result.Comments = MapSlice(rows, func(r store.Comment_GetChildCommentsRow) CommentDto {
			return moderationSafeComment(CommentDto{
				ID:             r.ID,
				Content:        r.Content,
				User:           CommentUserDto{ID: uuidDbToDomain(r.UserID), Name: r.UserName, Avatar: getUserAvatar(r.UserName, 84)},
				CreatedAt:      timeDbToDomain(r.CreatedAt),
				UpdatedAt:      timeNullableDbToDomain(r.UpdatedAt),
				Subcomments:    int(r.Subcomments),
				Likes:          int64(r.Likes),
				LikesUpdatedAt: timeDbToDomain(r.LikesRecalculatedAt),
			}, r.DeletedAt.Valid)
		})
	} else {
		ts := time.Unix(int64(query.Cursor), 0)

		var rows []store.Comment_GetChildCommentsAfterRow
		rows, err = queries.Comment_GetChildCommentsAfter(ctx, store.Comment_GetChildCommentsAfterParams{
			ParentID:  query.CommentID,
			Limit:     query.Limit,
			CreatedAt: timeToTimestamptz(ts),
		})
		if err != nil {
			err = apperror.WrapUnexpectedDBError(err)
			return
		}
		result.Comments = MapSlice(rows, func(r store.Comment_GetChildCommentsAfterRow) CommentDto {
			return moderationSafeComment(CommentDto{
				ID:          r.ID,
				Content:     r.Content,
				User:        CommentUserDto{ID: uuidDbToDomain(r.UserID), Name: r.UserName, Avatar: getUserAvatar(r.UserName, 84)},
				CreatedAt:   timeDbToDomain(r.CreatedAt),
				UpdatedAt:   timeNullableDbToDomain(r.UpdatedAt),
				Subcomments: int(r.Subcomments),
			}, r.DeletedAt.Valid)
		})
	}

	if len(result.Comments) == 0 {
		result.NextCursor = 0
	} else {
		unixTs := result.Comments[len(result.Comments)-1].CreatedAt.Unix()
		result.NextCursor = uint32(unixTs)

		if query.ActorUserID.Valid {
			err = c.fillWithLikedAtData(ctx, queries, result.Comments, query.ActorUserID.UUID)
			if err != nil {
				return
			}
		}
	}

	return
}

func (c *commentsService) fillWithLikedAtData(ctx context.Context, queries *store.Queries, comments []CommentDto, userID uuid.UUID) error {
	commentIds := make([]int64, len(comments))
	for i := range comments {
		commentIds[i] = comments[i].ID
	}

	likedComments, err := queries.Comment_GetLikedComments(ctx, store.Comment_GetLikedCommentsParams{
		UserID: uuidDomainToDb(userID),
		Ids:    commentIds,
	})
	if err != nil {
		err = apperror.WrapUnexpectedDBError(err)
		return nil
	}
	likedCommentsMapping := map[int64]time.Time{}

	for _, likedComment := range likedComments {
		likedCommentsMapping[likedComment.CommentID] = timeDbToDomain(likedComment.LikedAt)
	}

	for i := range comments {
		comment := comments[i]
		if likedAt, ok := likedCommentsMapping[comment.ID]; ok {
			comment.LikedAt = Value(likedAt)
			comments[i] = comment
		}
	}

	return nil
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

	comment, err := c.getByID(ctx, command.ID, command.UserID)
	if err != nil {
		return UpdateCommentResult{}, apperror.UnexpectedError.New(err.Error())
	}

	return UpdateCommentResult{Comment: comment}, err
}

func (c *commentsService) getByID(ctx context.Context, id int64, userID uuid.UUID) (CommentDto, error) {
	queries := store.New(c.db)
	row, err := queries.Comment_GetWithUserByID(ctx, id)
	if err != nil {
		if err == store.ErrNoRows {
			return CommentDto{}, ErrTypeCommentNotFound.New(fmt.Sprintf("comment with id %d not found", id))
		}
		return CommentDto{}, apperror.WrapUnexpectedDBError(err)
	}

	likedComments, err := queries.Comment_GetLikedComments(ctx, store.Comment_GetLikedCommentsParams{
		UserID: uuidDomainToDb(userID),
		Ids:    []int64{id},
	})
	if err != nil {
		return CommentDto{}, apperror.WrapUnexpectedDBError(err)
	}

	dto := CommentDto{
		ID:             row.ID,
		Content:        row.Content,
		User:           CommentUserDto{ID: uuidDbToDomain(row.UserID), Name: row.UserName, Avatar: getUserAvatar(row.UserName, 84)},
		CreatedAt:      timeDbToDomain(row.CreatedAt),
		UpdatedAt:      timeNullableDbToDomain(row.UpdatedAt),
		Subcomments:    int(row.Subcomments),
		Likes:          int64(row.Likes),
		LikesUpdatedAt: timeDbToDomain(row.LikesRecalculatedAt),
	}
	dto = moderationSafeComment(dto, row.DeletedAt.Valid)

	if len(likedComments) > 0 {
		dto.LikedAt = Value(timeDbToDomain(likedComments[0].LikedAt))
	}

	return dto, nil
}

func moderationSafeComment(dto CommentDto, deleted bool) CommentDto {
	dto.Deleted = deleted
	if deleted {
		dto.Content = ""
	}
	return dto
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
