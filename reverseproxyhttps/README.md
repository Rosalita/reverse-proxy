# Reverse Proxy HTTPS

This code is an evolution of the code in the `reverseproxy` folder.

I found this article on stack overflow:
https://stackoverflow.com/questions/35390726/confirm-tls-certificate-while-performing-reverseproxy-in-golang

This example uses https://github.com/chromium/badssl.com which hosts `https://sha256.badssl.com/`

I have modified my reverse proxy so that it can proxy `https://sha256.badssl.com/` based on the stack overflow example.

To run this code
`cd reverseproxyhttps/ && go build && ./reverseproxyhttps`

The open a browser and navigate to `http://localhost:8443/`
You should see the page loaded from `https://sha256.badssl.com/`
