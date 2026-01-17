package elasticstore

type FieldValue any

type Query struct {
	Term        map[string]TermQuery  `json:"term,omitempty"`
	Bool        *BoolQuery            `json:"bool,omitempty"`
	QueryString *QueryStringQuery     `json:"query_string,omitempty"`
	Range       map[string]RangeQuery `json:"range,omitempty"`
}

type SearchReqBody struct {
	Query Query `json:"query"`
	From  int   `json:"from,omitempty"`
	Size  int   `json:"from,omitempty"`
}

type TermQuery struct {
	Value FieldValue `json:"value"`
}

type BoolQuery struct {
	Must    []Query `json:"must,omitempty"`
	MustNot []Query `json:"must_not,omitempty"`
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
