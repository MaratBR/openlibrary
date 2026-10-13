package content

import (
	"strings"
	"testing"
)

func TestBookPreviewWidgetPersistence(t *testing.T) {
	engine := NewDefaultEngine()
	cleaned, err := engine.Clean(`<ol-widget name="book-preview" data='{"bookId":"9223372036854775807"}'><script>alert(1)</script><p>forged preview</p></ol-widget>`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(cleaned.Sanitized, "book-preview") || !strings.Contains(cleaned.Sanitized, "9223372036854775807") {
		t.Fatalf("reference lost: %s", cleaned.Sanitized)
	}
	expanded, err := engine.Expand(cleaned.Sanitized)
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{`data-widget="book-preview"`, `data-book-id="9223372036854775807"`, `href="/book/9223372036854775807"`} {
		if !strings.Contains(expanded, required) {
			t.Errorf("missing %s: %s", required, expanded)
		}
	}
	if strings.Contains(expanded, "forged preview") || strings.Contains(expanded, "alert(") {
		t.Fatalf("saved widget content was trusted: %s", expanded)
	}
}

func TestBookPreviewWidgetInvalidReferences(t *testing.T) {
	for _, data := range []string{`{}`, `null`, `{"bookId":null}`, `{"bookId":42}`, `{"bookId":"0"}`, `{"bookId":"-1"}`, `{"bookId":"9223372036854775808"}`, `{"bookId":"https://example.com"}`, `not-json`} {
		t.Run(data, func(t *testing.T) {
			result, err := NewDefaultEngine().Expand(`<ol-widget name="book-preview" data='` + data + `'></ol-widget>`)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(result, `data-widget-error="true"`) || strings.Contains(result, `href=`) {
				t.Fatalf("invalid reference expanded: %s", result)
			}
		})
	}
}
