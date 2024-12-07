package uiserver

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/koding/websocketproxy"
)

type proxy struct {
	target     *url.URL
	httpClient *http.Client
	wsProxy    http.Handler
	options    DevServerOptions
}

// ServeHTTP implements http.Handler.
func (p proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Upgrade") == "websocket" {
		p.wsProxy.ServeHTTP(w, req)
		return
	}
	p.proxy(w, req)
}

func (p proxy) proxy(w http.ResponseWriter, req *http.Request) {
	// we need to buffer the body if we want to read it here and send it
	// in the request.
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// you can reassign the body if you need to parse it as multipart
	req.Body = io.NopCloser(bytes.NewReader(body))

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

	if !(strings.HasPrefix(path, "/node_modules/") || strings.HasPrefix(path, "/src/")) {
		// completely kill all cache
		proxyReq.Header["Cache-Control"] = []string{"no-cache"}
		proxyReq.Header["Pragma"] = []string{"no-cache"}
		delete(proxyReq.Header, "If-None-Match")
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
		var serverData []byte
		if p.options.GetInjectedHTMLSegment != nil {
			serverData = p.options.GetInjectedHTMLSegment(req)
		}
		injectAtHeadTailOfResponse(resp.StatusCode, w, resp.Body, serverData)
	} else {
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func injectAtHeadTail(response, headTail []byte) ([]byte, bool) {
	idx := bytes.LastIndex(response, []byte("</head>"))
	if idx == -1 {
		return response, false
	}

	injectedBytes := []byte(headTail)
	newResponse := make([]byte, len(response)+len(injectedBytes))
	copy(newResponse, response[:idx])
	copy(newResponse[idx:], injectedBytes)
	copy(newResponse[idx+len(injectedBytes):], response[idx:])
	return newResponse, true
}

func injectAtHeadTailOfResponse(statusCode int, w http.ResponseWriter, r io.Reader, injectedBytes []byte) error {
	if len(injectedBytes) == 0 {
		w.WriteHeader(statusCode)
		_, err := io.Copy(w, r)
		return err
	}

	html, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	idx := bytes.LastIndex(html, []byte("</head>"))
	if idx == -1 {
		w.WriteHeader(statusCode)
		_, err = w.Write(html)
		if err != nil {
			return err
		}
	} else {
		contentLength, _ := strconv.ParseInt(w.Header().Get("Content-Length"), 10, 32)

		html, _ := injectAtHeadTail(html, injectedBytes)
		contentLength = int64(len(html))
		w.Header().Add("x-is-proxy", "yes")
		w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))

		w.WriteHeader(statusCode)
		w.Write(html)
	}

	return nil
}

func newProxy(targetHost string, options DevServerOptions) proxy {
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
	return proxy{
		target: u,
		httpClient: &http.Client{
			Timeout: time.Second * 60,
		},
		wsProxy: websocketproxy.NewProxy(&wsUrl),
		options: options,
	}
}

type DevServerOptions struct {
	GetInjectedHTMLSegment func(*http.Request) []byte
}
