package main

import (
	"crypto/tls"
	"log"
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
	t.DialTLS = dialTLS
	t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return t
}

func dialTLS(network, addr string) (net.Conn, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{ServerName: host}

	tlsConn := tls.Client(conn, cfg)
	if err := tlsConn.Handshake(); err != nil {
		conn.Close()
		return nil, err
	}

	cs := tlsConn.ConnectionState()
	cert := cs.PeerCertificates[0]

	cert.VerifyHostname(host)
	log.Println(cert.Subject)

	return tlsConn, nil
}
