package main

import "fmt"

//var (
//	location = prometheus.NewGaugeVec(
//		prometheus.GaugeOpts{
//			Name: "ns_exporter_tcp_established",
//			Help: "Namespace tcp sockets established",
//		},[]string{"container_name", "ip", "port"},
//	)
//)
//
//func main() {
//	addr := flag.String("web.listen-address", ":9300", "Address on which to expose metrics")
//	interval := flag.Int("interval", 10, "Interval fo metrics collection in seconds")
//	flag.Parse()
//	prometheus.MustRegister(location)
//	http.Handle("/metrics", prometheus.Handler())
//	go run(int(*interval))
//	log.Fatal(http.ListenAndServe(*addr, nil))
//
//}
//
//func run(interval int)  {
//	for {
//		location.Reset()
//
//		for code, number := range c.ConnectionsByCode {
//			location.With(prometheus.Labels{"location": code}).Set(float64(number))
//		}
//
//		time.Sleep(time.Duration(interval) * time.Second)
//	}
//}

func main()  {
	c := getContainers()
	fmt.Println(c)
}
