package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace   = "jupyterhub"
	metricsPath = "/metrics"
)

type myCollector struct{}

var (
	activeUserDesc = prometheus.NewDesc(
		"active_user",
		"Current active users.",
		[]string{"userName"}, nil,
	)
)

func (cc myCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cc, ch)
}

func (cc *myCollector) GetActiveUser() (
	activeUsers map[string]int,
) {
	activeUsers = map[string]int{
		"test":  0,
		"test2": 1,
	}
	return
}

func (cc myCollector) Collect(ch chan<- prometheus.Metric) {
	activeUsers := cc.GetActiveUser()

	for userName, lastActivity := range activeUsers {
		ch <- prometheus.MustNewConstMetric(
			activeUserDesc,
			prometheus.CounterValue,
			float64(lastActivity),
			userName,
		)
	}
}

func main() {
	reg := prometheus.NewPedanticRegistry()
	cc := myCollector{}
	reg.MustRegister(cc)

	http.Handle(metricsPath, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>Jupyterhub Exporter</title></head>
			<body>
			<h1>Jupyterhub Exporter</h1>
			<p><a href="` + metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
