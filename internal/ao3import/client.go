package ao3import

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"sync"
	"time"

	httplim "github.com/MaratBR/openlibrary/internal/http-lim"
	"golang.org/x/time/rate"
)

type Client struct {
	httpClient *httplim.HttpWorker
	mx         sync.Mutex
}

func NewClient() *Client {
	limiter := rate.NewLimiter(rate.Every(time.Second), 8)
	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}
	u, err := url.Parse("archiveofourown.org")
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("u is nil")
	}
	httpClient.Jar, err = cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	httpClient.Jar.SetCookies(u, []*http.Cookie{
		{
			Value:    "true",
			Name:     "view_adult",
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			Expires:  time.Now().Add(time.Hour * 24 * 30 * 12),
		},
	})
	return &Client{
		httpClient: httplim.NewHttpWorker(10, httpClient, limiter),
	}
}

func (c *Client) Close() {
	slog.Debug("closing http worker...")
	c.httpClient.Close()
	slog.Debug("http worker closed...")
}

func (c *Client) Run() {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.httpClient.Run()
}

func makeGetReq(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("host", "archiveofourown.org")
	req.Header.Set("authority", "archiveofourown.org")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="89", "Chromium";v="89", ";Not A Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("dnt", "1")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("accept-language", "en-GB,en;q=0.9")

	return req, nil
}

func doRequestWithAttempts(
	w *httplim.HttpWorker,
	req *http.Request,
	attempts int,
) (*http.Response, error) {
	if attempts < 1 {
		attempts = 1
	}

	var (
		err  error
		resp *http.Response
	)

	for attempts > 0 {
		attempts -= 1

		resp, err = w.Do(req)
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				time.Sleep(time.Millisecond * 1000)
				continue
			}
			return nil, err
		}
		if resp != nil && resp.StatusCode == 429 {
			time.Sleep(time.Millisecond * 500)
			continue
		}

		return resp, err
	}

	if err == nil {
		err = context.Canceled
	}

	return nil, err

}

func (c *Client) DownloadBook(id string) (*http.Response, error) {
	slog.Debug("fetching ao3 book", "id", id)

	var (
		downloadUrl string
	)

	{
		url := fmt.Sprintf("https://archiveofourown.org/works/%s?view_adult=true", id)
		var resp *http.Response
		req, err := makeGetReq(url)
		if err != nil {
			return nil, err
		}
		resp, err = doRequestWithAttempts(c.httpClient, req, 0)
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
		return resp, nil
	}
}

func (c *Client) GetBook(id string) (*Ao3Book, error) {
	resp, err := c.DownloadBook(id)
	if err != nil {
		return nil, err
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	writeBookTmp(id, respBytes)

	book, err := ParseBook(bytes.NewReader(respBytes))
	if err != nil {
		return nil, err
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
		slog.Debug("fetching", "url", u.String())
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

		slog.Debug("submitted ids", "count", len(ids))

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
