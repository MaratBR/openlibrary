package elasticstore

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

type Int32 struct {
	Valid bool
	Int32 int32
}

type Range struct {
	Min Int32
	Max Int32
}

func (r Range) Has() bool {
	return r.Max.Valid || r.Min.Valid
}

type SearchRequest struct {
	Query           string
	IncludeUsers    []string
	ExcludeUsers    []string
	IncludeTags     []int64
	ExcludeTags     []int64
	Words           Range
	Chapters        Range
	WordsPerChapter Range
	Page            int32
	PageSize        int32
	IncludeHidden   bool
}

type SearchRow struct {
	BookIndex
	ID int64
}

type SearchResult struct {
	Hits   []SearchRow
	TookMS int64
	Total  int64
}

func createNumberRangeQuery(rng Range) *types.NumberRangeQuery {
	rangeQuery := types.NewNumberRangeQuery()

	if rng.Min.Valid {
		value := types.Float64(float64(rng.Min.Int32))
		rangeQuery.Gte = &value
	}
	if rng.Max.Valid {
		value := types.Float64(float64(rng.Max.Int32))
		rangeQuery.Lte = &value
	}

	return rangeQuery
}

func createAndTermQuery(term string, values []types.FieldValue) types.Query {
	subQueries := make([]types.Query, len(values))
	for i := 0; i < len(values); i++ {
		subQueries[i] = types.Query{
			Term: map[string]types.TermQuery{
				term: {
					Value: values[i],
				},
			},
		}
	}

	return types.Query{
		Bool: &types.BoolQuery{
			Must: subQueries,
		},
	}
}

func createNegativeAndTermQuery(term string, values []types.FieldValue) types.Query {
	subQueries := make([]types.Query, len(values))
	for i := 0; i < len(values); i++ {
		subQueries[i] = types.Query{
			Term: map[string]types.TermQuery{
				term: {
					Value: values[i],
				},
			},
		}
	}

	return types.Query{
		Bool: &types.BoolQuery{
			MustNot: subQueries,
		},
	}
}

func Search(
	ctx context.Context,
	client *elasticsearch.TypedClient,
	req SearchRequest,
) (SearchResult, error) {
	must := []types.Query{
		{
			Term: map[string]types.TermQuery{
				"isTrashed": {
					Value: true,
				},
			},
		},
		{
			Term: map[string]types.TermQuery{
				"isPubliclyVisible": {
					Value: true,
				},
			},
		},
	}

	if req.Query != "" {
		must = append(must, types.Query{
			QueryString: &types.QueryStringQuery{
				Query: req.Query,
			},
		})
	}

	rangeQueries := map[string]types.RangeQuery{}

	if req.Words.Has() {
		rangeQueries["words"] = createNumberRangeQuery(req.Words)
	}
	if req.Chapters.Has() {
		rangeQueries["chapters"] = createNumberRangeQuery(req.Chapters)
	}
	if req.WordsPerChapter.Has() {
		rangeQueries["wordsPerChapter"] = createNumberRangeQuery(req.WordsPerChapter)
	}

	if len(rangeQueries) > 0 {
		must = append(must, types.Query{
			Range: rangeQueries,
		})
	}

	if len(req.IncludeUsers) > 0 {
		ids := make([]types.FieldValue, len(req.IncludeUsers))
		for i := 0; i < len(req.IncludeUsers); i++ {
			ids[i] = req.IncludeUsers[i]
		}
		must = append(must, createAndTermQuery("authorId", ids))
	}

	if len(req.ExcludeUsers) > 0 {
		ids := make([]types.FieldValue, len(req.IncludeUsers))
		for i := 0; i < len(req.ExcludeUsers); i++ {
			ids[i] = req.ExcludeUsers[i]
		}
		must = append(must, createNegativeAndTermQuery("authorId", ids))
	}

	if len(req.IncludeTags) > 0 {
		ids := make([]types.FieldValue, len(req.IncludeTags))
		for i := 0; i < len(req.IncludeTags); i++ {
			ids[i] = req.IncludeTags[i]
		}
		must = append(must, createAndTermQuery("tags", ids))
	}

	if len(req.ExcludeTags) > 0 {
		ids := make([]types.FieldValue, len(req.ExcludeTags))
		for i := 0; i < len(req.ExcludeTags); i++ {
			ids[i] = req.ExcludeTags[i]
		}
		must = append(must, createNegativeAndTermQuery("tags", ids))
	}

	size := int(req.PageSize)
	from := int((req.Page - 1) * req.PageSize)
	if from < 0 {
		from = 0
	}

	esReq := &search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Must: must,
			},
		},
		From: &from,
		Size: &size,
	}
	resp, err := client.Search().Request(esReq).Do(ctx)

	if err != nil {
		onSearchFail(esReq, err)
		return SearchResult{}, err
	}

	results := make([]SearchRow, 0, len(resp.Hits.Hits))

	if len(resp.Hits.Hits) == 0 {
		onSearchEmpty(esReq)
	}

	for i := 0; i < len(resp.Hits.Hits); i++ {
		hit := resp.Hits.Hits[i]
		var result SearchRow

		if hit.Id_ == nil {
			continue
		}

		id, err := strconv.ParseInt(*hit.Id_, 10, 64)
		if err != nil {
			continue
		}

		err = json.Unmarshal(hit.Source_, &result)

		if err != nil {
			continue
		}

		result.ID = id

		results = append(results, result)
	}

	return SearchResult{
		Hits:   results,
		TookMS: resp.Took,
		Total:  resp.Hits.Total.Value,
	}, nil
}

func onSearchFail(req *search.Request, err error) {
	jsonText, err_ := json.Marshal(req)
	if err_ != nil {
		slog.Error("failed to serialize elastic request to JSON", "err", err_)
		slog.Error("failed to execute elastic request", "err", err)
	} else {
		slog.Error("failed to execute elastic request", "err", err, "req", string(jsonText))
	}
}

func onSearchEmpty(req *search.Request) {
	jsonText, err := json.Marshal(req)
	if err == nil {
		slog.Debug("empty search result", "req", string(jsonText))
	}
}
