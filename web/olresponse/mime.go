package olresponse

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// PreferredMimeTypeIsJSON checks if the client prefers JSON mime types more than other types
func PreferredMimeTypeIsJSON(r *http.Request) bool {
	// Parse the Accept header
	acceptHeader := r.Header.Get("Accept")
	if acceptHeader == "" {
		return false
	}

	// Split the Accept header into individual mime types with their quality values
	mimeTypes := strings.Split(acceptHeader, ",")

	// Create a slice to store parsed mime types and their weights
	type MimeWeight struct {
		mime   string
		weight float64
	}
	parsedMimes := make([]MimeWeight, 0, len(mimeTypes))

	// Parse each mime type and its quality value
	for _, mimeType := range mimeTypes {
		parts := strings.Split(strings.TrimSpace(mimeType), ";")
		mime := strings.TrimSpace(parts[0])
		weight := 1.0 // Default quality value

		// Check for explicit quality value
		if len(parts) > 1 {
			qParts := strings.Split(parts[1], "=")
			if len(qParts) == 2 && strings.TrimSpace(qParts[0]) == "q" {
				// Parse the quality value
				var q float64
				_, err := fmt.Sscanf(strings.TrimSpace(qParts[1]), "%f", &q)
				if err == nil {
					weight = q
				}
			}
		}

		// Normalize mime type
		parsedMimes = append(parsedMimes, MimeWeight{
			mime:   mime,
			weight: weight,
		})
	}

	// Sort mime types by weight in descending order
	sort.Slice(parsedMimes, func(i, j int) bool {
		return parsedMimes[i].weight > parsedMimes[j].weight
	})

	// Check if the top mime type is JSON
	if len(parsedMimes) > 0 {
		topMime := parsedMimes[0].mime
		return topMime == "application/json" || topMime == "text/json"
	}

	return false
}
