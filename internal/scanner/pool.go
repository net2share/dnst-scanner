package scanner

import (
	"sync"
	"time"
)

func RunBasicScanPool(ips []string, workers int, tunnelDomain string, timeout time.Duration) []ScanResult {
	return runPool(ips, workers, timeout, func(ip string) ScanResult {
		return BasicScan(ip, tunnelDomain, timeout)
	})
}

func RunScanPoolWithE2E(
	ips []string,
	workers int,
	tunnelDomain string,
	timeout time.Duration,
	e2eCfg E2EConfig,
) []ScanResult {
	return runPool(ips, workers, timeout, func(ip string) ScanResult {
		res := BasicScan(ip, tunnelDomain, timeout)
		if e2eCfg.Enable {
			res.E2E = RunE2E(ip, timeout, e2eCfg)
		}
		return res
	})
}

func runPool(
	ips []string,
	workers int,
	timeout time.Duration,
	fn func(string) ScanResult,
) []ScanResult {
	in := make(chan string)
	out := make(chan ScanResult)

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range in {
				out <- fn(ip)
			}
		}()
	}

	go func() {
		for _, ip := range ips {
			in <- ip
		}
		close(in)
		wg.Wait()
		close(out)
	}()

	results := make([]ScanResult, 0, len(ips))
	for r := range out {
		results = append(results, r)
	}
	return results
}
