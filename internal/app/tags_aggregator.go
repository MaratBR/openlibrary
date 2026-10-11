package app

import "context"

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
