package app

import (
	"io"
	"strings"

	"github.com/MaratBR/openlibrary/internal/app/htmlsanitizer"
	"github.com/k3a/html2text"
	"golang.org/x/net/html"
)

// SanitizeHtml takes a string of HTML content and returns a sanitized version of it,
// free of potentially malicious tags and attributes.
func SanitizeHtml(html string) string {
	return htmlsanitizer.Sanitize(html)
}

type ProcessedContentData struct {
	Sanitized string
	Words     int32
}

// ProcessContent takes a string of HTML content and returns a ProcessedContentData
// containing both a sanitized version of the content (i.e. HTML tags removed and
// unsafe content stripped), and a count of the number of words in the content.
func ProcessContent(content string) (ProcessedContentData, error) {
	fixedHtml := SanitizeHtml(content)
	fixedHtml, err := FixHTML(fixedHtml)
	if err != nil {
		return ProcessedContentData{}, err
	}
	words := CountWordsHtml(fixedHtml)
	return ProcessedContentData{Sanitized: fixedHtml, Words: words}, nil
}

func CountWordsHtml(html string) int32 {
	text := html2text.HTML2Text(html)
	words := countWordsPlainText(text)
	return words
}

func countWordsPlainText(content string) int32 {

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
			if r == ' ' || r == '\n' || r == '\t' {
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
