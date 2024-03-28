package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	DomainsBlocked = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "domains_being_blocked",
			Namespace: "adguard",
			Help:      "This represents the number of domains being blocked",
		},
		[]string{"hostname"},
	)

	DNSQueriesToday = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "dns_queries_today",
			Namespace: "adguard",
			Help:      "This represents the number of DNS queries made over the last 24 hours",
		},
		[]string{"hostname"},
	)

	AdsBlockedToday = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "ads_blocked_today",
			Namespace: "adguard",
			Help:      "This represents the number of ads blocked over the last 24 hours",
		},
		[]string{"hostname"},
	)

	AdsPercentageToday = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "ads_percentage_today",
			Namespace: "adguard",
			Help:      "This represents the percentage of ads blocked over the last 24 hours",
		},
		[]string{"hostname"},
	)

	UniqueDomains = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "unique_domains",
			Namespace: "adguard",
			Help:      "This represents the number of unique domains seen",
		},
		[]string{"hostname"},
	)

	QueriesForwarded = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "queries_forwarded",
			Namespace: "adguard",
			Help:      "This represents the number of queries forwarded",
		},
		[]string{"hostname"},
	)

	QueriesCached = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "queries_cached",
			Namespace: "adguard",
			Help:      "This represents the number of queries cached",
		},
		[]string{"hostname"},
	)

	ClientsEverSeen = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "clients_ever_seen",
			Namespace: "adguard",
			Help:      "This represents the number of clients ever seen",
		},
		[]string{"hostname"},
	)

	UniqueClients = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "unique_clients",
			Namespace: "adguard",
			Help:      "This represents the number of unique clients seen",
		},
		[]string{"hostname"},
	)

	DNSQueriesAllTypes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "dns_queries_all_types",
			Namespace: "adguard",
			Help:      "This represents the number of DNS queries made for all types",
		},
		[]string{"hostname"},
	)

	Reply = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "reply",
			Namespace: "adguard",
			Help:      "This represents the number of replies made for all types",
		},
		[]string{"hostname", "type"},
	)

	TopQueries = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "top_queries",
			Namespace: "adguard",
			Help:      "This represents the number of top queries made by Adguard by domain",
		},
		[]string{"hostname", "domain"},
	)

	TopAds = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "top_ads",
			Namespace: "adguard",
			Help:      "This represents the number of top ads made by Pi-hole by domain",
		},
		[]string{"hostname", "domain"},
	)

	TopSources = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "top_sources",
			Namespace: "adguard",
			Help:      "This represents the number of top sources requests made by Adguard by source host",
		},
		[]string{"hostname", "source"},
	)

	ForwardDestinations = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "forward_destinations",
			Namespace: "adguard",
			Help:      "This represents the number of forward destinations requests made by Adguard by destination",
		},
		[]string{"hostname", "destination"},
	)

	QueryTypes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "querytypes",
			Namespace: "adguard",
			Help:      "This represents the number of queries made by Adguard by type",
		},
		[]string{"hostname", "type"},
	)

	Status = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "status",
			Namespace: "adguard",
			Help:      "This if Adguard is enabled",
		},
		[]string{"hostname"},
	)
)

// Init initializes all Prometheus metrics made available by Pi-hole exporter.
func Init() {
	initMetric("domains_blocked", DomainsBlocked)
	initMetric("dns_queries_today", DNSQueriesToday)
	initMetric("ads_blocked_today", AdsBlockedToday)
	initMetric("ads_percentag_today", AdsPercentageToday)
	initMetric("unique_domains", UniqueDomains)
	initMetric("queries_forwarded", QueriesForwarded)
	initMetric("queries_cached", QueriesCached)
	initMetric("clients_ever_seen", ClientsEverSeen)
	initMetric("unique_clients", UniqueClients)
	initMetric("dns_queries_all_types", DNSQueriesAllTypes)
	initMetric("reply", Reply)
	initMetric("top_queries", TopQueries)
	initMetric("top_ads", TopAds)
	initMetric("top_sources", TopSources)
	initMetric("forward_destinations", ForwardDestinations)
	initMetric("querytypes", QueryTypes)
	initMetric("status", Status)
}

func initMetric(name string, metric *prometheus.GaugeVec) {
	prometheus.MustRegister(metric)
	log.Info("New Prometheus metric registered: ", name)
}
