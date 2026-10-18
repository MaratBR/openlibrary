package elasticstore

import (
	"encoding/json"
	"testing"
)

func TestIndexSearchFilters(t *testing.T) {
	min, max := 1.25, 10.0
	req, err := buildSearchRequest(SearchRequest{
		Page: 2, PageSize: 50, Rating: "PG",
		Words:         Range{Min: Int32{Valid: true, Int32: 100}},
		NumberFilters: map[string]*NumberRangeQuery{"popularity": {Gte: &min}, "views": {Lte: &max}},
		DateFilters:   map[string]DateRangeQuery{"createdAt": {Gte: "2026-01-01", Lte: "2026-01-31"}},
		ExcludeUsers:  []string{"excluded-author"},
	})
	if err != nil {
		t.Fatal(err)
	}
	var body SearchReqBody
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.From != 50 || body.Size != 50 {
		t.Fatalf("wrong pagination: %+v", body)
	}
	ranges := map[string]bool{}
	for _, clause := range body.Query.Bool.Must {
		if len(clause.Range) > 1 {
			t.Fatal("OpenSearch range clause contains multiple fields")
		}
		for field := range clause.Range {
			ranges[field] = true
		}
	}
	for _, field := range []string{"words", "popularity", "views", "createdAt"} {
		if !ranges[field] {
			t.Fatalf("missing range: %s", field)
		}
	}
	if len(body.Query.Bool.Must) != 8 {
		t.Fatalf("missing rating, visibility, or author exclusions: %+v", body.Query.Bool.Must)
	}
}
