package scanner

import "time"

func BasicScan(ip string, tunnelDomain string, timeout time.Duration) ScanResult {
	results := make([]DomainResult, 0)

	// 1) Ping
	pingOK := PingCheck(ip, timeout)

	// 2) Domains für Klassifizierung
	classificationDomains := make([]string, 0)
	classificationDomains = append(classificationDomains, NormalDomains...)
	classificationDomains = append(classificationDomains, BlockedDomains...)

	// 3) Alle Domains, die wir wirklich testen
	testDomains := append([]string{}, classificationDomains...)
	if tunnelDomain != "" {
		testDomains = append(testDomains, tunnelDomain)
	}

	// 4) DNS Tests
	for _, d := range testDomains {
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

	// 5) NUR Normal + Blocked Domains klassifizieren
	classificationInput := make([]DomainResult, 0)
	for _, r := range results {
		if r.Domain == tunnelDomain {
			continue
		}
		classificationInput = append(classificationInput, r)
	}

	class := classify(classificationInput)

	return ScanResult{
		IP:             ip,
		PingOK:         pingOK,
		Classification: class,
		Domains:        results,
	}
}
