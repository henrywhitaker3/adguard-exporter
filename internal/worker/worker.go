package worker

import (
	"context"
	"log"
	"time"

	"github.com/henrywhitaker3/adguard-exporter/internal/adguard"
	"github.com/henrywhitaker3/adguard-exporter/internal/metrics"
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
	stats, err := client.GetStats(ctx)
	if err != nil {
		log.Printf("ERROR - could not get stats: %v\n", err)
		return err
	}

	metrics.TotalQueries.WithLabelValues(client.Url()).Set(float64(stats.TotalQueries))
	metrics.BlockedFiltered.WithLabelValues(client.Url()).Set(float64(stats.BlockedFilteredQueries))
	metrics.BlockedSafesearch.WithLabelValues(client.Url()).Set(float64(stats.BlockedSafesearchQueries))

	return nil
}
