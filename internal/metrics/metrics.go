package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TotalQueries = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "queries_total",
		Namespace: "adguard",
		Help:      "Total queries processed in the last 24 hours",
	}, []string{"server"})
)

func Init() {
	prometheus.MustRegister(TotalQueries)
}
