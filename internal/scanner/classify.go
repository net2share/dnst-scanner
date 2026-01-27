package scanner

func classify(results []DomainResult) Classification {
	for _, r := range results {
		// nur normale + blocked Domains berücksichtigen
		if isTunnelDomain(r.Domain) {
			continue
		}

		if !r.Resolved {
			return ClassBroken
		}
		if r.Hijacked {
			return ClassCensored
		}
	}
	return ClassClean
}

func isTunnelDomain(domain string) bool {
	return domain != "google.com" &&
		domain != "microsoft.com" &&
		domain != "facebook.com" &&
		domain != "x.com"
}
