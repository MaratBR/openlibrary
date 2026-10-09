package siteconfig

import (
	"context"
	"encoding/json"
	"reflect"
)

type SiteConfigEntry json.RawMessage

func (d ConfigData) Equal(other ConfigData) bool {
	return reflect.DeepEqual(d, other)
}

// Store manages the active site configuration and its persistence.
type Store interface {
	Load(ctx context.Context) error
	Get() *ConfigData
	Save(ctx context.Context, force bool) error
}
