package uiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Preload struct {
	As   string
	Href string
}

type Data struct {
	ServerMetadata map[string]any
	Preloads       []Preload
}

func newData() *Data {
	return &Data{
		ServerMetadata: map[string]any{
			"_preload": map[string]any{},
		},
		Preloads: []Preload{},
	}
}

func (s *Data) AddPreloadURL(as, href string) {
	s.Preloads = append(s.Preloads, Preload{
		As:   as,
		Href: href,
	})
}

func (s *Data) AddPreloadedData(
	key string,
	data any,
) {
	var (
		preloadedValues map[string]any
	)

	if preloadedValuesAny, ok := s.ServerMetadata["_preload"]; ok {
		preloadedValues = preloadedValuesAny.(map[string]any)
		preloadedValues[key] = data
	} else {
		preloadedValues = map[string]any{key: data}
		s.ServerMetadata["_preload"] = preloadedValues
	}
}

func getInjectedHTMLSegment(serverData *Data) []byte {
	var buffer bytes.Buffer

	buffer.WriteString("<!-- START: inject server data -->\n")

	for _, preload := range serverData.Preloads {
		buffer.WriteString(fmt.Sprintf("<link rel=\"preload\" as=\"%s\" href=\"%s\" />\n", preload.As, preload.Href))
	}

	if len(serverData.ServerMetadata) > 0 {
		metaJson, err := json.MarshalIndent(serverData.ServerMetadata, "", "\t")
		if err == nil {
			buffer.WriteString("<script type=\"application/javascript\" id=\"server-data\">\nwindow.__server__=" + string(metaJson) + "\n</script>\n")
		}
	}

	buffer.WriteString("<!-- END: inject server data -->\n")

	return buffer.Bytes()
}
