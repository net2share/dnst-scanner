package scanner

func classify(results []DomainResult) Classification {
	for _, r := range results {
		// Tunnel-Domain ignorieren
		if isTunnelDomain(r.Domain) {
			continue
		}

		// Normale Domain nicht auflösbar => broken
		if !r.Resolved {
			return ClassBroken
		}

		// Hijack => censored
		if r.Hijacked {
			return ClassCensored
		}
	}
	return ClassClean
}

// Aktuell simple Heuristik, README-konform.
// Später ersetzbar durch Domain-Typen.
func isTunnelDomain(domain string) bool {
	switch domain {
	case "google.com", "microsoft.com", "facebook.com", "x.com":
		return false
	default:
		return true
	}
}
