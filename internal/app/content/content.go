package content

import (
	"io"
	"strings"
	"unicode"

	"github.com/MaratBR/openlibrary/internal/app/htmlsanitizer"
	"github.com/k3a/html2text"
	"github.com/tdewolff/minify/v2"
	minifyhtml "github.com/tdewolff/minify/v2/html"
	"golang.org/x/net/html"
)

var contentMinifier = func() *minify.M {
	m := minify.New()
	m.AddFunc("text/html", minifyhtml.Minify)
	return m
}()

// SanitizeHtml takes a string of HTML content and returns a sanitized version of it,
// free of potentially malicious tags and attributes.
func SanitizeHtml(html string) string {
	return htmlsanitizer.Sanitize(html)
}

func CountWordsHtml(html string) int32 {
	text := html2text.HTML2Text(html)
	words := countWordsPlainText(text)
	return words
}

func countWordsPlainText(content string) int32 {

	println(content)

	r := strings.NewReader(content)
	var (
		words        int32
		isWithinWord bool = false
	)

	for {
		if r, _, err := r.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		} else {
			if unicode.IsSpace(r) {
				if isWithinWord {
					isWithinWord = false
					words += 1
				}
			} else {
				if !isWithinWord {
					isWithinWord = true
				}
			}
		}
	}

	return words
}

func FixHTML(htmlSnippet string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlSnippet))
	if err != nil {
		return "", err
	}

	// Find the body node
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

	// Render only the body content
	var fixedHTML strings.Builder
	if body != nil {
		for c := body.FirstChild; c != nil; c = c.NextSibling {
			err = html.Render(&fixedHTML, c)
			if err != nil {
				return "", err
			}
		}
	}

	return fixedHTML.String(), nil
}
