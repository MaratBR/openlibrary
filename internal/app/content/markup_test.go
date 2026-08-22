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
	if result.Words != 3 {
		t.Errorf("Process().Words = %d, want 3", result.Words)
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

	result, err := engine.Expand("<p>Content</p>")
	if err != nil {
		t.Fatalf("Expand() error = %v", err)
	}

	const wantSanitized = "<p>Content<em> expanded</em></p>"
	if result != wantSanitized {
		t.Errorf("Expand() = %q, want %q", result, wantSanitized)
	}
}

func TestNewEngineCopiesTagExpanders(t *testing.T) {
	t.Parallel()

	expanders := map[string]ExpandTag{
		"strong": ExpandTagFunc(func(node *html.Node) { node.Data = "em" }),
	}
	engine := NewEngine(MarkupEngineOptions{TagExapanders: expanders})
	expanders["strong"] = ExpandTagFunc(func(node *html.Node) { node.Data = "i" })

	result, err := engine.Expand("<strong>Content</strong>")
	if err != nil {
		t.Fatalf("Expand() error = %v", err)
	}

	const want = "<em>Content</em>"
	if result != want {
		t.Errorf("Expand() = %q, want %q", result, want)
	}
}
