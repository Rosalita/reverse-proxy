package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "https",
		Host:   "sha256.badssl.com",
	})
	proxy.Transport = customTransport()

	director := proxy.Director
	proxy.Director = func(req *http.Request) {
		director(req)
		req.Host = req.URL.Host
	}

	http.Handle("/", filterIP(&ProxyHandler{proxy: proxy}))
	log.Fatal(http.ListenAndServe(":8443", nil))
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
