package app

import (
	"context"
	"strings"

	"github.com/MaratBR/openlibrary/internal/store"
)

func dbTagTypeToTagsCategory(t store.TagType) TagsCategory {
	switch t {
	case store.TagTypeFreeform:
		return TagsCategoryOther
	case store.TagTypeWarning:
		return TagsCategoryWarning
	case store.TagTypeFandom:
		return TagsCategoryFandom
	case store.TagTypeRel:
		return TagsCategoryRelationship
	case store.TagTypeReltype:
		return TagsCategoryRelationshipType
	default:
		return TagsCategoryOther
	}
}

func tagsCategoryToDbTagType(cat TagsCategory) store.TagType {
	switch cat {
	case TagsCategoryOther:
		return store.TagTypeFreeform
	case TagsCategoryWarning:
		return store.TagTypeWarning
	case TagsCategoryFandom:
		return store.TagTypeFandom
	case TagsCategoryRelationship:
		return store.TagTypeRel
	case TagsCategoryRelationshipType:
		return store.TagTypeReltype
	default:
		return store.TagTypeFreeform
	}
}

func tagsCategoryName(cat TagsCategory) string {
	switch cat {
	case TagsCategoryOther:
		return "other"
	case TagsCategoryWarning:
		return "warning"
	case TagsCategoryFandom:
		return "fandom"
	case TagsCategoryRelationship:
		return "rel"
	case TagsCategoryRelationshipType:
		return "reltype"
	default:
		return "unknown"
	}
}

func tagsCategoryFromName(name string) TagsCategory {
	switch name {
	case "other":
		return TagsCategoryOther
	case "warning":
		return TagsCategoryWarning
	case "fandom":
		return TagsCategoryFandom
	case "rel":
		return TagsCategoryRelationship
	case "reltype":
		return TagsCategoryRelationshipType
	default:
		return TagsCategoryOther
	}
}

func definedTagToTagDto(t store.DefinedTag) DefinedTagDto {
	return DefinedTagDto{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		IsAdult:     t.IsAdult,
		IsSpoiler:   t.IsSpoiler,
		Category:    dbTagTypeToTagsCategory(t.TagType),
	}
}

type tagsService struct {
	queries *store.Queries
}

func NewTagsService(db store.DBTX) TagsService {
	return &tagsService{
		queries: store.New(db),
	}
}

func (t *tagsService) FindBookTags(ctx context.Context, names []string) (BookTags, error) {
	if len(names) == 0 {
		return BookTags{TagIds: []int64{}, ParentTagIds: []int64{}}, nil
	}

	tags, err := t.queries.GetTagsByName(ctx, names)
	if err != nil {
		return BookTags{}, err
	}
	parentTagIds := make([]int64, len(tags))
	for i, tag := range tags {
		if tag.SynonymOf.Valid {
			parentTagIds[i] = tag.SynonymOf.Int64
		} else {
			parentTagIds[i] = tag.ID
		}
	}

	tagIDs := mapSlice(tags, func(tag store.DefinedTag) int64 { return tag.ID })

	return BookTags{
		TagIds:       tagIDs,
		ParentTagIds: parentTagIds,
	}, nil
}

func (t *tagsService) GetTagsByIds(ctx context.Context, ids []int64) ([]DefinedTagDto, error) {
	tags, err := t.queries.GetTagsByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	return mapSlice(tags, definedTagToTagDto), nil
}

const (
	searchTagsLimit = 20
)

func (t *tagsService) SearchTags(ctx context.Context, query string) ([]DefinedTagDto, error) {
	query = strings.Trim(query, " \n\t")
	query = strings.ToLower(query)
	tags, err := t.queries.SearchDefinedTags(ctx, store.SearchDefinedTagsParams{
		LowercasedName: escapeSqlLikeValue(query) + "%",
		Limit:          searchTagsLimit,
	})
	if err != nil {
		return nil, err
	}

	return mapSlice(tags, definedTagToTagDto), nil
}

func escapeSqlLikeValue(v string) string {
	v = strings.ReplaceAll(v, "\\", "\\\\")
	v = strings.ReplaceAll(v, "_", "\\_")
	v = strings.ReplaceAll(v, "%", "\\%")
	return v
}

type tagsAggregator struct {
	service  TagsService
	tags     map[int64]struct{}
	bookTags map[int64][]int64
}

func newTagsAggregator(service TagsService) tagsAggregator {
	return tagsAggregator{service: service, tags: make(map[int64]struct{}), bookTags: make(map[int64][]int64)}
}

func (agg *tagsAggregator) Add(bookID int64, tagIDs []int64) {
	if len(tagIDs) == 0 {
		return
	}
	agg.bookTags[bookID] = tagIDs
	for _, id := range tagIDs {
		agg.tags[id] = struct{}{}
	}
}

func (agg *tagsAggregator) Fetch(ctx context.Context) (map[int64]DefinedTagDto, error) {
	ids := mapKeys(agg.tags)
	tags, err := agg.service.GetTagsByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	tagsMap := make(map[int64]DefinedTagDto, len(tags))
	for _, t := range tags {
		tagsMap[t.ID] = t
	}
	return tagsMap, nil
}

func (agg *tagsAggregator) BookTags(bookID int64) []int64 {
	ids, _ := agg.bookTags[bookID]
	return ids
}
