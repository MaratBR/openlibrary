package app

import (
	"context"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

var (
	CommentErrors                = AppErrors.NewSubNamespace("comment")
	ErrTypeCommentNotFound       = CommentErrors.NewType("not_found", ErrTraitEntityNotFound)
	ErrTypeCommentContentInvalid = CommentErrors.NewType("invalid_content", ErrTraitValidationError)
	ErrCommentContentEmpty       = ErrTypeCommentContentInvalid.New("comment content is empty")
	ErrCommentContentTooLarge    = ErrTypeCommentContentInvalid.New("comment content is too large")
)

type GetCommentsQuery struct {
	ActorUserID uuid.NullUUID

	ChapterID int64
	Limit     int32
	Cursor    uint32
}

type CommentDto struct {
	ID        int64               `json:"id,string"`
	Content   string              `json:"content"`
	User      CommentUserDto      `json:"user"`
	CreatedAt time.Time           `json:"createdAt"`
	UpdatedAt Nullable[time.Time] `json:"updatedAt"`
}

type CommentUserDto struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Avatar string    `json:"avatar"`
}

type GetCommentsResult struct {
	Cursor     uint32
	NextCursor uint32
	Comments   []CommentDto
}

type AddCommentCommand struct {
	UserID          uuid.UUID
	ChapterID       int64
	ParentCommentID Nullable[int64]
	Content         string
}

func (c *AddCommentCommand) Validate() error {
	if c.Content == "" {
		return ErrCommentContentEmpty
	}

	c.Content = strings.Trim(c.Content, " \n\t")

	if c.Content == "" {
		return ErrCommentContentEmpty
	}

	if len(c.Content) > 2000 {
		return ErrCommentContentTooLarge
	}

	return nil
}

type AddCommentResult struct {
	Comment CommentDto
}

type UpdateCommentCommand struct {
	ID      int64
	Content string
	UserID  uuid.UUID
}

type UpdateCommentResult struct {
	Comment CommentDto
}

type CommentsService interface {
	GetList(ctx context.Context, query GetCommentsQuery) (GetCommentsResult, error)

	AddComment(ctx context.Context, command AddCommentCommand) (AddCommentResult, error)
	UpdateComment(ctx context.Context, command UpdateCommentCommand) (UpdateCommentResult, error)
}
