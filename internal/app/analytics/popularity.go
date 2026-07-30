package analytics

import (
	"encoding/json"
)

type PopularityWeights struct {
	weights map[string]float64
	json    []byte
}

func NewPopularityWeights(weights map[string]float64) *PopularityWeights {
	return &PopularityWeights{weights: weights}
}

func (w *PopularityWeights) JSON() []byte {
	if w.json != nil {
		return w.json
	}

	var err error
	w.json, err = json.Marshal(w.weights)
	if err != nil {
		panic(err)
	}
	return w.json
}

var DefaultPopularityWeightsValues = NewPopularityWeights(map[string]float64{
	MetricViews:        1,
	MetricImpressions:  0.02,
	MetricSearchClicks: 0.25,
})
