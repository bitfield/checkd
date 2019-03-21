package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/bitfield/checkd"
	"github.com/spf13/viper"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

type quotaChecker struct {
	sync.Mutex
	project  string
	gcp      *compute.Service
	metrics  map[string]prometheus.Gauge
	interval time.Duration
}

func (q *quotaChecker) Name() string {
	return "quota"
}

func (q *quotaChecker) Interval() time.Duration {
	return q.interval
}

func (q *quotaChecker) Init(v *viper.Viper) (err error) {
	q.project = v.GetString("project")
	if q.project == "" {
		return fmt.Errorf("quota: project must be set")
	}
	v.SetDefault("interval", "24hr")
	q.interval = v.GetDuration("interval")
	google, err := google.DefaultClient(context.Background(), compute.CloudPlatformScope)
	if err != nil {
		return fmt.Errorf("quota: couldn't create Google client: %v", err)
	}
	q.gcp, err = compute.New(google)
	if err != nil {
		return fmt.Errorf("quota: couldn't create compute client: %v", err)
	}
	p, err := q.gcp.Projects.Get(q.project).Do()
	if err != nil {
		return fmt.Errorf("quota: failed to get project info for %s: %v", q.project, err)
	}
	q.metrics = map[string]prometheus.Gauge{}
	log.Println("quota: registering metrics")
	for _, m := range p.Quotas {
		q.metrics[m.Metric] = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: fmt.Sprintf("czquota_%s_ratio", strings.ToLower(m.Metric)),
			Help: fmt.Sprintf("Ratio of quota usage to limit for %q", m.Metric),
		})
		if err = prometheus.Register(q.metrics[m.Metric]); err != nil {
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				return fmt.Errorf("quota: failed to register metric: %t %v", err, err)
			}
		}
	}
	return nil
}

func (q *quotaChecker) Check() error {
	log.Println("quota: updating quotas")
	p, err := q.gcp.Projects.Get(q.project).Do()
	if err != nil {
		return fmt.Errorf("quota: failed to get project info for %s: %v", q.project, err)
	}
	q.Lock()
	for _, m := range p.Quotas {
		quotaPercentUsed := m.Usage / m.Limit
		q.metrics[m.Metric].Set(quotaPercentUsed)
	}
	q.Unlock()
	log.Println("quota: update complete")
	return nil
}

func init() {
	checkd.Register(&quotaChecker{})
}
