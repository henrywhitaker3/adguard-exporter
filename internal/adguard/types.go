package adguard

type Bool bool

func (b Bool) Int() int {
	if b {
		return 1
	}
	return 0
}

type Stats struct {
	TotalQueries               int                  `json:"num_dns_queries"`
	BlockedFilteredQueries     int                  `json:"num_blocked_filtering"`
	BlockedSafebrowsingQueries int                  `json:"num_blocked_safebrowsing"`
	BlockedSafesearchQueries   int                  `json:"num_blocked_safesearch"`
	BlockedParentalQueries     int                  `json:"num_blocked_parental"`
	AvgProcessingTime          float32              `json:"avg_processing_time"`
	TopQueriedDomains          []map[string]int     `json:"top_queired_domains"`
	TopBlockedDomains          []map[string]int     `json:"top_blocked_domains"`
	TopClients                 []map[string]int     `json:"top_clients"`
	TopUpstreamsResponses      []map[string]int     `json:"top_upstreams_responses"`
	TopUpstreamsAvgTimes       []map[string]float32 `json:"top_upstreams_avg_time"`
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
}

type DhcpStatus struct {
	Enabled      Bool        ` json:"enabled"`
	Leases       []DhcpLease `json:"leases"`
	StaticLeases []DhcpLease `json:"static_leases"`
}
