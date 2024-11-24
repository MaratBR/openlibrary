package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Preload struct {
	As   string
	Href string
}

type serverData struct {
	ServerMetadata map[string]any
	Preloads       []Preload
}

func (s *serverData) AddPreloadedData(
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

type serverDataKeyType struct{}

var serverDataKey serverDataKeyType

func getInjectedHTMLSegment(r *http.Request, serverData *serverData) []byte {
	var buffer bytes.Buffer

	buffer.WriteString(`<!-- START: inject server data -->`)

	for _, preload := range serverData.Preloads {
		buffer.WriteString(fmt.Sprintf(`<link rel="preload" as="%s" href="%s" />`, preload.As, preload.Href))
	}

	if len(serverData.ServerMetadata) > 0 {
		metaJson, err := json.Marshal(serverData.ServerMetadata)
		if err == nil {
			buffer.WriteString(`<script type="application/javascript" id="server-data">window.__server__=` + string(metaJson) + `</script>`)
		}
	}

	buffer.WriteString(`<!-- END: inject server data -->`)

	return buffer.Bytes()
}
