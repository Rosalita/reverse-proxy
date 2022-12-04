# reverse-proxy

* A small go program which acts as a reverse proxy and adds a header to all responses. 
* The proxy has been provisioned with a Transport using jusifiable configuration choices.
* The proxy will block any requests from a list of "bad" client ip addresses.

# Design Decisions

I want to keep things as simple I can, so I am going to use Go's excellent standard library to help me. 

This project will need a http server that responds to requests. I've used the `net/http` package 
in the standard library to write a CLI tool that will create http servers for me to test my reverse proxy with.

## Modifying Response

Looking at the documentation for the standard library reverse proxy there's an optional function
called `ModifyResponse` that modifies the response from the origin server.
It's signature is `ModifyResponse func(*http.Response) error` so will I use this to add the header.

## HTTP Transport

The standard library reverse proxy implementation by default will use `http.DefaultTransport`. 
The default transport has a dialer which will wait upto 30 seconds to establish a connection.
This timeout feels like it's too long. A TCP connection is initiated with a three way handshake. 
This being client to server, server to client and then client back to server. This means it will take 
take 1.5 times the round trip time to establish a connection. Assuming the worst possible round trip time,
possibly a satellite connection of 750ms, that would mean handshake should be done in 1125ms. 
So I have reduced this timeout to 2 seconds as it felt excessive. 

The default transport dialer also has a setting for `KeepAlive`. This is the frequency at which keep alive 
packets are sent. When there is a response to a keep alive packet, the idle timer is reset. 
This means that if the `KeepAlive` interval is greater than the `IdleConnTimeout` then keep alive packets
will never be sent. The default transport sets `KeepAlive` to 30 seconds, meaning every 30 seconds a 
keep alive packet is sent. At this time I can't see a reason to change this setting.

The default transport has `MaxIdleConns` set to 100, this is the maximum number of idle (keep alive)
connections across all hosts. The default transport also has `IdleConnTimeout` set to 90 seconds. 
This means if keep alive packets aren't responded to for 90 seconds, the connection will timeout.
At this time I can't see a reason to increase the connection pool of `MaxIdleConns` beyond 100.
The default `IdleConnTimeout` also feels sensible, if three keep alive packets, sent at 30 second intervals
haven't been responded to, then the connection should timeout.

The default transport uses `DefaultMaxIdleConnsPerHost`, which is set to 2. It also does not specify
`MaxConnsPerHost` which means this defaults to `DefaultMaxIdleConnsPerHost`. This means a single host can only 
have 2 connections. As many websites do parallel loading using multiple TCP connections e.g. loading html, css,
images, fonts etc. This setting feels like it needs to be greater than 2. As the `MaxIdleConns` pool 
of total connections is set to 100. To avoid using `DefaultMaxIdleConnsPerHost`, I have set 
`MaxConnsPerHost` to 100 and `MaxIdleConnsPerHost` to 100. This means a single host can take advantage
of all available connections in the connection pool if needed.

## Blocking requests from bad IPs

Every request that is received will need to be checked to see if it is from a bad IP.
For this reason, the IP filter that blocks requests should be middleware that runs before the
reverse proxy handles the request. In a larger project we would probably want to move middlewares
to their own package. In a production environment, I would also want to write test(s) for
middleware however the task does not request any tests to be written so on this occassion they 
have been omitted to save time.

# Running the code in this repository
1. In a terminal window, start a new http server on port 1234 using the command
`cd httpserver/ && go run main.go -p  1234`

2. Open a second terminal window, and start the reverse proxy server using the command
`cd reverseproxy/ && go build && ./reverseproxy`

3. Once everything is up and running, open a third terminal window and make a request to the 
reverse proxy server using the command
`curl -k -v http://127.0.0.1:8080`

You should see the following result:
```
$ curl -k -v http://127.0.0.1:8080
*   Trying 127.0.0.1:8080...
* Connected to 127.0.0.1 (127.0.0.1) port 8080 (#0)
> GET / HTTP/1.1
> Host: 127.0.0.1:8080
> User-Agent: curl/7.71.1
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Length: 41
< Content-Type: text/plain; charset=utf-8
< Date: Tue, 18 Oct 2022 19:22:15 GMT
< X-Searchpilot-Interview: hello
<
Request received at: 18-10-2022, 20:22:15* Connection #0 to host 127.0.0.1 left intact
```

This shows that the request is not blocked, it's status is 200 OK. 
The reverse proxy has also added the `X-Searchpilot-Interview: hello` header.

In `blocklist.go` change line 9 to `127.0.0.1` and run steps 2 & 3 above again

You should see the following result:
```
$ curl -k -v http://127.0.0.1:8080
*   Trying 127.0.0.1:8080...
* Connected to 127.0.0.1 (127.0.0.1) port 8080 (#0)
> GET / HTTP/1.1
> Host: 127.0.0.1:8080
> User-Agent: curl/7.71.1
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 403 Forbidden
< Date: Tue, 18 Oct 2022 20:23:49 GMT
< Content-Length: 0
<
* Connection #0 to host 127.0.0.1 left intact
```

This shows that the request was blocked, it's status is 403 Forbidden. 
The terminal window where the reverse proxy server is running will also show a log
line stating that a request from a bad IP has been blocked.

It is also possible to configure the reverse proxy to serve the `httpforever.com` website.
This can be done by:

1. Changing line 13 of `main.go` to `http://httpforever.com`

2. Starting the reverse proxy server `cd reverseproxy/ && go build && ./reverseproxy`

3. In a new terminal window making a request to the reverse proxy server using the command
`curl -v http://127.0.0.1:8080 -H "Host: httpforever.com"`
