package ao3import

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 30,
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
		defer resp.Body.Close()

		book, err = parseBook(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	return book, nil
}
