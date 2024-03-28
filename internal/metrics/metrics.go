package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ScrapeErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "scrape_errors_total",
		Namespace: "adguard",
		Help:      "The number of errors scraping a target",
	}, []string{"server"})

	// Status
	ProtectionEnabled = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "protection_enabled",
		Namespace: "adguard",
		Help:      "Whether DNS filtering is enabled",
	}, []string{"server"})
	Running = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "running",
		Namespace: "adguard",
		Help:      "Whether adguard is running or not",
	}, []string{"server", "version"})

	// Stats
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
	BlockedSafebrowsing = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "blocked_safebrowsing",
		Namespace: "adguard",
		Help:      "Total queries that have been blocked due to safebrowsing",
	}, []string{"server"})
	AvgProcessingTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "avg_processing_time_seconds",
		Namespace: "adguard",
		Help:      "The average query processing time in seconds",
	}, []string{"server"})
	TopQueriedDomains = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "top_queried_domains",
		Namespace: "adguard",
		Help:      "The number of queries for the top domains",
	}, []string{"server", "domain"})
	TopBlockedDomains = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "top_blocked_domains",
		Namespace: "adguard",
		Help:      "The number of blocked queries for the top domains",
	}, []string{"server", "domain"})
	TopClients = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "top_clients",
		Namespace: "adguard",
		Help:      "The number of queries for the top clients",
	}, []string{"server", "client"})
	TopUpstreams = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "top_upstreams",
		Namespace: "adguard",
		Help:      "The number of repsonses for the top upstream servers",
	}, []string{"server", "upstream"})
	TopUpstreamsAvgTimes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "top_upstreams_avg_response_time_seconds",
		Namespace: "adguard",
		Help:      "The average response time for each of the top upstream servers",
	}, []string{"server", "upstream"})

	// DHCP
	DhcpEnabled = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "dhcp_enabled",
		Namespace: "adguard",
		Help:      "Whether dhcp is enabled",
	}, []string{"server"})
	DhcpLeases = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "dhcp_leases",
		Namespace: "adguard",
		Help:      "The number of dhcp leases",
	}, []string{"server"})
	DhcpStaticLeases = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "dhcp_static_leases",
		Namespace: "adguard",
		Help:      "The number of static dhcp leases",
	}, []string{"server"})
)

func Init() {
	prometheus.MustRegister(ScrapeErrors)

	// Stats
	prometheus.MustRegister(TotalQueries)
	prometheus.MustRegister(BlockedFiltered)
	prometheus.MustRegister(BlockedSafesearch)
	prometheus.MustRegister(BlockedSafebrowsing)
	prometheus.MustRegister(AvgProcessingTime)
	prometheus.MustRegister(TopBlockedDomains)
	prometheus.MustRegister(TopClients)
	prometheus.MustRegister(TopQueriedDomains)
	prometheus.MustRegister(TopUpstreams)
	prometheus.MustRegister(TopUpstreamsAvgTimes)

	// Status
	prometheus.MustRegister(Running)
	prometheus.MustRegister(ProtectionEnabled)

	// DHCP
	prometheus.MustRegister(DhcpEnabled)
	prometheus.MustRegister(DhcpLeases)
	prometheus.MustRegister(DhcpStaticLeases)
}
