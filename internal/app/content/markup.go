package content

import (
	"errors"
	"maps"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
)

type ExpandTag interface {
	Expand(node *html.Node)
	ConfigureSanitizer(policy *bluemonday.Policy)
}

// ExpandTagFunc adapts a function to an ExpandTag. It is useful for small
// tag expanders and tests that do not need to change the sanitizer policy.
type ExpandTagFunc func(node *html.Node)

func (f ExpandTagFunc) Expand(node *html.Node) {
	f(node)
}

func (ExpandTagFunc) ConfigureSanitizer(*bluemonday.Policy) {}

type MarkupEngine struct {
	tagExpanders    map[string]ExpandTag
	sanitizerPolicy *bluemonday.Policy
	urlFilter       func(url string) bool
}

type MarkupEngineOptions struct {
	TagExapanders map[string]ExpandTag

	// URLFilter decides whether a URL used by markup such as links and images is
	// permitted. A nil filter leaves URL validation to the sanitizer policy.
	URLFilter func(url string) bool

	// TextColorFilter decides whether a CSS text color value is permitted. A nil
	// filter leaves text color validation to the sanitizer policy.
	TextColorFilter func(color string) bool

	// AllowedFontFamilies is the fixed set of permitted font-family CSS values.
	AllowedFontFamilies []string

	// FontFamilyFilter permits dynamic font-family values.
	FontFamilyFilter func(fontFamily string) bool

	// AllowedFontSizes is the fixed set of permitted font-size CSS values.
	AllowedFontSizes []string
}

func NewEngine(options MarkupEngineOptions) *MarkupEngine {

	eng := &MarkupEngine{
		tagExpanders: maps.Clone(options.TagExapanders),
		urlFilter:    options.URLFilter,
	}
	eng.createSanitizerPolicy(options)
	return eng
}

func (e *MarkupEngine) createSanitizerPolicy(options MarkupEngineOptions) {
	policy := bluemonday.UGCPolicy()

	policy.AllowStyles("text-align").MatchingEnum("left", "right", "center", "justify").Globally()
	policy.AllowAttrs("target").Matching(regexp.MustCompile(`^_blank$`)).OnElements("a")
	policy.AllowAttrs("data-lang").Matching(regexp.MustCompile(`^[a-zA-Z0-9_+-]+$`)).OnElements("pre")

	if options.TextColorFilter != nil {
		policy.AllowStyles("color").MatchingHandler(options.TextColorFilter).Globally()
	}
	if options.FontFamilyFilter != nil {
		policy.AllowStyles("font-family").MatchingHandler(options.FontFamilyFilter).Globally()
	} else if len(options.AllowedFontFamilies) > 0 {
		policy.AllowStyles("font-family").MatchingEnum(options.AllowedFontFamilies...).Globally()
	}
	if len(options.AllowedFontSizes) > 0 {
		policy.AllowStyles("font-size").MatchingEnum(options.AllowedFontSizes...).Globally()
	}

	for _, expander := range e.tagExpanders {
		expander.ConfigureSanitizer(policy)
	}

	e.sanitizerPolicy = policy
}

type ProcessedContentData struct {
	Sanitized string
	Words     int32
	Fonts     []string
}

func (e *MarkupEngine) Clean(value string) (processed ProcessedContentData, err error) {
	processed.Sanitized, err = e.process(value, false, true)
	if err != nil {
		return ProcessedContentData{}, err
	}
	processed.Words = CountWordsHtml(processed.Sanitized)
	processed.Fonts, err = collectFontFamilies(processed.Sanitized)
	if err != nil {
		return ProcessedContentData{}, err
	}

	return processed, nil
}

func collectFontFamilies(content string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	fonts := make([]string, 0)
	seen := make(map[string]struct{})
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode {
			for _, attr := range node.Attr {
				if attr.Namespace != "" || attr.Key != "style" {
					continue
				}
				for _, declaration := range strings.Split(attr.Val, ";") {
					property, value, ok := strings.Cut(declaration, ":")
					if !ok || !strings.EqualFold(strings.TrimSpace(property), "font-family") {
						continue
					}
					for _, font := range splitFontFamilyList(value) {
						if _, ok := seen[font]; ok {
							continue
						}
						seen[font] = struct{}{}
						fonts = append(fonts, font)
					}
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(doc)
	return fonts, nil
}

func splitFontFamilyList(value string) []string {
	var fonts []string
	start := 0
	var quote rune
	escaped := false
	for index, char := range value {
		if escaped {
			escaped = false
			continue
		}
		if char == '\\' && quote != 0 {
			escaped = true
			continue
		}
		if char == '\'' || char == '"' {
			if quote == 0 {
				quote = char
			} else if quote == char {
				quote = 0
			}
			continue
		}
		if char == ',' && quote == 0 {
			fonts = appendFontFamily(fonts, value[start:index])
			start = index + 1
		}
	}
	return appendFontFamily(fonts, value[start:])
}

func appendFontFamily(fonts []string, value string) []string {
	font := strings.TrimSpace(value)
	if len(font) >= 2 && ((font[0] == '\'' && font[len(font)-1] == '\'') || (font[0] == '"' && font[len(font)-1] == '"')) {
		font = font[1 : len(font)-1]
	}
	if font != "" {
		fonts = append(fonts, font)
	}
	return fonts
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

	if e.urlFilter != nil {
		e.filterURLs(body)
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

func (e *MarkupEngine) filterURLs(node *html.Node) {
	if node.Type == html.ElementNode {
		for index := 0; index < len(node.Attr); {
			attr := node.Attr[index]
			if attr.Namespace == "" && (attr.Key == "href" || attr.Key == "src") && !e.urlFilter(attr.Val) {
				node.Attr = append(node.Attr[:index], node.Attr[index+1:]...)
				if attr.Key == "href" {
					e.removeAttribute(node, "rel")
				}
				continue
			}
			index++
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		e.filterURLs(child)
	}
}

func (e *MarkupEngine) removeAttribute(node *html.Node, name string) {
	for index, attr := range node.Attr {
		if attr.Namespace == "" && attr.Key == name {
			node.Attr = append(node.Attr[:index], node.Attr[index+1:]...)
			return
		}
	}
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
