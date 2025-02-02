package app_test

import (
	"testing"

	"github.com/MaratBR/openlibrary/internal/app"
)

var (
	sanitizationCases = map[string]string{
		"<b>bold</b>":           "<b>bold</b>",
		"<i>italic</i>":         "<i>italic</i>",
		"<em>italic</em>":       "<em>italic</em>",
		"<strong>bold</strong>": "<strong>bold</strong>",

		"<script>alert('lol xss go brr')</script>": "",

		"<span></span>": "<span></span>",

		"<strong>cannot fix unclosed tags": "<strong>cannot fix unclosed tags",
	}
)

func TestContentSanitization(t *testing.T) {
	for input, expectedOutput := range sanitizationCases {
		output := app.SanitizeHtml(input)
		if output != expectedOutput {
			t.Errorf("Expected %s, got %s", expectedOutput, output)
		}
	}
}

var (
	contentProcessingTestCases = map[string]string{
		"<b>bold</b>":           "<b>bold</b>",
		"<i>italic</i>":         "<i>italic</i>",
		"<em>italic</em>":       "<em>italic</em>",
		"<strong>bold</strong>": "<strong>bold</strong>",
		"<script>alert('lol xss go brr')</script>": "",
		"<span></span>": "<span></span>",

		"<strong>can fix unclosed tags": "<strong>can fix unclosed tags</strong>",
	}
)

func TestContentProcessing(t *testing.T) {
	for input, expectedOutput := range contentProcessingTestCases {
		output, err := app.ProcessContent(input)
		if err != nil {
			t.Errorf("Expected no error, got %s", err.Error())
		}
		if output.Sanitized != expectedOutput {
			t.Errorf("Expected %s, got %s", expectedOutput, output.Sanitized)
		}
	}
}
