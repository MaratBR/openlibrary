package ao3import

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 120,
		},
	}
}

func makeGetReq(url string) (*http.Request, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Host", "archiveofourown.org")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36")
	return req, nil
}

func (c *Client) GetBook(id string) (*Ao3Book, error) {
	slog.Info("fetching ao3 book", "id", id)

	var (
		book        *Ao3Book
		downloadUrl string
	)

	{
		url := fmt.Sprintf("https://archiveofourown.org/works/%s", id)
		var resp *http.Response
		req, err := makeGetReq(url)
		if err != nil {
			return nil, err
		}
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode == 404 {
			return nil, fmt.Errorf("book %s not found", id)
		}
		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("failed to fetch book %s: %s", url, resp.Status)
		}
		downloadUrl, err = getDownloadUrl(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to get download url for book %s: %s", url, err.Error())
		}
		slog.Info("book download url", "id", id, "url", downloadUrl)
	}

	{
		var resp *http.Response
		req, err := makeGetReq(downloadUrl)
		if err != nil {
			return nil, err
		}
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		resp.Body.Close()

		writeBookTmp(id, respBytes)

		book, err = parseBook(bytes.NewReader(respBytes))
		if err != nil {
			return nil, err
		}
	}

	return book, nil
}

func writeBookTmp(id string, b []byte) {
	err := os.MkdirAll("/tmp/ao3import", os.ModePerm)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(fmt.Sprintf("/tmp/ao3import/%s.html", id))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.Write(b)
	if err != nil {
		panic(err)
	}
}

func (c *Client) ScrapeBookIDs(urlString string, maxPages int, out chan<- string) error {

	var (
		remainingPages = maxPages
		u              *url.URL
	)

	u, err := url.Parse(urlString)
	if err != nil {
		return err
	}

	for remainingPages > 0 {
		var resp *http.Response
		req, err := makeGetReq(u.String())
		if err != nil {
			return err
		}
		slog.Info("fetching", "url", u.String())
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 404 {
			slog.Error("encountered 404", "url", u.String())
			break
		}
		if resp.StatusCode >= 400 {
			slog.Error("failed to fetch book", "url", u.String(), "status", resp.Status)
			break
		}

		ids, href, err := getBookIds(resp.Body)
		if err != nil {
			slog.Error("failed to process page")
			break
		}

		for _, id := range ids {
			out <- id
		}

		if href == "" {
			break
		}

		u, err = url.Parse(fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, href))
		if err != nil {
			slog.Error("failed to parse href of next page", "href", href)
			break
		}

		remainingPages -= 1
	}

	return nil
}
