package content

import (
	"errors"
	"maps"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
)

type ExpandTag interface {
	Expand(node *html.Node)
	ConfigureSanitizer(policy *bluemonday.Policy)
}

type MarkupEngine struct {
	tagExpanders    map[string]ExpandTag
	sanitizerPolicy *bluemonday.Policy
}

type MarkupEngineOptions struct {
	TagExapanders map[string]ExpandTag
}

func NewEngine(options MarkupEngineOptions) *MarkupEngine {

	eng := &MarkupEngine{
		tagExpanders: maps.Clone(options.TagExapanders),
	}
	eng.createSanitizerPolicy()
	return eng
}

func (e *MarkupEngine) createSanitizerPolicy() {
	policy := bluemonday.UGCPolicy()

	for _, expander := range e.tagExpanders {
		expander.ConfigureSanitizer(policy)
	}

	e.sanitizerPolicy = policy
}

type ProcessedContentData struct {
	Sanitized string
	Words     int32
}

func (e *MarkupEngine) Clean(value string) (processed ProcessedContentData, err error) {
	processed.Sanitized, err = e.process(value, false, true)
	processed.Words = CountWordsHtml(processed.Sanitized)

	return processed, nil
}

func (e *MarkupEngine) Expand(content string) (string, error) {
	return e.process(content, true, true)
}

func (e *MarkupEngine) process(content string, expandTags bool, sanitize bool) (processed string, err error) {
	if sanitize {
		content = e.sanitizerPolicy.Sanitize(content)
	}

	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return "", err
	}

	body := getBody(doc)
	if body == nil {
		return "", errors.New("cannot find body element")
	}

	if expandTags {
		e.expandRecursively(body, 0)
	}

	// Render only the body content
	var fixedHTML strings.Builder
	for c := body.FirstChild; c != nil; c = c.NextSibling {
		err = html.Render(&fixedHTML, c)
		if err != nil {
			return "", err
		}
	}

	html := fixedHTML.String()

	return html, nil
}

func (e *MarkupEngine) expandRecursively(node *html.Node, depth int) {
	if node.Type == html.ElementNode {
		e.expandTag(node)
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		e.expandRecursively(c, depth+1)
	}
}

func (e *MarkupEngine) expandTag(node *html.Node) {
	tag := node.Data

	expander, ok := e.tagExpanders[tag]

	if !ok {
		return
	}

	expander.Expand(node)
}

func getBody(doc *html.Node) *html.Node {
	var body *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			body = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return body
}
