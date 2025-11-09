package royalroadapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type BookPage struct {
	ID            int
	Name          string
	CoverURL      string
	Description   string
	Author        BookPageAuthor
	Chapters      []BookPageChapter
	Tags          []string
	ReviewsCount  int
	TotalViews    int
	Favortites    int
	Followers     int
	Ratings       int
	AverageRating float64
}

type BookPageAuthor struct {
	Name string
	ID   int
}

type BookPageChapter struct {
	ID         int
	Title      string
	ReleasedAt time.Time
}

func (c *Client) GetBookPage(bookID int) (*BookPage, error) {
	url := fmt.Sprintf("https://www.royalroad.com/fiction/%d", bookID)
	req, err := c.createGetRequest(url)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	book, err := parseBookPage(resp.Body)
	if err != nil {
		return nil, err
	}
	book.ID = bookID
	return book, nil
}

type bookSchema struct {
	Genre       []string `json:"genre"`
	Image       string   `json:"image"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
}

func parseBookPage(r io.Reader) (*BookPage, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	jsonMetadata := doc.Find("script[type=\"application/ld+json\"]").First().Text()
	if jsonMetadata == "" {
		return nil, fmt.Errorf("could not find book metadata: script[type=application/ld+json] selector is empty")
	}

	var schema bookSchema
	err = json.Unmarshal([]byte(jsonMetadata), &schema)
	if err != nil {
		return nil, fmt.Errorf("failed to parse book metadata: %s", err.Error())
	}

	page := &BookPage{
		Name:        schema.Name,
		CoverURL:    schema.Image,
		Description: schema.Description,
		Chapters:    []BookPageChapter{},
		Tags:        []string{},
	}

	doc.Find(".tags").Find(".fiction-tag").Each(func(i int, s *goquery.Selection) {
		page.Tags = append(page.Tags, s.Text())
	})

	dataReviewsPage, ok := doc.Find(".reviews-container ul.pagination").First().Find("li").Last().Find("a").First().Attr("data-page")
	if ok {
		reviewsPages, err := strconv.ParseInt(dataReviewsPage, 10, 64)
		if err == nil {
			page.ReviewsCount = int(10 * reviewsPages) // approx number of reviews
		}
	}

	statsContainer := doc.Find(".stats-content").First()
	overallScore := statsContainer.Find(".list-item").First().Find("span").First()
	label, _ := overallScore.Attr("aria-label")
	ratingStr := strings.TrimSuffix(label, " stars")
	rating, err := strconv.ParseFloat(ratingStr, 64)
	if err == nil {
		page.AverageRating = rating
	}

	secondColumnUL := statsContainer.Find(".col-sm-6").Last().Find("ul").First()

	page.TotalViews = parseIntStat(secondColumnUL.Find("li").Eq(1))
	page.Followers = parseIntStat(secondColumnUL.Find("li").Eq(5))
	page.Favortites = parseIntStat(secondColumnUL.Find("li").Eq(7))
	page.Ratings = parseIntStat(secondColumnUL.Find("li").Eq(9))

	doc.Find(".chapter-row").Each(func(i int, s *goquery.Selection) {
		link := s.Find("td").First().Find("a").First()
		title := link.Text()
		if title == "" {
			slog.Error("failed parse chapter name from title", "index", i)
			return
		}
		prefix := fmt.Sprintf("%d. ", i+1)
		title = strings.TrimPrefix(title, prefix)
		title = strings.Trim(title, " \n\t")

		url, _ := link.Attr("href")
		if url == "" {
			slog.Error("failed to et chapter url from data parameter")
			return
		}

		chapterID, err := parseChapterIDFromRelativeURL(url)
		if err != nil {
			slog.Error("failed to parse chapter id", "err", err, "index", i)
			return
		}

		timeTag := s.Find("time").First()
		timeStr, _ := timeTag.Attr("datetime")
		releasedAt, err := time.Parse(time.RFC3339Nano, timeStr)
		if err != nil {
			slog.Error("failed to parse chapter release date, ignoring this error", "err", err, "index", i)
		}

		page.Chapters = append(page.Chapters, BookPageChapter{
			ID:         chapterID,
			Title:      title,
			ReleasedAt: releasedAt,
		})
	})

	return page, nil

}

func parseIntStat(s *goquery.Selection) int {
	text := s.Text()
	text = strings.Replace(text, ",", "", -1)
	integer, _ := strconv.ParseInt(text, 10, 64)
	return int(integer)
}

func parseChapterIDFromRelativeURL(s string) (int, error) {
	const prefix string = "/chapter/"
	idx := strings.Index(s, prefix)
	if idx == -1 {
		return 0, fmt.Errorf("could not find chapter id: %s", s)
	}
	s = s[idx+len(prefix):]
	idx = strings.IndexRune(s, '/')
	if idx == -1 {
		return 0, fmt.Errorf("could not find chapter id: %s", s)
	}
	chapterID, err := strconv.Atoi(s[:idx])
	if err != nil {
		return 0, err
	}
	return chapterID, nil
}

type BookWithChapters struct {
	Book     *BookPage
	Chapters []*ChapterPage
}

func (c *Client) GetBookWithChapters(bookID int) (*BookWithChapters, error) {
	page, err := c.GetBookPage(bookID)
	if err != nil {
		return nil, err
	}

	book := &BookWithChapters{
		Book:     page,
		Chapters: make([]*ChapterPage, len(page.Chapters)),
	}

	// for i := range page.Chapters {
	// 	slog.Debug("downloading chapter", "chapterID", page.Chapters[i].ID)
	// 	book.Chapters[i], err = c.GetChapterPage(bookID, page.Chapters[i].ID)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to fetching chapter %d (id=%d): %s", i+1, page.Chapters[i].ID, err.Error())
	// 	}
	// }

	return book, nil
}
