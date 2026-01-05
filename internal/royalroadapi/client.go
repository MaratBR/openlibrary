package royalroadapi

import (
	"log/slog"
	"net/http"
	"time"

	httplim "github.com/MaratBR/openlibrary/internal/http-lim"
	"golang.org/x/time/rate"
)

type Client struct {
	httpClient  *httplim.HttpWorker
	requestHost string
}

func (c *Client) createGetRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("host", c.requestHost)
	req.Header.Set("authority", c.requestHost)
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

func NewClient() *Client {
	httpClient := &http.Client{
		Timeout: time.Second * 120,
	}
	limiter := rate.NewLimiter(rate.Every(time.Millisecond*1000), 7)

	return &Client{
		httpClient:  httplim.NewHttpWorker(4, httpClient, limiter),
		requestHost: "www.royalroad.com",
	}
}

func (c *Client) Close() {
	slog.Debug("closing http worker...")
	c.httpClient.Close()
	slog.Debug("http worker closed...")
}

func (c *Client) Run() {
	c.httpClient.Run()
}
