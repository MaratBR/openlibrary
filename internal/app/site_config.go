package app

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
)

type SiteConfigEntry json.RawMessage

type SiteConfig struct {
	mx          sync.Mutex
	changedKeys map[string]struct{}
	db          DB
	lastFetched time.Time
	map_        map[string]SiteConfigEntry
}

func (c *SiteConfig) get(key string, dest any) (error, bool) {
	bytes, ok := c.map_[key]
	if !ok || bytes == nil {
		return nil, false
	}
	err := json.Unmarshal(bytes, dest)
	if err != nil {
		slog.Error("failed to unmarshal site config", "key", key, "json", string(bytes))
		return err, false
	}

	return nil, true
}

func (c *SiteConfig) set(key string, v any) error {
	json, err := json.Marshal(v)
	if err != nil {
		return err
	}
	c.mx.Lock()
	defer c.mx.Unlock()
	c.map_[key] = json
	if c.changedKeys == nil {
		c.changedKeys = map[string]struct{}{}
	}
	c.changedKeys[key] = struct{}{}
	return nil
}

func NewSiteConfig(db DB) *SiteConfig {
	cfg := &SiteConfig{
		map_: map[string]SiteConfigEntry{},
		db:   db,
	}
	return cfg
}

func (s *SiteConfig) Load(ctx context.Context) error {
	if s.shouldRefresh() {
		return s.fetch(ctx)
	}
	return nil
}

func (s *SiteConfig) Save(ctx context.Context) error {
	if len(s.changedKeys) == 0 {
		return nil
	}

	for changedKey := range s.changedKeys {
		queries := store.New(s.db)
		value, ok := s.map_[changedKey]
		if !ok {
			continue
		}
		err := queries.SiteConfig_Set(ctx, store.SiteConfig_SetParams{
			Key:   changedKey,
			Value: value,
		})
		if err != nil {
			slog.Error("failed to update site settings", "err", err)
		}
	}

	return s.Load(ctx)
}

func (s *SiteConfig) shouldRefresh() bool {
	return s.lastFetched == time.Time{} || s.lastFetched.Before(time.Now().Add(-time.Second*30))
}

func (s *SiteConfig) fetch(ctx context.Context) error {
	queries := store.New(s.db)
	cfg, err := queries.SiteConfig_All(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	s.mx.Lock()
	defer s.mx.Unlock()

	s.changedKeys = nil
	s.map_ = map[string]SiteConfigEntry{}
	for _, row := range cfg {
		s.map_[row.Key] = SiteConfigEntry(row.Value)
	}

	return nil
}
