package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TotalQueries = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "queries",
		Namespace: "adguard",
		Help:      "Total queries processed in the last 24 hours",
	}, []string{"server"})
	BlockedFiltered = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "blocked_filtered",
		Namespace: "adguard",
		Help:      "Total queries that have been blocked from filter lists",
	}, []string{"server"})
	BlockedSafesearch = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "blocked_safesearch",
		Namespace: "adguard",
		Help:      "Total queries that have been blocked due to safesearch",
	}, []string{"server"})
)

func Init() {
	prometheus.MustRegister(TotalQueries)
	prometheus.MustRegister(BlockedFiltered)
	prometheus.MustRegister(BlockedSafesearch)
}
