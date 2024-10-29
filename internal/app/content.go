package app

import (
	"io"
	"strings"
)

func CountWords(content string) int32 {
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

func CleanUpContent(content string) string {
	// TODO implement
	return content
}
