package checkd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Port is the metrics listener port (default 8666).
var Port = 8666

// check represents a check, containing a check function and a run interval.
type check struct {
	check    func()
	interval time.Duration
}

// checks is the list of checks to be run.
var checks = []check{}

// gauges is the map of registered gauges.
var gauges = map[string]prometheus.Gauge{}

// counters is the map of registered counters.
var counters = map[string]prometheus.Counter{}

// Every schedules the specified check function to run at the specified interval.
func Every(interval time.Duration, checkFunc func()) {
	checks = append(checks, check{checkFunc, interval})
}

// Start runs all checks concurrently.
func Start() {
	log.Printf("starting %d checks", len(checks))
	for _, c := range checks {
		go func(c check) {
			for {
				c.check()
				time.Sleep(c.interval)
			}
		}(c)
	}
	log.Printf("starting metrics listener on port %d", Port)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Port), nil))
}

// Gauge registers and returns a new Prometheus gauge metric. If the gauge has
// already been registered, Gauge returns the existing gauge.
func Gauge(name, help string) prometheus.Gauge {
	if _, ok := gauges[name]; !ok {
		gauges[name] = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: name,
			Help: help,
		})
		prometheus.MustRegister(gauges[name])
	}
	return gauges[name]
}

// Counter registers and returns a new Prometheus counter metric. If the counter has
// already been registered, Counter returns the existing counter.
func Counter(name, help string) prometheus.Counter {
	if _, ok := counters[name]; !ok {
		counters[name] = prometheus.NewCounter(prometheus.CounterOpts{
			Name: name,
			Help: help,
		})
		prometheus.MustRegister(counters[name])
	}
	return counters[name]
}
