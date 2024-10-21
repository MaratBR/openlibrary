package uiserver

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/koding/websocketproxy"
)

type proxy struct {
	target     *url.URL
	httpClient *http.Client
	wsProxy    http.Handler
}

// ServeHTTP implements http.Handler.
func (p *proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Upgrade") == "websocket" {
		p.wsProxy.ServeHTTP(w, req)
		return
	}
	p.proxy(w, req)
}

func (p *proxy) proxy(w http.ResponseWriter, req *http.Request) {
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
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

}

func newProxy(targetHost string) http.Handler {
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
	return &proxy{
		target: u,
		httpClient: &http.Client{
			Timeout: time.Second * 60,
		},
		wsProxy: websocketproxy.NewProxy(&wsUrl),
	}
}

func NewDevServerProxy() http.Handler {
	return newProxy("http://localhost:5173")
}
