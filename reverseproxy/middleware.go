package main

import (
	"log"
	"net"
	"net/http"
)

// filterIP is a middleware handler that will filter out requests from bad IPs.
func filterIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, err := getRequestIP(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		badIPs := blocklist()

		if _, ok := badIPs[ip]; ok {
			w.WriteHeader(http.StatusForbidden)
			log.Println("blocked request from bad ip")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getRequestIP extracts the host from a http request.
func getRequestIP(req *http.Request) (string, error) {
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", err
	}
	return host, nil
}
