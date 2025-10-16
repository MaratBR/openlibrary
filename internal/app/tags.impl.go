package app

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/jackc/pgx/v5/pgtype"
)

type tagsService struct {
	db      store.DBTX
	queries *store.Queries
}

// GetTag implements TagsService.
func (t *tagsService) GetTag(ctx context.Context, id int64) (TagDetailsItemDto, error) {
	tag, err := t.queries.GetTag(ctx, id)
	if err != nil {
		if err == store.ErrNoRows {
			return TagDetailsItemDto{}, ErrTagNotFound
		}

		return TagDetailsItemDto{}, wrapUnexpectedDBError(err)
	}

	var synonym Nullable[struct {
		ID   int64
		Name string
	}]

	if tag.SynonymOf.Valid {
		synonym = Value(struct {
			ID   int64
			Name string
		}{
			ID:   tag.SynonymOf.Int64,
			Name: tag.SynonymName.String,
		})
	}

	return TagDetailsItemDto{
		DefinedTagDto: DefinedTagDto{
			ID:          tag.ID,
			Name:        tag.Name,
			Description: tag.Description,
			IsAdult:     tag.IsAdult,
			IsSpoiler:   tag.IsSpoiler,
			Category:    dbTagTypeToTagsCategory(tag.TagType),
		},
		CreatedAt: timeDbToDomain(tag.CreatedAt),
		IsDefault: tag.IsDefault,
		SynonymOf: synonym,
	}, nil
}

func (t *tagsService) GetTagsByIds(ctx context.Context, ids []int64) ([]DefinedTagDto, error) {
	tags, err := t.queries.GetTagsByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	return MapSlice(tags, definedTagToTagDto), nil
}

func (t *tagsService) SearchTags(ctx context.Context, query string) ([]DefinedTagDto, error) {
	query = strings.Trim(query, " \n\t")
	query = strings.ToLower(query)

	var searchTagsLimit int32 = 20

	if query == "" {
		searchTagsLimit = 100
	}

	tags, err := t.queries.SearchDefinedTags(ctx, store.SearchDefinedTagsParams{
		LowercasedName: escapeSqlLikeValue(query) + "%",
		Limit:          searchTagsLimit,
	})
	if err != nil {
		return nil, err
	}

	return MapSlice(tags, definedTagToTagDto), nil
}

func (t *tagsService) FindParentTagIds(ctx context.Context, ids []int64) (r BookTags, err error) {
	if len(ids) == 0 {
		r.ParentTagIds = []int64{}
		r.TagIds = []int64{}
		return
	}

	tags, err := t.queries.GetTagsByIds(ctx, ids)
	if err != nil {
		return
	}
	r.ParentTagIds = make([]int64, len(tags))
	r.TagIds = make([]int64, len(tags))

	for i, tag := range tags {
		r.TagIds[i] = tag.ID
		if tag.SynonymOf.Valid {
			r.ParentTagIds[i] = tag.SynonymOf.Int64
		} else {
			r.ParentTagIds[i] = tag.ID
		}
	}

	return
}

// CreateTags implements TagsService.
func (t *tagsService) CreateTags(ctx context.Context, cmd CreateTagsCommand) ([]DefinedTagDto, error) {
	tags := make([]tagImportRow, len(cmd.Tags))
	tagNames := make([]string, len(cmd.Tags))

	for i, tag := range cmd.Tags {
		tagNames[i] = tag.Name
		tags[i] = tagImportRow{
			ID:          GenID(),
			Name:        tag.Name,
			IsAdult:     tag.IsAdult,
			IsSpoiler:   tag.IsSpoiler,
			TagType:     tagsCategoryToDbTagType(tag.Category),
			Description: "",
			SynonymOf:   pgtype.Int8{Valid: false},
			CreatedAt:   time.Now(),
		}
	}

	err := importTags(ctx, t.queries, tags)
	if err != nil {
		return nil, err
	}

	tagsDto, err := t.queries.GetTagsByName(ctx, tagNames)
	return MapSlice(tagsDto, definedTagToTagDto), err
}

// List implements TagsService.
func (t *tagsService) List(ctx context.Context, query ListTagsQuery) (ListTagsResult, error) {
	limit := min(1000, max(query.PageSize, 1))
	offset := (max(1, query.Page) - 1) * limit

	dbQuery := store.ListTagsQuery{
		Query:          query.SearchQuery,
		Limit:          uint(limit),
		Offset:         uint(offset),
		OnlyParentTags: query.OnlyParentTags,
		OnlyAdultTags:  query.OnlyAdultTags,
	}
	tags, err := store.ListTags(ctx, t.db, dbQuery)
	if err != nil {
		return ListTagsResult{}, wrapUnexpectedDBError(err)
	}

	count, err := store.CountTags(ctx, t.db, dbQuery)
	totalPages := uint32(math.Ceil(float64(count) / float64(limit)))

	tagDtos := MapSlice(tags, func(t store.TagRow) TagDetailsItemDto {
		var synonym Nullable[struct {
			ID   int64
			Name string
		}]

		if t.SynonymOf.Valid && t.SynonymOfName.Valid {
			synonym = Value(struct {
				ID   int64
				Name string
			}{
				ID:   t.SynonymOf.Int64,
				Name: t.SynonymOfName.String,
			})
		}

		return TagDetailsItemDto{
			DefinedTagDto: DefinedTagDto{
				ID:          t.ID,
				Name:        t.Name,
				Description: t.Description,
				IsAdult:     t.IsAdult,
				IsSpoiler:   t.IsSpoiler,
				Category:    dbTagTypeToTagsCategory(t.TagType),
			},
			IsDefault: t.IsDefault,
			CreatedAt: t.CreatedAt,
			SynonymOf: synonym,
		}
	})
	return ListTagsResult{
		Tags:       tagDtos,
		TotalPages: totalPages,
		Page:       query.Page,
	}, nil
}

func (t *tagsService) UpdateTag(ctx context.Context, command UpdateTagCommand) error {
	if err := validateTagName(command.Name); err != nil {
		return err
	}

	if err := validateTagDescription(command.Description); err != nil {
		return err
	}

	err := t.queries.UpdateTag(ctx, store.UpdateTagParams{
		ID:          command.ID,
		Name:        command.Name,
		Description: command.Description,
		TagType:     tagsCategoryToDbTagType(command.Type),
		IsAdult:     command.IsAdult,
		IsSpoiler:   command.IsSpoiler,
		SynonymOf:   int64NullableDomainToDb(command.SynonymOfTagID),
	})
	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	return nil
}

func escapeSqlLikeValue(v string) string {
	v = strings.ReplaceAll(v, "\\", "\\\\")
	v = strings.ReplaceAll(v, "_", "\\_")
	v = strings.ReplaceAll(v, "%", "\\%")
	return v
}

func NewTagsService(db store.DBTX) TagsService {
	return &tagsService{
		db:      db,
		queries: store.New(db),
	}
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

func validateTagName(name string) error {
	if name == "" {
		return ValidationError.New("tag name cannot be empty")
	}

	if len(name) > 50 {
		return ValidationError.New("tag name cannot be larger than 50 characters")
	}

	return nil
}

func validateTagDescription(value string) error {
	if len(value) > 500 {
		return ValidationError.New("tag description cannot be larger than 500 characters")
	}

	return nil
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
