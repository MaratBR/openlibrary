package analytics

import (
	"encoding/json"
	"testing"
)

func TestPopularityWeightsJSON(t *testing.T) {
	weights := NewPopularityWeights(map[string]float64{
		MetricViews:       1,
		MetricImpressions: 0.02,
	})

	var got map[string]float64
	if err := json.Unmarshal(weights.JSON(), &got); err != nil {
		t.Fatalf("JSON() returned invalid JSON: %v", err)
	}
	if got[MetricViews] != 1 || got[MetricImpressions] != 0.02 || len(got) != 2 {
		t.Errorf("JSON() decoded to %#v", got)
	}
}
