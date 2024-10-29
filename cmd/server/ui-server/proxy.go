package uiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/koding/websocketproxy"
)

type devProxy struct {
	target     *url.URL
	httpClient *http.Client
	wsProxy    http.Handler
	options    DevServerOptions
}

// ServeHTTP implements http.Handler.
func (p *devProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Upgrade") == "websocket" {
		p.wsProxy.ServeHTTP(w, req)
		return
	}
	p.proxy(w, req)
}

func (p *devProxy) proxy(w http.ResponseWriter, req *http.Request) {
	// we need to buffer the body if we want to read it here and send it
	// in the request.
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// you can reassign the body if you need to parse it as multipart
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	// create a new url from the raw RequestURI sent by the client
	path := p.target.Path + req.RequestURI
	url := fmt.Sprintf("%s://%s%s", p.target.Scheme, p.target.Host, path)

	proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(body))

	// We may want to filter some headers, otherwise we could just use a shallow copy
	// proxyReq.Header = req.Header
	proxyReq.Header = make(http.Header)
	for h, val := range req.Header {
		proxyReq.Header[h] = val
	}

	resp, err := p.httpClient.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for h, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(h, value)
		}
	}

	if resp.Header.Get("Content-Type") == "text/html" {
		var serverData map[string]any
		if p.options.GetServerPushedData != nil {
			serverData = p.options.GetServerPushedData(req)
		}
		writeWithServerData(resp.StatusCode, w, resp.Body, serverData)
	} else {
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func injectStringAtEndOfBody(response []byte, injectedString string) ([]byte, bool) {
	idx := bytes.LastIndex(response, []byte("</body>"))
	if idx == -1 {
		return response, false
	}

	injectedBytes := []byte(injectedString)
	newResponse := make([]byte, len(response)+len(injectedBytes))
	copy(newResponse, response[:idx])
	copy(newResponse[idx:], injectedBytes)
	copy(newResponse[idx+len(injectedBytes):], response[idx:])
	return newResponse, true
}

func writeWithServerData(statusCode int, w http.ResponseWriter, r io.Reader, serverData map[string]any) error {
	if serverData == nil {
		w.WriteHeader(statusCode)
		_, err := io.Copy(w, r)
		return err
	}

	jsonString, err := json.Marshal(serverData)
	if err != nil {
		slog.Error("failed to marshal server data", "err", err)
		w.WriteHeader(statusCode)
		_, err = io.Copy(w, r)
		return err
	}

	html, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	idx := bytes.LastIndex(html, []byte("</body>"))
	if idx == -1 {
		w.WriteHeader(statusCode)
		_, err = w.Write(html)
		if err != nil {
			return err
		}
	} else {
		contentLength, _ := strconv.ParseInt(w.Header().Get("Content-Length"), 10, 32)

		injectedString := fmt.Sprintf("<script data-ts=\"%d\">window.SERVER_DATA=%s</script>", time.Now().Unix(), jsonString)
		html, _ := injectStringAtEndOfBody(html, injectedString)
		contentLength = int64(len(html))
		w.Header().Add("x-is-proxy", "yes")
		w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))

		w.WriteHeader(statusCode)
		w.Write(html)
	}

	return nil
}

func newDevProxy(targetHost string, options DevServerOptions) http.Handler {
	u, err := url.Parse(targetHost)
	if err != nil {
		panic(err)
	}

	wsUrl := *u
	if wsUrl.Scheme == "https" {
		wsUrl.Scheme = "wss"
	} else {
		wsUrl.Scheme = "ws"
	}
	return &devProxy{
		target: u,
		httpClient: &http.Client{
			Timeout: time.Second * 60,
		},
		wsProxy: websocketproxy.NewProxy(&wsUrl),
		options: options,
	}
}

type DevServerOptions struct {
	GetServerPushedData func(*http.Request) map[string]any
}

func NewDevServerProxy(options DevServerOptions) http.Handler {
	return newDevProxy("http://localhost:5173", options)
}
