package webinfra

import "time"

// ServerRuntime is populated when the HTTP listener starts, before serving requests.
// StartedAt remains immutable for the lifetime of that server run.
type ServerRuntime struct {
	StartedAt time.Time
}
