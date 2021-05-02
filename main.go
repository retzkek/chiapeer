package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen", ":9133", "The address to listen on for HTTP requests.")

func main() {
	flag.Parse()

	if err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Subsystem: "chia",
			Name:      "peers_count",
			Help:      "Number of peers currently connected.",
		},
		countPeers,
	)); err != nil {
		log.Fatal(err)
	}

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func countPeers() float64 {
	c := exec.Command("sh", "-c", "cat /proc/net/tcp | grep :20FC | wc -l")
	o, err := c.CombinedOutput()
	if err != nil {
		log.Printf("error counting peers from /proc/net/tcp: %s", err)
		return -1.0
	}

	p, err := strconv.ParseFloat(string(bytes.TrimSpace(o)), 64)
	if err != nil {
		log.Printf("error converting %s to float: %s", o, err)
		return -1.0
	}
	// magic number: six local connections?
	return p - 6.0
}
