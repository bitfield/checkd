package checkd

import (
	"bytes"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func TestRunChecks(t *testing.T) {
	checks = []check{}
	var func1Runs, func2Runs int
	Logger.SetOutput(ioutil.Discard)
	Every(2*time.Millisecond, func() { func1Runs++ })
	Every(6*time.Millisecond, func() { func2Runs++ })
	go Start()
	time.Sleep(5 * time.Millisecond)
	if func1Runs < 2 {
		t.Errorf("want func1 to run at least 2 times, got %d", func1Runs)
	}
	time.Sleep(10 * time.Millisecond)
	if func2Runs < 2 {
		t.Errorf("want func2 to run at least 2 times, got %d", func2Runs)
	}
}

func TestGauge(t *testing.T) {
	gauges = map[string]prometheus.Gauge{}
	Gauge("test_set_gauge", "")
	if _, ok := gauges["test_set_gauge"]; !ok {
		t.Fatalf("gauge not registered")
	}
	Gauge("test_set_gauge", "").Set(1)
	if len(gauges) != 1 {
		t.Fatalf("gauge not cached")
	}
}

func TestGaugeVec(t *testing.T) {
	gaugevecs = map[string]prometheus.GaugeVec{}
	GaugeVec("test_set_gaugevec", "", []string{"testLabel"})
	if _, ok := gaugevecs["test_set_gaugevec"]; !ok {
		t.Fatalf("gaugevec not registered")
	}
	g := GaugeVec("test_set_gaugevec", "", []string{"testLabel"})
	g.WithLabelValues("foo").Set(1)
	if len(gaugevecs) != 1 {
		t.Fatalf("gauge not cached")
	}
}

func TestCounter(t *testing.T) {
	counters = map[string]prometheus.Counter{}
	Counter("test_set_counter", "")
	if _, ok := counters["test_set_counter"]; !ok {
		t.Fatalf("counter not registered")
	}
	Counter("test_set_counter", "").Inc()
	if len(gauges) != 1 {
		t.Fatalf("counter not cached")
	}
}

func TestLogAndDebug(t *testing.T) {
	buf := &bytes.Buffer{}
	Logger = *log.New(buf, "", 0)
	Logger.Printf("foo")
	if buf.String() != "foo\n" {
		t.Fatalf("want %q, got %q", "foo\n", buf.String())
	}
	buf2 := &bytes.Buffer{}
	Logger = *log.New(buf2, "", 0)
	Debug("foo")
	if buf2.String() != "" {
		t.Fatalf("expected no log output with debugging disabled, got %q", buf2.String())
	}
	EnableDebug()
	Debug("foo")
	if buf2.String() != "foo\n" {
		t.Fatalf("want %q with debugging enabled, got %q", "foo\n", buf2.String())
	}
}
