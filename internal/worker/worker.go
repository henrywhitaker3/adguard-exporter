package worker

import (
	"context"
	"log"
	"time"

	"github.com/henrywhitaker3/adguard-exporter/internal/adguard"
	"github.com/henrywhitaker3/adguard-exporter/internal/metrics"
)

func Work(ctx context.Context, interval time.Duration, clients []*adguard.Client) {
	log.Printf("Collectin metrics every %s\n", interval)
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
	log.Printf("collecting metrics for %s\n", client.Url())

	stats, err := client.GetStats(ctx)
	if err != nil {
		log.Printf("ERROR - could not get stats: %v\n", err)
		return err
	}

	metrics.TotalQueries.WithLabelValues(client.Url()).Set(float64(stats.TotalQueries))

	return nil
}
