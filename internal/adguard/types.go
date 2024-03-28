package adguard

type Bool bool

func (b Bool) Int() int {
	if b {
		return 1
	}
	return 0
}

type Stats struct {
	TotalQueries                int                  `json:"num_dns_queries"`
	BlockedFilteredQueries      int                  `json:"num_blocked_filtering"`
	ReplacedSafebrowsingQueries int                  `json:"num_replaced_safebrowsing"`
	ReplacedSafesearchQueries   int                  `json:"num_replaced_safesearch"`
	ReplacedParentalQueries     int                  `json:"num_replaced_parental"`
	AvgProcessingTime           float32              `json:"avg_processing_time"`
	TopQueriedDomains           []map[string]int     `json:"top_queried_domains"`
	TopBlockedDomains           []map[string]int     `json:"top_blocked_domains"`
	TopClients                  []map[string]int     `json:"top_clients"`
	TopUpstreamsResponses       []map[string]int     `json:"top_upstreams_responses"`
	TopUpstreamsAvgTimes        []map[string]float32 `json:"top_upstreams_avg_time"`
}

type Status struct {
	ProtectionEnabled Bool   `json:"protection_enabled"`
	Version           string `json:"version"`
	Running           Bool   `json:"running"`
}

type DhcpLease struct {
	Mac      string `json:"mac"`
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	Expires  string `json:"expires,omitempty"`
	Type     string
}

type DhcpStatus struct {
	Enabled       Bool        ` json:"enabled"`
	DynamicLeases []DhcpLease `json:"leases"`
	StaticLeases  []DhcpLease `json:"static_leases"`
	Leases        []DhcpLease
}
