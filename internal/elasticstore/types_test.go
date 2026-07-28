package elasticstore

import (
	"encoding/json"
	"testing"
)

func TestFuzzyNameMatchQueryJSON(t *testing.T) {
	query := Query{
		Match: map[string]MatchQuery{
			"name": {
				Query:     "mistkae",
				Fuzziness: "AUTO",
			},
		},
	}

	got, err := json.Marshal(query)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	want := `{"match":{"name":{"query":"mistkae","fuzziness":"AUTO"}}}`
	if string(got) != want {
		t.Fatalf("json.Marshal() = %s, want %s", got, want)
	}
}
