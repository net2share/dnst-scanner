package scanner

import "time"

func BasicScan(ip string, tunnelDomain string, timeout time.Duration) ScanResult {
	results := make([]DomainResult, 0)

	pingOK := PingCheck(ip, timeout)

	allDomains := append([]string{}, NormalDomains...)
	allDomains = append(allDomains, BlockedDomains...)
	allDomains = append(allDomains, tunnelDomain)

	for _, d := range allDomains {
		msg, err := ResolveWithRetry(ip, d, timeout)
		if err != nil {
			results = append(results, DomainResult{
				Domain:   d,
				Resolved: false,
				Hijacked: false,
			})
			continue
		}
		results = append(results, analyzeResponse(d, msg))
	}

	class := classify(results)

	return ScanResult{
		IP:             ip,
		PingOK:         pingOK,
		Classification: class,
		Domains:        results,
	}
}
