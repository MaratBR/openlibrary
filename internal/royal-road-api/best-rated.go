package royalroadapi

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type BestRatedBook struct {
	ID   int
	Name string
}

type BestRatedBooksResponse struct {
	Books    []BestRatedBook
	LastPage int
}

func (c *Client) GetBestRatedBooks(page int) (BestRatedBooksResponse, error) {
	if page < 1 {
		page = 1
	}
	url := fmt.Sprintf("https://www.royalroad.com/fictions/best-rated?page=%d", page)
	req, err := c.createGetRequest(url)
	if err != nil {
		return BestRatedBooksResponse{}, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return BestRatedBooksResponse{}, err
	}
	defer resp.Body.Close()
	return parseBestRatedBooks(resp.Body)
}

func parseBestRatedBooks(readCloser io.ReadCloser) (BestRatedBooksResponse, error) {
	doc, err := goquery.NewDocumentFromReader(readCloser)
	if err != nil {
		return BestRatedBooksResponse{}, err
	}

	li := doc.Find(".pagination").Find("li a").Last()
	pageStr, _ := li.Attr("data-page")
	lastPage, err := strconv.Atoi(pageStr)
	if err != nil {
		return BestRatedBooksResponse{}, err
	}

	books := []BestRatedBook{}
	doc.Find(".fiction-list-item").Each(func(i int, s *goquery.Selection) {
		link := s.Find(".fiction-title a").First()
		href, _ := link.Attr("href")
		href = strings.TrimPrefix(href, "/fiction/")
		idx := strings.IndexRune(href, '/')
		if idx == -1 {
			return
		}
		id, err := strconv.Atoi(href[:idx])
		if err != nil {
			return
		}
		books = append(books, BestRatedBook{
			ID:   id,
			Name: link.Text(),
		})
	})

	return BestRatedBooksResponse{
		Books:    books,
		LastPage: lastPage,
	}, nil
}
