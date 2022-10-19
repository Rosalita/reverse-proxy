package main

import (
	"net"
	"net/http"
	"time"
)

// customTransport configures a transport based on default transport
// but with a few modifications, see README.md for configuration explanation.
func customTransport() *http.Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()

	t.DialContext = (&net.Dialer{
		Timeout:   2 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	return t
}
