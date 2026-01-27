package scanner

import "time"

func RunBasicScanPool(
	ips []string,
	workers int,
	tunnelDomain string,
	timeout time.Duration,
) []ScanResult {
	jobs := make(chan string)
	results := make(chan ScanResult)

	for i := 0; i < workers; i++ {
		go func() {
			for ip := range jobs {
				results <- BasicScan(ip, tunnelDomain, timeout)
			}
		}()
	}

	go func() {
		for _, ip := range ips {
			jobs <- ip
		}
		close(jobs)
	}()

	out := make([]ScanResult, 0, len(ips))
	for i := 0; i < len(ips); i++ {
		out = append(out, <-results)
	}
	return out
}
