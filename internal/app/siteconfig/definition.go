package siteconfig

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// SettingsDefinition binds an independently versioned row to a snapshot field.
// Published migrations must only be appended: index i upgrades version i to i+1.
type SettingsDefinition interface {
	Key() string
	Version() int
	Load(*ConfigData, json.RawMessage) (bool, error)
	Encode(ConfigData) (json.RawMessage, error)
}

type Definition[T any] struct {
	Name       string
	Field      func(*ConfigData) *T
	Migrations []func(*T) error
}

type settingsRow struct {
	Version int
	Data    json.RawMessage
}

func (d Definition[T]) Key() string  { return d.Name }
func (d Definition[T]) Version() int { return len(d.Migrations) }

func (d Definition[T]) Load(snapshot *ConfigData, raw json.RawMessage) (bool, error) {
	row := settingsRow{}
	value := *d.Field(snapshot)
	if len(raw) > 0 {
		if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			return false, fmt.Errorf("null settings row")
		}
		if err := json.Unmarshal(raw, &row); err != nil {
			return false, err
		}
		if len(row.Data) == 0 || bytes.Equal(bytes.TrimSpace(row.Data), []byte("null")) {
			return false, fmt.Errorf("missing settings data")
		}
		if err := json.Unmarshal(row.Data, &value); err != nil {
			return false, err
		}
	}
	if row.Version < 0 || row.Version > d.Version() {
		return false, fmt.Errorf("unsupported version %d (latest %d)", row.Version, d.Version())
	}
	for i := row.Version; i < d.Version(); i++ {
		if err := d.Migrations[i](&value); err != nil {
			return false, fmt.Errorf("migration %d: %w", i, err)
		}
	}
	*d.Field(snapshot) = value
	return len(raw) == 0 || row.Version != d.Version(), nil
}

func (d Definition[T]) Encode(snapshot ConfigData) (json.RawMessage, error) {
	data, err := json.Marshal(d.Field(&snapshot))
	if err != nil {
		return nil, err
	}
	return json.Marshal(settingsRow{Version: d.Version(), Data: data})
}
