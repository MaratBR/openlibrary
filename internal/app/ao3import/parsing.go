package ao3import

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getDownloadUrl(r io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}

	var url string

	doc.Find(".download li a").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "HTML" {
			url, _ = s.Attr("href")
		}
	})

	url = strings.Trim(url, " \n\t")
	if url == "" {
		return "", errors.New("could not find download url")
	}

	return fmt.Sprintf("https://download.archiveofourown.org%s", url), nil
}

func parseBook(r io.Reader) (*Ao3Book, error) {

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	name := doc.Find("h1").First().Text()
	book := &Ao3Book{
		Name:    name,
		Summary: getBookSummary(doc),
	}

	doc.Find("h2.heading").Each(func(i int, s *goquery.Selection) {
		html, err := s.Parent().Next().Html()
		if err != nil {
			return
		}

		book.Chapters = append(book.Chapters, Ao3Chapter{
			Title:   s.Text(),
			Content: html,
		})
	})

	parseTags(doc, book)

	{
		bookUrl := doc.Find(".message a:nth-child(2)").First().Text()
		parts := strings.Split(bookUrl, "/")
		book.ID = parts[len(parts)-1]
	}

	book.AuthorName = doc.Find("a[rel=author]").Text()

	return book, nil
}

func getBookSummary(doc *goquery.Document) string {
	byline := doc.Find(".byline").First()
	summaryLabel := byline.Next()
	if summaryLabel.Text() != "Summary" {
		return ""
	}

	html, err := summaryLabel.Next().Html()
	if err != nil {
		return ""
	}
	return html
}

func parseTags(doc *goquery.Document, book *Ao3Book) {
	tags := map[string][]string{}

	doc.Find(".tags dt").Each(func(i int, s *goquery.Selection) {
		tagType := s.Text()
		tagType = tagType[:len(tagType)-1]

		if tagType == "Stats" {
			return
		}

		if tagType == "Language" {
			book.Language = s.Next().Text()
			return
		}

		if tagType == "Rating" {
			book.Rating = s.Next().Text()
			return
		}

		tags[tagType] = []string{}

		s.Next().Find("a").Each(func(i int, s *goquery.Selection) {
			tags[tagType] = append(tags[tagType], s.Text())
		})
	})

	book.Tags = tags
}
