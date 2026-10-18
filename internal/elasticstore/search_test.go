package elasticstore

import (
	"encoding/json"
	"testing"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

func TestSearchRewardsExactWordsAndKeepsFuzzyFallback(t *testing.T) {
	for _, tc := range []struct {
		name  string
		query string
		sort  SearchSort
	}{
		{name: "single word", query: "Runebound"},
		{name: "multiple words", query: "Runebound Professor"},
		{name: "explicit sort", query: "Runebound", sort: SearchSortWords},
		{name: "browse"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request, err := buildSearchRequest(SearchRequest{
				Query: tc.query, Sort: tc.sort, Page: 1, PageSize: 20,
			})
			if err != nil {
				t.Fatal(err)
			}
			var body SearchReqBody
			if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if tc.query == "" {
				if len(body.Query.Bool.Should) != 0 || len(body.Query.Bool.Must) != 2 {
					t.Fatal("browse should only require visibility filters")
				}
				return
			}
			if len(body.Query.Bool.Must) != 3 || len(body.Query.Bool.Should) != 1 {
				t.Fatalf("expected required fuzzy match and optional exact boost: %+v", body.Query.Bool)
			}
			fuzzy := body.Query.Bool.Must[2].Match["name"]
			exact := body.Query.Bool.Should[0].Match["name"]
			if fuzzy.Query != tc.query || fuzzy.Fuzziness != "AUTO" {
				t.Fatalf("missing fuzzy fallback: %+v", fuzzy)
			}
			if exact.Query != tc.query || exact.Fuzziness != "" || exact.Operator != "and" || exact.Boost <= 1 {
				t.Fatalf("expected boosted match on all exact words: %+v", exact)
			}
			if tc.sort != "" && (len(body.Sort) != 1 || body.Sort[0][tc.sort.field()].Order != "desc") {
				t.Fatalf("explicit sort was not preserved: %+v", body.Sort)
			}
		})
	}
}

func TestMapSearchResultSkipsInvalidHits(t *testing.T) {
	var response opensearchapi.SearchResp
	if err := json.Unmarshal([]byte(`{
		"took": 7,
		"hits": {
			"total": {"value": 5, "relation": "eq"},
			"hits": [
				{"_id": "42", "_source": {"name": "Runebound"}},
				{"_id": "", "_source": {"name": "Missing ID"}},
				{"_id": "invalid", "_source": {}},
				{"_id": "43", "_source": {"name": 123}},
				{"_id": "44", "_source": {"name": "Unbound"}}
			]
		}
	}`), &response); err != nil {
		t.Fatal(err)
	}
	result := mapSearchResult(&response)
	if result.TookMS != 7 || result.Total != 5 {
		t.Fatalf("unexpected metadata: %+v", result)
	}
	if len(result.Hits) != 2 || result.Hits[0].ID != 42 || result.Hits[0].Name != "Runebound" || result.Hits[1].ID != 44 {
		t.Fatalf("unexpected mapped hits: %+v", result.Hits)
	}
}

func TestMapSearchResultEmptyHits(t *testing.T) {
	result := mapSearchResult(&opensearchapi.SearchResp{})
	if result.Hits == nil || len(result.Hits) != 0 {
		t.Fatalf("expected non-nil empty hits: %+v", result.Hits)
	}
}

func TestSearchTitleHighlights(t *testing.T) {
	for _, query := range []string{"", "Hobit", "Hobbit Returns"} {
		request, err := buildSearchRequest(SearchRequest{Query: query, Page: 1, PageSize: 20})
		if err != nil {
			t.Fatal(err)
		}
		var body SearchReqBody
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if query == "" {
			if body.Highlight != nil {
				t.Fatal("browse should not request highlights")
			}
			continue
		}
		h := body.Highlight
		if h == nil || h.Encoder != "html" || len(h.PreTags) != 1 || h.PreTags[0] != "<mark>" || len(h.PostTags) != 1 || h.PostTags[0] != "</mark>" {
			t.Fatalf("unexpected highlight configuration: %+v", h)
		}
		field, ok := h.Fields["name"]
		if !ok || field.NumberOfFragments != 0 {
			t.Fatal("expected full title highlighting")
		}
	}
	var response opensearchapi.SearchResp
	if err := json.Unmarshal([]byte(`{"hits":{"hits":[{"_id":"1","_source":{"name":"The Hobbit"},"highlight":{"name":["The <mark>Hobbit</mark>"]}},{"_id":"2","_source":{"name":"Other"}}]}}`), &response); err != nil {
		t.Fatal(err)
	}
	result := mapSearchResult(&response)
	if len(result.Hits) != 2 || result.Hits[0].NameHighlight != "The <mark>Hobbit</mark>" || result.Hits[1].NameHighlight != "" {
		t.Fatalf("unexpected highlights: %+v", result.Hits)
	}
}
