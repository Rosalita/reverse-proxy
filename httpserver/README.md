# http server

A simple command line utility to create simple http servers.

## Usage
To start a http server on port 1234
`go run main.go -p 1234`

When the http server receives a request, it logs request details.
Each request will be responded to with the date and time the request was received.