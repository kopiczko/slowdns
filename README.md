# slowdns

A simple DNS server implementation waiting for a given time (default 30s) and responding with SERVFAIL to all requests.

It helps debugging what happens when you have too slow DNS server.

Options:

```
$ slowdns -h
Usage of slowdns:
  -p int
        port to run on (default 8053)
  -w string
        time to wait before returning SERVFAIL (default "30s")
```

## Building

To build `slowdns` binary:

```
make build
```

To build `kopiczko/slowdns` docker image:

```
make docker-build
```

## Example CoreDNS configuration

Assuming all defaults in `values.yaml`:

```
    (cache-snip) {
      cache 30
    }

    (log-snip) {
      log . {
        class denial
        class error
      }
    }

    (loadbalance-snip) {
      loadbalance round_robin
    }

    my.test.io:1053 {
      import cache-snip
      forward . 172.31.0.11
      import log-snip
      import loadbalance-snip
    }
```

## Example usage

(see `dig` commands on the bottom)

```
$ docker run -p 8053:8053 --rm docker.io/kopiczko/slowdns
2024/05/10 18:52:21 Listening on port 8053
2024/05/10 18:52:21 Sleeping for 30s before responding with SERVFAIL
2024/05/10 18:52:29 Received msg[1]:
    ;; opcode: QUERY, status: NOERROR, id: 19780
    ;; flags: rd ad; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 1
    
    ;; OPT PSEUDOSECTION:
    ; EDNS: version 0; flags:; udp: 1232
    ; COOKIE: 915ca9ee4230c4bc
    
    ;; QUESTION SECTION:
    ;a.com.     IN       A
    
2024/05/10 18:52:39 Received msg[2]:
    ;; opcode: QUERY, status: NOERROR, id: 352
    ;; flags: rd ad; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 1
    
    ;; OPT PSEUDOSECTION:
    ; EDNS: version 0; flags:; udp: 1232
    ; COOKIE: 915ca9ee4230c4bc
    
    ;; QUESTION SECTION:
    ;a.com.     IN       A
    
2024/05/10 18:52:49 Received msg[3]:
    ;; opcode: QUERY, status: NOERROR, id: 24658
    ;; flags: rd ad; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 1
    
    ;; OPT PSEUDOSECTION:
    ; EDNS: version 0; flags:; udp: 1232
    ; COOKIE: 915ca9ee4230c4bc
    
    ;; QUESTION SECTION:
    ;a.com.     IN       A
    
2024/05/10 18:52:59 Responded to msg[1]
2024/05/10 18:53:09 Responded to msg[2]
2024/05/10 18:53:14 Received msg[4]:
    ;; opcode: QUERY, status: NOERROR, id: 3948
    ;; flags: rd ad; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 1
    
    ;; OPT PSEUDOSECTION:
    ; EDNS: version 0; flags:; udp: 1232
    ; COOKIE: 01c831cd4c2c3b9c
    
    ;; QUESTION SECTION:
    ;a.com.     IN       A
    
2024/05/10 18:53:19 Responded to msg[3]
2024/05/10 18:53:44 Responded to msg[4]
```

```
$ dig a.com +tcp  @127.0.0.1 -p 8053
;; communications error to 127.0.0.1#8053: timed out
;; communications error to 127.0.0.1#8053: timed out
;; communications error to 127.0.0.1#8053: timed out

; <<>> DiG 9.18.24 <<>> a.com +tcp @127.0.0.1 -p 8053
;; global options: +cmd
;; no servers could be reached
```

```
$ dig a.com +tcp  @127.0.0.1 -p 8053 +timeout=32 +tries=1

; <<>> DiG 9.18.24 <<>> a.com +tcp @127.0.0.1 -p 8053 +timeout=32 +tries=1
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: SERVFAIL, id: 3948
;; flags: qr aa rd; QUERY: 1, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 0
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;a.com.                         IN      A

;; Query time: 30002 msec
;; SERVER: 127.0.0.1#8053(127.0.0.1) (TCP)
;; WHEN: Fri May 10 19:53:44 BST 2024
;; MSG SIZE  rcvd: 23
```
