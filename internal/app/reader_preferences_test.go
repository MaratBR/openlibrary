package app

import "testing"

func TestReaderContentWidthValidation(t *testing.T) {
	for _, width := range []int16{0, 48, 72, 100} {
		p := DefaultReaderPreferences()
		p.ContentWidth = width
		if err := p.Validate(); err != nil {
			t.Errorf("width %d: %v", width, err)
		}
	}
	for _, width := range []int16{-1, 46, 49, 101, 102} {
		p := DefaultReaderPreferences()
		p.ContentWidth = width
		if err := p.Validate(); err == nil {
			t.Errorf("accepted invalid width %d", width)
		}
	}
}

func TestReaderFontCatalogValidation(t *testing.T) {
	fonts := NewReaderFontService().List()
	if len(fonts) == 0 {
		t.Fatal("empty font catalog")
	}
	seen := map[string]bool{}
	for _, font := range fonts {
		if seen[font.ID] || font.LabelKey == "" || font.Family == "" {
			t.Fatalf("invalid font: %+v", font)
		}
		seen[font.ID] = true
		p := DefaultReaderPreferences()
		p.FontFamily = font.ID
		if err := p.Validate(); err != nil {
			t.Fatal(err)
		}
	}
	p := DefaultReaderPreferences()
	p.FontFamily = "unknown"
	if p.Validate() == nil {
		t.Fatal("accepted unknown font")
	}
}

func TestReaderPreferencesUseSuppliedFontCatalog(t *testing.T) {
	p := DefaultReaderPreferences()
	p.FontFamily = "custom"
	if err := p.ValidateWithFonts([]ReaderFont{{ID: "custom", LabelKey: "reader.custom", Family: "Custom, serif"}}); err != nil {
		t.Fatal(err)
	}
	if p.ValidateWithFonts(nil) == nil {
		t.Fatal("accepted font outside supplied catalog")
	}
}

func TestReaderThemeValidation(t *testing.T) {
	for _, theme := range []string{ReaderThemeSystem, ReaderThemeLight, ReaderThemeDark, ReaderThemeOLED} {
		p := DefaultReaderPreferences()
		p.Theme = theme
		if err := p.Validate(); err != nil {
			t.Errorf("theme %s: %v", theme, err)
		}
	}
	p := DefaultReaderPreferences()
	p.Theme = "unknown"
	if p.Validate() == nil {
		t.Fatal("accepted unknown theme")
	}
}
