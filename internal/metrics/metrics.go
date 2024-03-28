package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TotalQueries = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "queries_total",
		Namespace: "adguard",
		Help:      "Total queries processed in the last 24 hours",
	})
)

func Init() {
	prometheus.MustRegister(TotalQueries)
}
