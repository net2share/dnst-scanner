package scanner

import "time"

func BasicScan(ip string, tunnelDomain string, timeout time.Duration) ScanResult {
	results := make([]DomainResult, 0)

	allDomains := append([]string{}, NormalDomains...)
	allDomains = append(allDomains, BlockedDomains...)
	allDomains = append(allDomains, tunnelDomain)

	for _, d := range allDomains {
		msg, err := QueryAWithResponse(ip, d, timeout)
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
		Classification: class,
		Domains:        results,
	}
}
