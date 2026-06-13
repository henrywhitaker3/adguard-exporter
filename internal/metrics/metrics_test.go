package metrics

import (
	"testing"

	"github.com/henrywhitaker3/adguard-exporter/internal/adguard"
	"github.com/prometheus/client_golang/prometheus"
)

// AdGuard can return duplicate lease rows (in-place self-update rebuilds the
// lease table; clients behind a MAC-NAT'ing repeater collapse onto one bridge
// MAC). Before the dedup fix this poisoned the whole scrape with a hard 500.
func TestDhcpLeasesCollectDeduplicates(t *testing.T) {
	server := "http://192.168.3.2"
	dup := adguard.DhcpLease{
		Mac:      "2a:fc:97:00:a6:e0",
		IP:       "192.168.3.114",
		Hostname: "iphone",
		Type:     "dynamic",
	}
	distinct := adguard.DhcpLease{
		Mac:      "0a:8b:f4:09:7d:63",
		IP:       "192.168.3.113",
		Hostname: "watch",
		Type:     "dynamic",
	}

	d := NewDhcpLeasesServer(DhcpLeasesMetric)
	// Same label tuple appears twice, plus one distinct lease.
	d.Record(server, []adguard.DhcpLease{dup, dup, distinct})

	reg := prometheus.NewPedanticRegistry()
	if err := reg.Register(d); err != nil {
		t.Fatalf("register: %v", err)
	}

	// Gather mirrors what the /metrics handler does. Before the fix this
	// returned the "was collected before with the same name and label values" error.
	mfs, err := reg.Gather()
	if err != nil {
		t.Fatalf("Gather returned error (duplicate leases not deduped): %v", err)
	}

	for _, mf := range mfs {
		if mf.GetName() != "adguard_dhcp_leases" {
			continue
		}
		if got := len(mf.GetMetric()); got != 2 {
			t.Fatalf("expected 2 deduplicated lease series, got %d", got)
		}
		return
	}
	t.Fatal("adguard_dhcp_leases metric family not found")
}
