package adguard

import "time"

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
	Mac      string     `json:"mac"`
	IP       string     `json:"ip"`
	Hostname string     `json:"hostname"`
	Expires  *time.Time `json:"expires,omitempty"`
	Type     string
}

type DhcpStatus struct {
	Enabled       Bool        ` json:"enabled"`
	DynamicLeases []DhcpLease `json:"leases"`
	StaticLeases  []DhcpLease `json:"static_leases"`
	Leases        []DhcpLease
}

type query struct {
	Class string `json:"class"`
	Host  string `json:"name"`
	Type  string `json:"type"`
}

type answer struct {
	Type  string  `json:"type"`
	TTL   float64 `json:"ttl"`
	Value any     `json:"value"`
}

type dnsHeader struct {
	Name     string `json:"Name"`
	Rrtype   int    `json:"Rrtype"`
	Class    int    `json:"Class"`
	TTL      int    `json:"Ttl"`
	Rdlength int    `json:"Rdlength"`
}

type type65 struct {
	Hdr   dnsHeader `json:"Hdr"`
	RData string    `json:"Rdata"`
}

type ClientInfo struct {
	Whois          map[string]any `json:"whois"`
	Name           string         `json:"name"`
	DisallowedRule string         `json:"disallowed_rule"`
	Disallowed     Bool           `json:"disallowed"`
}

type logEntry struct {
	Answer      []answer   `json:"answer"`
	DNSSec      Bool       `json:"answer_dnssec"`
	Client      string     `json:"client"`
	ClientProto string     `json:"client_proto"`
	Elapsed     string     `json:"elapsedMs"`
	Question    query      `json:"question"`
	Reason      string     `json:"reason"`
	Status      string     `json:"status"`
	Time        string     `json:"time"`
	Upstream    string     `json:"upstream"`
	ClientInfo  ClientInfo `json:"client_info"`
}

type QueryTime struct {
	Elapsed  time.Duration
	Client   string
	Upstream string
}

type queryLog struct {
	Log    []logEntry `json:"data"`
	Oldest string     `json:"oldest"`
}
