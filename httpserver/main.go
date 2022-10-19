package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("you must provide a port, see readme.md for usage.")
	}

	p := flag.String("p", "", "the port http server will run on")
	flag.Parse()
	port := fmt.Sprintf(":%s", *p)

	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(port, nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	log.Println("Request received")
	log.Printf("Request method: %s\n", r.Method)
	log.Printf("Request url: %s\n", r.URL.String())
	log.Printf("Request headers: %v\n", r.Header)

	serverDate := now.Format("02-01-2006")
	serverTime := now.Format("15:04:05")

	fmt.Fprintf(w, "Request received at: %s, %s", serverDate, serverTime)
}
