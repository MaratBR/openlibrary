package content

import (
	"encoding/json"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
	"strconv"
)

type WidgetRegistry struct {
}

const WidgetElement = "ol-widget"

func NewWidgetRegistry() *WidgetRegistry {
	return &WidgetRegistry{}
}

func (w *WidgetRegistry) Expand(node *html.Node) {
	attr := getAttributes(node)
	name, ok := attr["name"]
	if !ok {
		w.writeError(node, "no 'name' attribute")
		return
	}
	if name == "book-preview" {
		w.expandBookPreview(node, attr["data"])
		return
	}

	clearElement(node)
	node.Data = "div"
	node.Attr = []html.Attribute{
		{Key: "class", Val: "widget"},
		{Key: "data-widget", Val: name},
	}
	node.AppendChild(&html.Node{
		Type: html.TextNode,
		Data: "widget " + name,
	})
}

func (w *WidgetRegistry) expandBookPreview(node *html.Node, data string) {
	var reference struct {
		BookID string `json:"bookId"`
	}
	if err := json.Unmarshal([]byte(data), &reference); err != nil {
		w.writeError(node, "invalid book reference")
		return
	}
	id, err := strconv.ParseInt(reference.BookID, 10, 64)
	if err != nil || id <= 0 {
		w.writeError(node, "invalid book reference")
		return
	}
	bookID := strconv.FormatInt(id, 10)
	clearElement(node)
	node.Data = "div"
	node.Attr = []html.Attribute{
		{Key: "class", Val: "BookPreviewWidget"},
		{Key: "data-widget", Val: "book-preview"},
		{Key: "data-book-id", Val: bookID},
	}
	link := &html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{{Key: "href", Val: "/book/" + bookID}},
	}
	link.AppendChild(&html.Node{Type: html.TextNode, Data: "/book/" + bookID})
	node.AppendChild(link)
}

func (w *WidgetRegistry) ConfigureSanitizer(policy *bluemonday.Policy) {
	policy.AllowElements(WidgetElement)
	policy.AllowAttrs("name", "data").OnElements(WidgetElement)
}

func (w *WidgetRegistry) writeError(node *html.Node, text string) {

	clearElement(node)
	node.Data = "span"
	node.Attr = []html.Attribute{
		{Key: "data-widget-error", Val: "true"},
	}
	node.AppendChild(&html.Node{
		Type: html.TextNode,
		Data: "widget error: " + text,
	})
}

func getAttributes(node *html.Node) map[string]string {
	m := make(map[string]string)

	for _, attr := range node.Attr {
		if attr.Namespace != "" {
			continue
		}

		m[attr.Key] = attr.Val
	}

	return m
}

func getAttr(node *html.Node, name string) (val string) {
	for _, attr := range node.Attr {
		if attr.Key == name && attr.Namespace == "" {
			val = attr.Val
			return
		}
	}

	return
}

func clearElement(node *html.Node) {
	for node.FirstChild != nil {
		node.RemoveChild(node.FirstChild)
	}

	node.Attr = nil
}
