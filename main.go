package main

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var onlineUsers = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "goapp_online_users",
	Help: "Online users",
	ConstLabels: map[string]string{
		"course": "fullcycle",
	},
})

var httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "goapp_http_requests_total",
	Help: "Count of all HTTP requests for goapp",
}, []string{})

func main() {
	r := prometheus.NewRegistry()
	r.MustRegister(onlineUsers)
	r.MustRegister(httpRequestsTotal)

	go func() {
		for {
			onlineUsers.Set(float64(rand.Intn(2000)))
		}
	}()

	home := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello FullCycle"))
	})

	http.Handle("/", promhttp.InstrumentHandlerCounter(httpRequestsTotal, home))
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":8181", nil))
}
