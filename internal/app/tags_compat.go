package app

// Compatibility with existing callers; new feature code belongs in tags.
import (
	"github.com/MaratBR/openlibrary/internal/app/tags"
	store "github.com/MaratBR/openlibrary/internal/store"
)

var ErrTagNotFound = tags.ErrTagNotFound

type DefinedTagDto = tags.DefinedTagDto
type TagDetailsItemDto = tags.TagDetailsItemDto
type BookTags = tags.BookTags
type TagDescriptor = tags.TagDescriptor
type CreateTagsCommand = tags.CreateTagsCommand
type ListTagsQuery = tags.ListTagsQuery
type ListTagsResult = tags.ListTagsResult
type UpdateTagCommand = tags.UpdateTagCommand
type TagsService = tags.TagsService
type TagsCategory = tags.TagsCategory

const TagsCategoryOther = tags.TagsCategoryOther
const TagsCategoryWarning = tags.TagsCategoryWarning
const TagsCategoryFandom = tags.TagsCategoryFandom
const TagsCategoryRelationship = tags.TagsCategoryRelationship
const TagsCategoryRelationshipType = tags.TagsCategoryRelationshipType
const TagsCategoryGenre = tags.TagsCategoryGenre

var TagsCategoryList = tags.TagsCategoryList

func TagsCategoryFromName(name string) TagsCategory { return tags.TagsCategoryFromName(name) }
func NewTagsService(db store.DBTX) TagsService      { return tags.NewTagsService(db) }
