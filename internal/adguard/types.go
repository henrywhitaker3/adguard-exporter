package adguard

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
