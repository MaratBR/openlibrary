package content

import (
	"testing"

	"golang.org/x/net/html"
)

func TestMarkupEngineProcess(t *testing.T) {
	t.Parallel()

	engine := NewEngine(MarkupEngineOptions{})
	result, err := engine.Clean(`
		<p class="intro">Hello <strong>world</strong>.</p>
		<script>alert("unsafe")</script>
		<strong>Unclosed
	`)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	const wantSanitized = "<p>Hello <strong>world</strong>.</p>\n\t\t\n\t\t<strong>Unclosed\n\t</strong>"
	if result.Sanitized != wantSanitized {
		t.Errorf("Process().Sanitized = %q, want %q", result.Sanitized, wantSanitized)
	}
	if result.Words != 4 {
		t.Errorf("Process().Words = %d, want 4", result.Words)
	}
}

func TestMarkupEngineProcessExpandsTagsRecursively(t *testing.T) {
	t.Parallel()

	engine := NewEngine(MarkupEngineOptions{
		TagExapanders: map[string]ExpandTag{
			"p": ExpandTagFunc(func(node *html.Node) {
				node.AppendChild(&html.Node{Type: html.ElementNode, Data: "strong"})
			}),
			"strong": ExpandTagFunc(func(node *html.Node) {
				node.Data = "em"
				node.AppendChild(&html.Node{Type: html.TextNode, Data: " expanded"})
			}),
		},
	})

	result, err := engine.Clean("<p>Content</p>")
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	const wantSanitized = "<p>Content<em> expanded</em></p>"
	if result.Sanitized != wantSanitized {
		t.Errorf("Process().Sanitized = %q, want %q", result.Sanitized, wantSanitized)
	}
	if result.Words != 1 {
		t.Errorf("Process().Words = %d, want 1", result.Words)
	}
}

func TestNewEngineCopiesTagExpanders(t *testing.T) {
	t.Parallel()

	expanders := map[string]ExpandTag{
		"strong": ExpandTagFunc(func(node *html.Node) { node.Data = "em" }),
	}
	engine := NewEngine(MarkupEngineOptions{TagExapanders: expanders})
	expanders["strong"] = ExpandTagFunc(func(node *html.Node) { node.Data = "i" })

	result, err := engine.Clean("<strong>Content</strong>")
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	const want = "<em>Content</em>"
	if result.Sanitized != want {
		t.Errorf("Process().Sanitized = %q, want %q", result.Sanitized, want)
	}
}
