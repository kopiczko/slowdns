FROM alpine:3.19

WORKDIR /app

ADD slowdns /app/slowdns

ENTRYPOINT ["/app/slowdns"]
