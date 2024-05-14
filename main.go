package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/miekg/dns"
)

const rcode = dns.RcodeServerFailure

func main() {
	port := flag.Int("p", 8053, "port to run on")
	wait := flag.String("w", "30s", "time to wait before returning SERVFAIL")
	flag.Parse()

	sleepFor, err := time.ParseDuration(*wait)
	if err != nil {
		log.Fatalf("Can't parse %q (-w flag) as duration: %s\n", *wait, err)
	}

	go func() {
		srv := &dns.Server{Addr: ":" + strconv.Itoa(*port), Net: "udp"}
		srv.Handler = newMux("UDP", sleepFor)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to set udp listener %s\n", err.Error())
		}
	}()

	go func() {
		srv := &dns.Server{Addr: ":" + strconv.Itoa(*port), Net: "tcp"}
		srv.Handler = newMux("TCP", sleepFor)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Failed to set tcp listener %s\n", err.Error())
		}
	}()

	log.Printf("Listening on port %d\n", *port)
	log.Printf("Sleeping for %s before responding with %s\n", *wait, dns.RcodeToString[rcode])

	sig := make(chan os.Signal, 1)
	defer close(sig)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Printf("Signal (%v) received, stopping\n", s)
}

func newMux(proto string, sleepFor time.Duration) *dns.ServeMux {
	handleIDSeq := &atomic.Int64{}

	mux := dns.NewServeMux()

	mux.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		id := handleIDSeq.Add(1)
		log.Printf("Received msg[%s:%d]:\n%s\n", proto, id, indent(4, r.String()))
		m := new(dns.Msg)
		m.SetRcode(r, rcode)
		m.Authoritative = true
		time.Sleep(sleepFor)
		w.WriteMsg(m)
		log.Printf("Responded to msg[%d]\n", id)
	})

	return mux
}

func indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}
