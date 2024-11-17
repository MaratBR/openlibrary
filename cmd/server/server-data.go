package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MaratBR/openlibrary/internal/app"
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

func getInjectedHTMLSegment(r *http.Request) []byte {
	serverData, ok := getServerData(r)
	if !ok {
		return nil
	}

	var buffer bytes.Buffer

	buffer.WriteString(`<!-- START: inject server data -->`)

	for _, preload := range serverData.Preloads {
		buffer.WriteString(fmt.Sprintf(`<link rel="preload" as="%s" href="%s">`, preload.As, preload.Href))
	}

	metaJson, err := json.Marshal(serverData.ServerMetadata)
	if err != nil {
	} else {
		buffer.WriteString(`<script type="application/javascript" id="server-data">window.SERVER_DATA=` + string(metaJson) + `</script>`)
	}

	buffer.WriteString(`<!-- END: inject server data -->`)

	return buffer.Bytes()
}

func getServerData(r *http.Request) (*serverData, bool) {
	value := r.Context().Value(serverDataKey)
	if value == nil {
		return nil, false
	}

	serverData := value.(*serverData)
	return serverData, true
}

func getDefaultServerMetadata(r *http.Request) map[string]any {
	now := time.Now()
	tz, offset := now.Zone()
	meta := map[string]any{
		"ts":    now.Unix(),
		"tz":    tz,
		"tzoff": offset,
		"id":    app.GenID(),
	}

	sessionInfo, ok := getSession(r)
	if ok {
		meta["session"] = map[string]any{
			"id":           sessionInfo.UserID,
			"username":     sessionInfo.UserName,
			"createdAt":    sessionInfo.CreatedAt,
			"expiresAt":    sessionInfo.ExpiresAt,
			"userJoinedAt": sessionInfo.UserJoinedAt,
		}
	}

	return meta
}

func injectServerData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverData := &serverData{
			ServerMetadata: getDefaultServerMetadata(r),
		}
		newContext := context.WithValue(r.Context(), serverDataKey, serverData)
		r = r.WithContext(newContext)
		next.ServeHTTP(w, r)
	})
}
