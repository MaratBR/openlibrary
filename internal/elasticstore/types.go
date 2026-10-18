package elasticstore

type FieldValue any

type Query struct {
	Term        map[string]TermQuery  `json:"term,omitempty"`
	Match       map[string]MatchQuery `json:"match,omitempty"`
	Bool        *BoolQuery            `json:"bool,omitempty"`
	QueryString *QueryStringQuery     `json:"query_string,omitempty"`
	Range       map[string]RangeQuery `json:"range,omitempty"`
}

type HighlightRequest struct {
	PreTags  []string                  `json:"pre_tags"`
	PostTags []string                  `json:"post_tags"`
	Encoder  string                    `json:"encoder"`
	Fields   map[string]HighlightField `json:"fields"`
}

type HighlightField struct {
	NumberOfFragments int `json:"number_of_fragments"`
}

type SearchReqBody struct {
	Highlight *HighlightRequest       `json:"highlight,omitempty"`
	Query     Query                   `json:"query"`
	From      int                     `json:"from,omitempty"`
	Size      int                     `json:"size,omitempty"`
	Sort      []map[string]SortOption `json:"sort,omitempty"`
}

type SortOption struct {
	Order string `json:"order"`
}

type TermQuery struct {
	Value FieldValue `json:"value"`
}

type MatchQuery struct {
	Query     string  `json:"query"`
	Fuzziness string  `json:"fuzziness,omitempty"`
	Boost     float64 `json:"boost,omitempty"`
	Operator  string  `json:"operator,omitempty"`
}

type BoolQuery struct {
	Must    []Query `json:"must,omitempty"`
	MustNot []Query `json:"must_not,omitempty"`
	Should  []Query `json:"should,omitempty"`
}

type QueryStringQuery struct {
	Query string `json:"query"`
}

type NumberRangeQuery struct {
	Gt  *float64 `json:"gt,omitempty"`
	Gte *float64 `json:"gte,omitempty"`
	Lt  *float64 `json:"lt,omitempty"`
	Lte *float64 `json:"lte,omitempty"`
}

type RangeQuery any
