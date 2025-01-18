package templates

import (
	"log/slog"
	"math"
	"regexp"
	"strings"
	"time"
)

func approximateLines(width, textSize int, html string) int {
	t := time.Now()
	// Remove HTML tags to get the raw text
	text := stripHTMLTags(html)

	// Estimate the average character width based on the text size (pixels)
	// This is an approximation and can vary depending on the font.
	avgCharWidth := float64(textSize) * 0.6 // Assuming 0.6 of text size per character

	// Calculate the total width of the text in pixels
	totalTextWidth := int(math.Ceil(float64(len(text)) * avgCharWidth))

	// Calculate the approximate number of lines
	numLines := int(math.Ceil(float64(totalTextWidth) / float64(width)))

	took := time.Now().Sub(t)
	if took.Milliseconds() > 3 {
		slog.Warn("approximateLines took too long", "took", took.Milliseconds())
	}

	return numLines
}

// stripHTMLTags removes supported HTML tags from the input string
func stripHTMLTags(input string) string {
	re := regexp.MustCompile(`(?i)<(\/?(?:p|b|strong|em|i|span|br)[^>]*?)>`)
	output := re.ReplaceAllString(input, "")

	// Handle <br> as a line break replacement
	output = strings.ReplaceAll(output, "<br>", "\n")
	output = strings.ReplaceAll(output, "<br/>", "\n")
	output = strings.ReplaceAll(output, "<br />", "\n")

	return output
}
