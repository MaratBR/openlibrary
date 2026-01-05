package app

import (
	"context"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/gofrs/uuid"
)

var (
	ErrTagNotFound = apperror.AppErrors.NewType("tag404", apperror.ErrTraitEntityNotFound).New("tag with this ID cannot be found")
)

type DefinedTagDto struct {
	ID          int64        `json:"id,string"`
	Name        string       `json:"name"`
	Description string       `json:"desc"`
	IsAdult     bool         `json:"adult"`
	IsSpoiler   bool         `json:"spoiler"`
	Category    TagsCategory `json:"cat"`
}

type TagDetailsItemDto struct {
	DefinedTagDto

	SynonymOf Nullable[struct {
		ID   int64
		Name string
	}]
	CreatedAt time.Time
	IsDefault bool
}

type BookTags struct {
	ParentTagIds []int64
	TagIds       []int64
}

type TagDescriptor struct {
	Name      string
	Category  TagsCategory
	IsAdult   bool
	IsWarning bool
	IsSpoiler bool
}

type CreateTagsCommand struct {
	Tags []TagDescriptor
}

type ListTagsQuery struct {
	SearchQuery    string
	Page           uint32
	PageSize       uint32
	OnlyParentTags bool
	OnlyAdultTags  bool
}

type ListTagsResult struct {
	Tags       []TagDetailsItemDto
	Page       uint32
	TotalPages uint32
}

type UpdateTagCommand struct {
	ID             int64
	Name           string
	Description    string
	IsAdult        bool
	IsSpoiler      bool
	SynonymOfTagID Nullable[int64]
	Type           TagsCategory
	UserID         uuid.UUID
}

type TagsService interface {
	GetTag(ctx context.Context, id int64) (TagDetailsItemDto, error)
	GetTagsByIds(ctx context.Context, ids []int64) ([]DefinedTagDto, error)
	SearchTags(ctx context.Context, query string) ([]DefinedTagDto, error)
	FindParentTagIds(ctx context.Context, names []int64) (BookTags, error)
	CreateTags(ctx context.Context, cmd CreateTagsCommand) ([]DefinedTagDto, error)
	List(ctx context.Context, query ListTagsQuery) (ListTagsResult, error)
	UpdateTag(ctx context.Context, command UpdateTagCommand) error
}
