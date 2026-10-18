package i18n

import (
	"github.com/eduardolat/goeasyi18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"testing"
)

func TestIncompleteCatalogFallback(t *testing.T) {
	dir := t.TempDir()
	files := map[language.Tag][]string{}
	for lang, content := range map[language.Tag]string{
		language.English: "[editor]\nfontTools = \"Font\"\nname = \"English\"\nwelcome = \"Hello {{.Name}}\"\n",
		language.Spanish: "[editor]\nname = \"Español\"\n",
	} {
		path := filepath.Join(dir, lang.String()+".toml")
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
		files[lang] = []string{path}
	}
	// The selected/default language must not replace English as the fallback.
	engine, def := newI18N(language.Spanish, files, zap.NewNop().Sugar())
	localizer := newLocalizer(engine, def, language.Spanish, false)
	for key, want := range map[string]string{"editor.name": "Español", "editor.fontTools": "Font", "missing.key": "missing.key"} {
		if got := localizer.T(key); got != want {
			t.Errorf("%s: got %q want %q", key, got, want)
		}
	}
	if got := localizer.T("editor.welcome", goeasyi18n.Options{Data: map[string]string{"Name": "Reader"}}); got != "Hello Reader" {
		t.Fatal(got)
	}
	if got := localizer.TT("editor")["editor.fontTools"]; got != "Font" {
		t.Fatalf("frontend extraction omitted fallback: %q", got)
	}
	if got := newLocalizer(engine, def, language.German, false).T("editor.fontTools"); got != "Font" {
		t.Fatal(got)
	}
	if got := newLocalizer(engine, def, language.Spanish, true).T("editor.name"); got != "editor.name" {
		t.Fatal(got)
	}
}
