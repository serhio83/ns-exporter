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
	tcpConnEstabCollector = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ns_exporter_tcp_established",
			Help: "Namespace tcp sockets established",
		}, []string{"container_name", "ip", "port"},
	)
	tcpConnTWCollector = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ns_exporter_tcp_timewait",
			Help: "Namespace tcp sockets time-wait",
		}, []string{"container_name", "ip", "port"},
	)
)

func main() {
	addr := flag.String("web.listen-address", ":9515", "Address on which to expose metrics")
	interval := flag.Int("interval", 10, "Interval for metrics collection in seconds")
	flag.Parse()
	containers := &Containers{}
	prometheus.MustRegister(tcpConnEstabCollector)
	prometheus.MustRegister(tcpConnTWCollector)
	log.Println("Listen on:", *addr)
	http.Handle("/metrics", promhttp.Handler())
	go run(int(*interval), containers)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func run(interval int, containers *Containers) {
	for {
		containers.getContainers()
		tcpConnEstabCollector.Reset()
		tcpConnTWCollector.Reset()

		// collect tcp connections
		for _, container := range containers.contList {

			estabCount := map[string]float64{}
			twCount := map[string]float64{}

			for _, tcpEstab := range container.tcpConnEstab {
				for ip, port := range tcpEstab {
					if _, ok := estabCount[strconv.Itoa(port)+ip]; ok {
						estabCount[strconv.Itoa(port)+ip] += 1
					} else {
						estabCount[strconv.Itoa(port)+ip] = 1
					}
					tcpConnEstabCollector.With(prometheus.Labels{"container_name": container.name, "ip": ip, "port": strconv.Itoa(port)}).Set(estabCount[strconv.Itoa(port)+ip])
				}
			}

			for _, tcpTW := range container.tcpConnTW {
				for ip, port := range tcpTW {
					if _, ok := twCount[strconv.Itoa(port)+ip]; ok {
						twCount[strconv.Itoa(port)+ip] += 1
					} else {
						twCount[strconv.Itoa(port)+ip] = 1
					}
					tcpConnTWCollector.With(prometheus.Labels{"container_name": container.name, "ip": ip, "port": strconv.Itoa(port)}).Set(twCount[strconv.Itoa(port)+ip])
				}
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
