package scanner

type Classification string

const (
	ClassClean    Classification = "clean"
	ClassCensored Classification = "censored"
	ClassBroken   Classification = "broken"
)

type DomainResult struct {
	Domain   string `json:"domain"`
	Resolved bool   `json:"resolved"`
	Hijacked bool   `json:"hijacked"`
}

type E2EResult struct {
	SlipstreamOK bool   `json:"slipstream_ok,omitempty"`
	DNSTTOK      bool   `json:"dnstt_ok,omitempty"`
	Error        string `json:"error,omitempty"`
}

type ScanResult struct {
	IP             string         `json:"ip"`
	PingOK         bool           `json:"ping_ok"`
	Classification Classification  `json:"classification"`
	Domains        []DomainResult `json:"domains"`
	E2E            *E2EResult     `json:"e2e,omitempty"`
}
