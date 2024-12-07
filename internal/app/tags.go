package app

import (
	"context"
	"encoding/json"
	"fmt"
)

type TagsCategory uint8

const (
	TagsCategoryOther TagsCategory = iota
	TagsCategoryWarning
	TagsCategoryFandom
	TagsCategoryRelationship
	TagsCategoryRelationshipType
)

func (t TagsCategory) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, tagsCategoryName(t))), nil
}

func (t *TagsCategory) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*t = tagsCategoryFromName(s)
	return nil
}

type DefinedTagDto struct {
	ID          int64        `json:"id,string"`
	Name        string       `json:"name"`
	Description string       `json:"desc"`
	IsAdult     bool         `json:"adult"`
	IsSpoiler   bool         `json:"spoiler"`
	Category    TagsCategory `json:"cat"`
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

type TagsService interface {
	GetTagsByIds(ctx context.Context, ids []int64) ([]DefinedTagDto, error)
	SearchTags(ctx context.Context, query string) ([]DefinedTagDto, error)
	FindParentTagIds(ctx context.Context, names []int64) (BookTags, error)

	CreateTags(ctx context.Context, cmd CreateTagsCommand) ([]DefinedTagDto, error)
}
