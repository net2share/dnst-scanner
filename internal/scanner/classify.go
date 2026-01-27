package scanner

func classify(results []DomainResult) Classification {
	for _, r := range results {
		if !r.Resolved {
			return ClassBroken
		}
		if r.Hijacked {
			return ClassCensored
		}
	}
	return ClassClean
}
