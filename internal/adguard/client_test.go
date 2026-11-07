package adguard

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/henrywhitaker3/adguard-exporter/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStats(t *testing.T) {
	client := newTestClient(t, fixtureResponse(t, "stats.json", func(r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/control/stats", r.URL.Path)
		assert.Empty(t, r.URL.RawQuery)
	}))

	stats, err := client.GetStats(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 1042, stats.TotalQueries)
	assert.Equal(t, 51, stats.BlockedFilteredQueries)
	assert.Equal(t, float32(0.0042), stats.AvgProcessingTime)
	assert.Equal(t, []map[string]int{{"192.0.2.10": 100}, {"2001:db8::10": 20}}, stats.TopClients)
}

func TestGetStatus(t *testing.T) {
	client := newTestClient(t, fixtureResponse(t, "status.json", func(r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/control/status", r.URL.Path)
	}))

	status, err := client.GetStatus(context.Background())

	require.NoError(t, err)
	assert.True(t, bool(status.ProtectionEnabled))
	assert.True(t, bool(status.Running))
	assert.Equal(t, "v0.107.63", status.Version)
}

func TestGetDhcpCombinesAndLabelsLeases(t *testing.T) {
	client := newTestClient(t, fixtureResponse(t, "dhcp-status.json", func(r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/control/dhcp/status", r.URL.Path)
	}))

	status, err := client.GetDhcp(context.Background())

	require.NoError(t, err)
	assert.True(t, bool(status.Enabled))
	require.Len(t, status.Leases, 2)
	assert.Equal(t, "dynamic", status.Leases[0].Type)
	assert.Equal(t, "laptop", status.Leases[0].Hostname)
	assert.Equal(t, "static", status.Leases[1].Type)
	assert.Equal(t, "printer", status.Leases[1].Hostname)
	require.NotNil(t, status.Leases[0].Expires)
	assert.Equal(t, "2026-11-08T12:00:00Z", status.Leases[0].Expires.Format(time.RFC3339))
}

func TestGetQueryLogBuildsTypesAndTimings(t *testing.T) {
	client := newTestClient(t, fixtureResponse(t, "querylog.json", func(r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/control/querylog", r.URL.Path)
		assert.Equal(t, "limit=1000&response_status=all", r.URL.RawQuery)
	}))

	types, timings, entries, err := client.GetQueryLog(context.Background())

	require.NoError(t, err)
	assert.Equal(t, map[string]map[string]int{
		"192.0.2.10":   {"A": 1, "TYPE65": 1},
		"2001:db8::10": {"AAAA": 1},
	}, types)
	require.Len(t, timings, 2)
	assert.Equal(t, QueryTime{Elapsed: 3750 * time.Microsecond, Client: "192.0.2.10", Upstream: "https://dns.example/dns-query"}, timings[0])
	assert.Equal(t, QueryTime{Elapsed: 2 * time.Millisecond, Client: "2001:db8::10", Upstream: "self"}, timings[1])
	require.Len(t, entries, 3)
	assert.Equal(t, "laptop", entries[0].ClientInfo.Name)
}

func TestSearchClients(t *testing.T) {
	client := newTestClient(t, fixtureResponse(t, "clients-search.json", func(r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/control/clients/search", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"clients":[{"id":"192.0.2.10"},{"id":"2001:db8::10"}]}`, string(body))
	}))

	clients, err := client.SearchClients(context.Background(), []map[string]int{
		{"192.0.2.10": 100},
		{"2001:db8::10": 20},
	})

	require.NoError(t, err)
	assert.Equal(t, map[string]string{
		"192.0.2.10":   "laptop",
		"laptop-id":    "laptop",
		"2001:db8::10": "phone",
	}, clients)
}

func TestClientReturnsResponseErrors(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "unavailable", http.StatusServiceUnavailable)
		})

		_, err := client.GetStats(context.Background())

		require.Error(t, err)
		assert.ErrorContains(t, err, "unexpected status code 503: unavailable")
	})

	t.Run("post", func(t *testing.T) {
		client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "invalid client", http.StatusBadRequest)
		})

		_, err := client.SearchClients(context.Background(), []map[string]int{{"192.0.2.10": 1}})

		require.Error(t, err)
		assert.ErrorContains(t, err, "unexpected status code 400: invalid client")
	})
}

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "exporter", username)
		assert.Equal(t, "password", password)
		handler(w, r)
	}))
	t.Cleanup(server.Close)

	return NewClient(config.Config{
		Url:      server.URL,
		Username: "exporter",
		Password: "password",
	})
}

func fixtureResponse(t *testing.T, name string, assertRequest func(*http.Request)) http.HandlerFunc {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)

	return func(w http.ResponseWriter, r *http.Request) {
		assertRequest(r)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(body)
		assert.NoError(t, err)
	}
}
