package scanner

import "time"

type Result struct {
	IP string `json:"ip"`
	OK bool   `json:"ok"`
}

func RunPool(ips []string, workers int, timeout time.Duration, domain string) []Result {
	jobs := make(chan string)
	results := make(chan Result)

	for i := 0; i < workers; i++ {
		go func() {
			for ip := range jobs {
				ok := QueryA(ip, domain, timeout)
				results <- Result{IP: ip, OK: ok}
			}
		}()
	}

	go func() {
		for _, ip := range ips {
			jobs <- ip
		}
		close(jobs)
	}()

	out := make([]Result, 0, len(ips))
	for i := 0; i < len(ips); i++ {
		out = append(out, <-results)
	}
	return out
}
