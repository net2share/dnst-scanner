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

type ScanResult struct {
	IP             string         `json:"ip"`
	PingOK         bool           `json:"ping_ok"`
	Classification Classification  `json:"classification"`
	Domains        []DomainResult `json:"domains"`
}
