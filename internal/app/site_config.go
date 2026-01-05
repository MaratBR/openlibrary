package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/knadh/koanf/v2"
)

type SiteConfigEntry json.RawMessage

type SiteConfigData struct {
	CaptchaSettings      CaptchaSettings
	PasswordRequirements PasswordRequirements
	ContentRestrictions  ContentRestrictions
}

func (d SiteConfigData) Equal(other SiteConfigData) bool {
	return reflect.DeepEqual(d, other)
}

type SiteConfig struct {
	changed     bool
	db          DB
	lastFetched time.Time
	staticCfg   *koanf.Koanf

	activeConfigMX sync.RWMutex
	saveMX         sync.Mutex

	activeConfig  SiteConfigData
	defaultConfig SiteConfigData
	dbConfig      *SiteConfigData
}

type CaptchaSettings struct {
	GoogleRecaptchaKey string
	Type               string
}

var (
	PasswordError = apperror.AppErrors.NewType("password")
)

type PasswordRequirements struct {
	Digits         bool
	Symbols        string
	SymbolsEnabled bool
	DifferentCases bool
	MinLength      int
}

var (
	regexDigits    = regexp.MustCompile("\\d")
	regexUppercase = regexp.MustCompile("[A-Z]")
	regexLowercase = regexp.MustCompile("[a-z]")
)

func ValidatePassword(pwd string, r PasswordRequirements) error {
	if r.Digits && !regexDigits.Match([]byte(pwd)) {
		return PasswordError.New("password must have digits")
	}
	if r.MinLength > 0 && len(pwd) < r.MinLength {
		return PasswordError.New(fmt.Sprintf("password must be at least %d characters long", r.MinLength))
	}
	if r.DifferentCases && !(regexLowercase.Match([]byte(pwd)) && regexUppercase.Match([]byte(pwd))) {
		return PasswordError.New("password must contain at least one uppercase and one lowercase letter")
	}
	if r.SymbolsEnabled && len(r.Symbols) > 0 {
		var containsSymbol bool
		for _, rune := range r.Symbols {
			if strings.ContainsRune(pwd, rune) {
				containsSymbol = true
				break
			}
		}

		if !containsSymbol {
			return PasswordError.New("password")
		}
	}
	return nil
}

type ContentRestrictions struct {
	// if true - whole website is considered "adult"
	AdultWebsite bool
}

func NewSiteConfig(db DB, staticCfg *koanf.Koanf) *SiteConfig {
	cfg := &SiteConfig{
		db:        db,
		staticCfg: staticCfg,
	}
	err := cfg.loadDefault()
	if err != nil {
		slog.Error("failed to load default site config", "err", err)
	}
	return cfg
}

func (c *SiteConfig) loadDefault() error {
	value := c.staticCfg.Get("siteConfig.default")
	jsonString, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonString, &c.defaultConfig)
	if err != nil {
		return err
	}
	c.activeConfig = c.defaultConfig
	return nil
}

func (s *SiteConfig) Load(ctx context.Context) error {
	queries := store.New(s.db)
	jsonCfg, err := queries.SiteConfig_Get(ctx)
	if err != nil {
		if err == store.ErrNoRows {
			return nil
		} else {
			return apperror.WrapUnexpectedDBError(err)
		}
	}

	cfg := new(SiteConfigData)
	err = json.Unmarshal(jsonCfg, cfg)

	s.activeConfigMX.Lock()
	defer s.activeConfigMX.Unlock()
	s.dbConfig = cfg
	s.activeConfig = *cfg

	return nil
}

func (s *SiteConfig) Get() *SiteConfigData {
	s.activeConfigMX.RLock()
	defer s.activeConfigMX.RUnlock()
	return &s.activeConfig
}

func (s *SiteConfig) Save(ctx context.Context, force bool) error {
	currentConfig := *s.Get()

	s.saveMX.Lock()
	defer s.saveMX.Unlock()

	if force {
		err := s.save(ctx, currentConfig)
		return err
	} else {
		if s.dbConfig != nil && currentConfig.Equal(*s.dbConfig) {
			// no changes
			return nil
		}

		err := s.save(ctx, currentConfig)
		return err
	}
}

func (s *SiteConfig) save(ctx context.Context, cfg SiteConfigData) error {
	newJson, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	queries := store.New(s.db)
	err = queries.SiteConfig_Set(ctx, newJson)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	return nil
}
