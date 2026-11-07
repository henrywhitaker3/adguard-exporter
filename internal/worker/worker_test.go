package worker

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/henrywhitaker3/adguard-exporter/internal/adguard"
	"github.com/henrywhitaker3/adguard-exporter/internal/config"
	"github.com/henrywhitaker3/adguard-exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectStatusRetainsRunningMetricsForEveryServerAfterVersionChanges(t *testing.T) {
	resetStatusState(t)

	first := &statusResponse{Version: "v0.107.79", Running: true, ProtectionEnabled: true}
	second := &statusResponse{Version: "v0.107.74", Running: true, ProtectionEnabled: true}
	firstClient := newStatusClient(t, first)
	secondClient := newStatusClient(t, second)

	collectStatus(context.Background(), firstClient)
	collectStatus(context.Background(), secondClient)

	first.Version = "v0.107.80"
	second.Version = "v0.107.80"
	collectStatus(context.Background(), firstClient)
	collectStatus(context.Background(), secondClient)

	assert.Equal(t, map[string]runningMetric{
		firstClient.Url():  {Version: "v0.107.80", Value: 1},
		secondClient.Url(): {Version: "v0.107.80", Value: 1},
	}, collectRunningMetrics(t))
}

func TestCollectStatusSetsRunningToZeroWhenStatusRequestFails(t *testing.T) {
	resetStatusState(t)

	response := &statusResponse{Version: "v0.107.80", Running: true, ProtectionEnabled: true}
	client := newStatusClient(t, response)

	collectStatus(context.Background(), client)
	response.StatusCode = http.StatusServiceUnavailable
	collectStatus(context.Background(), client)

	assert.Equal(t, map[string]runningMetric{
		client.Url(): {Version: "v0.107.80", Value: 0},
	}, collectRunningMetrics(t))
}

func TestCollectStatusReplacesUnknownVersionWhenStatusRequestRecovers(t *testing.T) {
	resetStatusState(t)

	response := &statusResponse{StatusCode: http.StatusServiceUnavailable}
	client := newStatusClient(t, response)

	collectStatus(context.Background(), client)
	response.StatusCode = 0
	response.Version = "v0.107.80"
	response.Running = true
	response.ProtectionEnabled = true
	collectStatus(context.Background(), client)

	assert.Equal(t, map[string]runningMetric{
		client.Url(): {Version: "v0.107.80", Value: 1},
	}, collectRunningMetrics(t))
}

func TestCollectStatusSupportsConcurrentCollections(t *testing.T) {
	resetStatusState(t)

	clients := make([]*adguard.Client, 8)
	for i := range clients {
		clients[i] = newStatusClient(t, &statusResponse{
			Version:           "v0.107.80",
			Running:           true,
			ProtectionEnabled: true,
		})
	}

	start := make(chan struct{})
	var group sync.WaitGroup
	for _, client := range clients {
		group.Go(func() {
			<-start
			collectStatus(context.Background(), client)
		})
	}
	close(start)
	group.Wait()

	assert.Len(t, collectRunningMetrics(t), len(clients))
}

func TestCollectQueryLogStatsReplacesStaleSeriesForEachServer(t *testing.T) {
	resetQueryLogMetrics(t)

	firstResponse := &queryLogResponse{Data: []queryLogEntryResponse{newQueryLogEntry("first-old")}}
	secondResponse := &queryLogResponse{Data: []queryLogEntryResponse{newQueryLogEntry("second-current")}}
	firstClient := newQueryLogClient(t, firstResponse)
	secondClient := newQueryLogClient(t, secondResponse)

	collectQueryLogStats(context.Background(), firstClient)
	collectQueryLogStats(context.Background(), secondClient)

	firstResponse.Data = []queryLogEntryResponse{newQueryLogEntry("first-current")}
	collectQueryLogStats(context.Background(), firstClient)

	expected := map[string]int{firstClient.Url(): 1, secondClient.Url(): 1}
	assert.Equal(t, expected, collectMetricSeriesByServer(t, metrics.TotalQueriesDetails))
	assert.Equal(t, expected, collectMetricSeriesByServer(t, metrics.TotalQueriesDetailsHistogram))
	assert.Equal(t, expected, collectMetricSeriesByServer(t, metrics.ProcessingTimeBucketMilli))
	assert.Equal(t, expected, collectMetricSeriesByServer(t, metrics.ProcessingTimeBucket))
}

type statusResponse struct {
	Version           string `json:"version"`
	Running           bool   `json:"running"`
	ProtectionEnabled bool   `json:"protection_enabled"`
	StatusCode        int    `json:"-"`
}

type runningMetric struct {
	Version string
	Value   float64
}

type queryLogResponse struct {
	Data []queryLogEntryResponse `json:"data"`
}

type queryLogEntryResponse struct {
	Client      string `json:"client"`
	ClientProto string `json:"client_proto"`
	Elapsed     string `json:"elapsedMs"`
	Reason      string `json:"reason"`
	Status      string `json:"status"`
	Upstream    string `json:"upstream"`
	ClientInfo  struct {
		Name string `json:"name"`
	} `json:"client_info"`
}

func newStatusClient(t *testing.T, response *statusResponse) *adguard.Client {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if response.StatusCode != 0 {
			http.Error(w, "unavailable", response.StatusCode)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	return adguard.NewClient(config.Config{Url: server.URL})
}

func newQueryLogClient(t *testing.T, response *queryLogResponse) *adguard.Client {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	return adguard.NewClient(config.Config{Url: server.URL})
}

func collectRunningMetrics(t *testing.T) map[string]runningMetric {
	t.Helper()

	registry := prometheus.NewRegistry()
	require.NoError(t, registry.Register(metrics.Running))
	metricFamilies, err := registry.Gather()
	require.NoError(t, err)

	running := map[string]runningMetric{}
	for _, metric := range metricFamilies[0].Metric {
		labels := map[string]string{}
		for _, label := range metric.Label {
			labels[label.GetName()] = label.GetValue()
		}
		running[labels["server"]] = runningMetric{
			Version: labels["version"],
			Value:   metric.GetGauge().GetValue(),
		}
	}

	return running
}

func resetStatusState(t *testing.T) {
	t.Helper()

	versions = map[string]string{}
	metrics.Running.Reset()
	metrics.ProtectionEnabled.Reset()
	metrics.ScrapeErrors.Reset()
	t.Cleanup(func() {
		versions = map[string]string{}
		metrics.Running.Reset()
		metrics.ProtectionEnabled.Reset()
		metrics.ScrapeErrors.Reset()
	})
}

func newQueryLogEntry(client string) queryLogEntryResponse {
	entry := queryLogEntryResponse{
		Client:      client,
		ClientProto: "udp",
		Elapsed:     "3.75",
		Reason:      "NotFilteredNotFound",
		Status:      "NOERROR",
		Upstream:    "https://dns.example/dns-query",
	}
	entry.ClientInfo.Name = client

	return entry
}

func collectMetricSeriesByServer(t *testing.T, collector prometheus.Collector) map[string]int {
	t.Helper()

	registry := prometheus.NewRegistry()
	require.NoError(t, registry.Register(collector))
	metricFamilies, err := registry.Gather()
	require.NoError(t, err)

	series := map[string]int{}
	for _, metric := range metricFamilies[0].Metric {
		for _, label := range metric.Label {
			if label.GetName() == "server" {
				series[label.GetValue()]++
			}
		}
	}

	return series
}

func resetQueryLogMetrics(t *testing.T) {
	t.Helper()

	metrics.TotalQueriesDetails.Reset()
	metrics.TotalQueriesDetailsHistogram.Reset()
	metrics.ProcessingTimeBucketMilli.Reset()
	metrics.ProcessingTimeBucket.Reset()
	t.Cleanup(func() {
		metrics.TotalQueriesDetails.Reset()
		metrics.TotalQueriesDetailsHistogram.Reset()
		metrics.ProcessingTimeBucketMilli.Reset()
		metrics.ProcessingTimeBucket.Reset()
	})
}
