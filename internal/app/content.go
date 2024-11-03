package app

import (
	"io"
	"strings"

	"github.com/MaratBR/openlibrary/internal/app/htmlsanitizer"
	"github.com/k3a/html2text"
)

func SanitizeHtml(html string) string {
	return htmlsanitizer.Sanitize(html)
}

type ProcessedContentData struct {
	Sanitized string
	Words     int32
}

func ProcessContent(content string) ProcessedContentData {
	sanitized := SanitizeHtml(content)
	words := CountWordsHtml(sanitized)
	return ProcessedContentData{Sanitized: sanitized, Words: words}
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
