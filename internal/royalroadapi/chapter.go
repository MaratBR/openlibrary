package royalroadapi

import (
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
)

type ChapterPage struct {
	Content string
}

func (c *Client) GetChapterPage(bookID, chapterID int) (*ChapterPage, error) {
	url := fmt.Sprintf("https://www.royalroad.com/fiction/%d/chapter/%d", bookID, chapterID)
	req, err := c.createGetRequest(url)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch chapter page: %s", resp.Status)
	}
	defer resp.Body.Close()
	return parseChapterPage(resp.Body)
}

func parseChapterPage(r io.Reader) (*ChapterPage, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	content, err := doc.Find(".chapter-content").First().Html()
	if err != nil {
		return nil, fmt.Errorf("could not find chapter content: %s", err.Error())
	}

	if content == "" {
		html, _ := doc.Html()
		return nil, fmt.Errorf("chapter content is empty, full page html: %s", html)
	}

	return &ChapterPage{
		Content: content,
	}, nil
}
