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

// Register adds a new check to the list of checks that will be run.
func Register(checkFunc func(), interval time.Duration) {
	log.Printf("registering check %v at interval %v\n", checkFunc, interval)
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
	if _, ok := gauges[name]; ok {
		return gauges[name]
	}
	gauges[name] = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	})
	log.Printf("registering gauge %q", name)
	if err := prometheus.Register(gauges[name]); err != nil {
		log.Fatalf("failed to register gauge %q: %v", name, err)
	}
	return gauges[name]
}

// Counter registers and returns a new Prometheus counter metric. If the counter has
// already been registered, Counter returns the existing counter.
func Counter(name, help string) prometheus.Counter {
	if _, ok := counters[name]; ok {
		return counters[name]
	}
	counters[name] = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: help,
	})
	log.Printf("registering counter %q", name)
	if err := prometheus.Register(counters[name]); err != nil {
		log.Fatalf("failed to register counter %q: %v", name, err)
	}
	return counters[name]
}
