package app

import (
	"context"
	"errors"

	"uuid"
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
	ReaderThemeOLED   = "oled"
)

type ReaderPreferences struct {
	ContentWidth int16  `json:"contentWidth"`
	FontSize     int16  `json:"fontSize"`
	FontFamily   string `json:"fontFamily"`
	PageColor    string `json:"pageColor"`
	Theme        string `json:"theme"`
}

func DefaultReaderPreferences() ReaderPreferences {
	return ReaderPreferences{
		ContentWidth: 72,
		FontSize:     18,
		FontFamily:   ReaderFontSerif,
		PageColor:    ReaderPageBackground,
		Theme:        ReaderThemeSystem,
	}
}

func (p ReaderPreferences) Validate() error {
	return p.ValidateWithFonts(NewReaderFontService().List())
}

func (p ReaderPreferences) ValidateWithFonts(fonts []ReaderFont) error {
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
	if p.ContentWidth != 0 && (p.ContentWidth < 48 || p.ContentWidth > 100 || p.ContentWidth%2 != 0) {
		return errors.New("invalid content width")
	}
	validFont := false
	for _, font := range fonts {
		if font.ID == p.FontFamily {
			validFont = true
			break
		}
	}
	if !validFont {
		return errors.New("invalid font family")
	}
	if p.PageColor != ReaderPageBackground && p.PageColor != ReaderPageSurface {
		return errors.New("invalid page color")
	}
	if p.Theme != ReaderThemeSystem && p.Theme != ReaderThemeLight && p.Theme != ReaderThemeDark && p.Theme != ReaderThemeOLED {
		return errors.New("invalid reader theme")
	}
	return nil
}

type ReaderPreferencesService interface {
	Get(ctx context.Context, userID uuid.UUID) (Nullable[ReaderPreferences], error)
	Save(ctx context.Context, userID uuid.UUID, preferences ReaderPreferences) error
}
