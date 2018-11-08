package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	tcpconn = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ns_exporter_tcp_established",
			Help: "Namespace tcp sockets established",
		},[]string{"container_name", "ip", "port"},
	)
)

func main() {
	addr := flag.String("web.listen-address", ":9515", "Address on which to expose metrics")
	interval := flag.Int("interval", 10, "Interval for metrics collection in seconds")
	flag.Parse()
	prometheus.MustRegister(tcpconn)
	log.Println("Listen on:", *addr)
	http.Handle("/metrics", promhttp.Handler())
	go run(int(*interval))
	log.Fatal(http.ListenAndServe(*addr, nil))

}

func run(interval int)  {
	for {
		tcpconn.Reset()
		c := getContainers()

		// collect tcp connections metrics
		for _, cont := range c.List {
			count := map[string]float64{}
			for _, skt := range cont.Sockets {
				for ip, port := range skt {
					if _, ok := count[strconv.Itoa(port) + ip]; ok {
						count[strconv.Itoa(port) + ip] += 1
					} else {
						count[strconv.Itoa(port) + ip] = 1
					}
					tcpconn.With(prometheus.Labels{"container_name": cont.Name, "ip": ip,"port": strconv.Itoa(port)}).Set(count[strconv.Itoa(port) + ip])
				}
			}
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}
}
