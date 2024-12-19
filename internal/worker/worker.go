package worker

import (
	"context"
	"log"
	"slices"
	"strconv"
	"time"

	"github.com/henrywhitaker3/adguard-exporter/internal/adguard"
	"github.com/henrywhitaker3/adguard-exporter/internal/metrics"
)

var (
	initialised = []string{}
	versions    = map[string]string{}
)

func Work(ctx context.Context, interval time.Duration, clients []*adguard.Client) {
	log.Printf("Collecting metrics every %s\n", interval)
	tick := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			for _, c := range clients {
				go collect(ctx, c)
			}
		}
	}
}

func collect(ctx context.Context, client *adguard.Client) error {
	// Initialise the scrape errors counter with a 0
	if !slices.Contains(initialised, client.Url()) {
		metrics.ScrapeErrors.WithLabelValues(client.Url())
		initialised = append(initialised, client.Url())
	}

	go collectStats(ctx, client)
	go collectStatus(ctx, client)
	go collectDhcp(ctx, client)
	go collectQueryLogStats(ctx, client)

	return nil
}

func collectStats(ctx context.Context, client *adguard.Client) {
	stats, err := client.GetStats(ctx)
	if err != nil {
		log.Printf("ERROR - could not get stats: %v\n", err)
		metrics.ScrapeErrors.WithLabelValues(client.Url()).Inc()
		return
	}
	metrics.TotalQueries.WithLabelValues(client.Url()).Set(float64(stats.TotalQueries))
	metrics.BlockedFiltered.WithLabelValues(client.Url()).Set(float64(stats.BlockedFilteredQueries))
	metrics.ReplacedSafesearch.WithLabelValues(client.Url()).Set(float64(stats.ReplacedSafesearchQueries))
	metrics.ReplacedSafebrowsing.WithLabelValues(client.Url()).Set(float64(stats.ReplacedSafebrowsingQueries))
	metrics.ReplacedParental.WithLabelValues(client.Url()).Set(float64(stats.ReplacedParentalQueries))
	metrics.AvgProcessingTime.WithLabelValues(client.Url()).Set(float64(stats.AvgProcessingTime))

	for _, c := range stats.TopClients {
		for key, val := range c {
			metrics.TopClients.WithLabelValues(client.Url(), key).Set(float64(val))
		}
	}
	for _, c := range stats.TopUpstreamsResponses {
		for key, val := range c {
			metrics.TopUpstreams.WithLabelValues(client.Url(), key).Set(float64(val))
		}
	}
	for _, c := range stats.TopQueriedDomains {
		for key, val := range c {
			metrics.TopQueriedDomains.WithLabelValues(client.Url(), key).Set(float64(val))
		}
	}
	for _, c := range stats.TopBlockedDomains {
		for key, val := range c {
			metrics.TopBlockedDomains.WithLabelValues(client.Url(), key).Set(float64(val))
		}
	}
	for _, c := range stats.TopUpstreamsAvgTimes {
		for key, val := range c {
			metrics.TopUpstreamsAvgTimes.WithLabelValues(client.Url(), key).Set(float64(val))
		}
	}
}

func collectStatus(ctx context.Context, client *adguard.Client) {
	status, err := client.GetStatus(ctx)
	if err != nil {
		log.Printf("ERROR - could not get status: %v\n", err)
		metrics.ScrapeErrors.WithLabelValues(client.Url()).Inc()
		return
	}
	// Persist the running version the first time
	if _, ok := versions[client.Url()]; !ok {
		versions[client.Url()] = status.Version
	}

	// Check if the adguard version has changed
	if versions[client.Url()] != status.Version {
		metrics.Running.Reset()
	}

	metrics.Running.WithLabelValues(client.Url(), status.Version).Set(float64(status.Running.Int()))
	metrics.ProtectionEnabled.WithLabelValues(client.Url()).Set(float64(status.ProtectionEnabled.Int()))
}

func collectDhcp(ctx context.Context, client *adguard.Client) {
	dhcp, err := client.GetDhcp(ctx)
	if err != nil {
		log.Printf("ERROR - could not get dhcp status: %v\n", err)
		metrics.ScrapeErrors.WithLabelValues(client.Url()).Inc()
		return
	}
	metrics.DhcpEnabled.WithLabelValues(client.Url()).Set(float64(dhcp.Enabled.Int()))
	metrics.DhcpLeases.Record(client.Url(), dhcp.Leases)
}

func collectQueryLogStats(ctx context.Context, client *adguard.Client) {
	stats, times, queries, err := client.GetQueryLog(ctx)
	if err != nil {
		log.Printf("ERROR - could not get query type stats: %v\n", err)
		metrics.ScrapeErrors.WithLabelValues(client.Url()).Inc()
		return
	}

	for c, v := range stats {
		for t, v := range v {
			metrics.QueryTypes.WithLabelValues(client.Url(), t, c).Set(float64(v))
		}
	}

	for _, l := range queries {
		elapsed, err := strconv.ParseFloat(l.Elapsed, 64)
		if err != nil {
			continue
		}
		metrics.TotalQueriesDetails.WithLabelValues(client.Url(), l.Client, l.Reason, l.Status, l.Upstream, l.ClientInfo.Name).Set(elapsed)
		metrics.TotalQueriesDetailsHistogram.WithLabelValues(client.Url(), l.Client, l.Reason, l.Status, l.Upstream, l.ClientInfo.Name).Observe(float64(elapsed))
	}

	for _, t := range times {
		metrics.ProcessingTimeBucketMilli.
			WithLabelValues(client.Url(), t.Client, t.Upstream).
			Observe(float64(t.Elapsed.Milliseconds()))
		metrics.ProcessingTimeBucket.
			WithLabelValues(client.Url(), t.Client, t.Upstream).
			Observe(t.Elapsed.Seconds())
	}
}
