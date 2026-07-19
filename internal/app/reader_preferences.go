package app

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
)

const (
	ReaderFontSerif    = "serif"
	ReaderFontSans     = "sans"
	ReaderFontDyslexic = "dyslexic"

	ReaderPageBackground = "background"
	ReaderPageSurface    = "surface"

	ReaderThemeSystem = "system"
	ReaderThemeLight  = "light"
	ReaderThemeDark   = "dark"
)

type ReaderPreferences struct {
	FontSize   int16  `json:"fontSize"`
	FontFamily string `json:"fontFamily"`
	PageColor  string `json:"pageColor"`
	Theme      string `json:"theme"`
}

func DefaultReaderPreferences() ReaderPreferences {
	return ReaderPreferences{
		FontSize:   18,
		FontFamily: ReaderFontSerif,
		PageColor:  ReaderPageBackground,
		Theme:      ReaderThemeSystem,
	}
}

func (p ReaderPreferences) Validate() error {
	validFontSize := false
	for _, size := range []int16{12, 14, 16, 18, 20, 22, 26, 30, 36, 42, 48} {
		if p.FontSize == size {
			validFontSize = true
			break
		}
	}
	if !validFontSize {
		return errors.New("invalid font size")
	}
	if p.FontFamily != ReaderFontSerif && p.FontFamily != ReaderFontSans && p.FontFamily != ReaderFontDyslexic {
		return errors.New("invalid font family")
	}
	if p.PageColor != ReaderPageBackground && p.PageColor != ReaderPageSurface {
		return errors.New("invalid page color")
	}
	if p.Theme != ReaderThemeSystem && p.Theme != ReaderThemeLight && p.Theme != ReaderThemeDark {
		return errors.New("invalid reader theme")
	}
	return nil
}

type ReaderPreferencesService interface {
	Get(ctx context.Context, userID uuid.UUID) (Nullable[ReaderPreferences], error)
	Save(ctx context.Context, userID uuid.UUID, preferences ReaderPreferences) error
}
