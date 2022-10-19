package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	// originServer is the httpserver that the reverse proxy
	// will be forwarding requests to and receiving responses from.
	originServer string = "http://localhost:1234"
)

func main() {
	url, err := url.Parse(originServer)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = customTransport()

	http.Handle("/", filterIP(&ProxyHandler{proxy: proxy}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type ProxyHandler struct {
	proxy *httputil.ReverseProxy
}

func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ph.proxy.ModifyResponse = modifyResponse
	ph.proxy.ServeHTTP(w, r)
}

func modifyResponse(resp *http.Response) error {
	resp.Header.Set("x-searchpilot-interview", "hello")
	return nil
}
