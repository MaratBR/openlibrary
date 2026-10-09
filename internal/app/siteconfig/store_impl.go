package siteconfig

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/MaratBR/openlibrary/internal/app/dal"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/knadh/koanf/v2"
)

type storeImpl struct {
	db       dal.DB
	mu       sync.Mutex
	defaults ConfigData
	active   ConfigData
	saved    json.RawMessage
}

// Get returns the editable active configuration. Callers must coordinate edits
// with Load and Save; Save persists edits made through this pointer.
func (s *storeImpl) Get() *ConfigData { return &s.active }

func (s *storeImpl) Load(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	raw, err := store.New(s.db).SiteConfig_Get(ctx)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("load site configuration: %w", err)
	}
	next := cloneConfig(s.defaults)
	if errors.Is(err, pgx.ErrNoRows) {
		s.active = next
		s.saved = nil
		return nil
	}
	var rows map[string]json.RawMessage
	if err := json.Unmarshal(raw, &rows); err != nil {
		return fmt.Errorf("decode site configuration: %w", err)
	}
	changed := false
	for _, definition := range settingsDefinitions {
		migrated, err := definition.Load(&next, rows[definition.Key()])
		if err != nil {
			return fmt.Errorf("load %s: %w", definition.Key(), err)
		}
		changed = changed || migrated
	}
	encoded, err := encodeSettings(next)
	if err != nil {
		return err
	}
	if changed {
		if err := store.New(s.db).SiteConfig_Set(ctx, encoded); err != nil {
			return fmt.Errorf("persist site configuration migrations: %w", err)
		}
	}
	s.active = next
	s.saved = encoded
	return nil
}

func (s *storeImpl) Save(ctx context.Context, force bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	encoded, err := encodeSettings(s.active)
	if err != nil {
		return err
	}
	if !force && string(encoded) == string(s.saved) {
		return nil
	}
	if err := store.New(s.db).SiteConfig_Set(ctx, encoded); err != nil {
		return fmt.Errorf("save site configuration: %w", err)
	}
	s.saved = encoded
	return nil
}

func encodeSettings(snapshot ConfigData) (json.RawMessage, error) {
	rows := make(map[string]json.RawMessage, len(settingsDefinitions))
	for _, definition := range settingsDefinitions {
		raw, err := definition.Encode(snapshot)
		if err != nil {
			return nil, fmt.Errorf("encode %s: %w", definition.Key(), err)
		}
		rows[definition.Key()] = raw
	}
	return json.Marshal(rows)
}

func cloneConfig(c ConfigData) ConfigData {
	if c.FontConfiguration.Whitelist != nil {
		c.FontConfiguration.Whitelist = append([]string{}, c.FontConfiguration.Whitelist...)
	}
	return c
}

func NewStore(db dal.DB, cfg *koanf.Koanf) Store {
	defaults := ConfigData{FontConfiguration: FontConfiguration{MaxPerChapter: 10}, CaptchaSettings: CaptchaSettings{Type: "none"}}
	if cfg != nil {
		// The fields contain only scalar values and string slices; koanf decoding
		// overlays configured values on the built-in defaults.
		if err := cfg.Unmarshal("siteConfig.default", &defaults); err != nil {
			panic(fmt.Errorf("decode default site configuration: %w", err))
		}
		if cfg.Exists("chapter-fonts.max-per-chapter") {
			defaults.FontConfiguration.MaxPerChapter = cfg.Int("chapter-fonts.max-per-chapter")
		}
		if cfg.Exists("chapter-fonts.whitelist") {
			defaults.FontConfiguration.Whitelist = cfg.Strings("chapter-fonts.whitelist")
		}
	}
	return &storeImpl{db: db, defaults: cloneConfig(defaults), active: cloneConfig(defaults)}
}
