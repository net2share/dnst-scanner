package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/net2share/dnst-scanner/internal/scanner"
)

func main() {
	input := flag.String("input", "resolvers.txt", "resolver IP list file")
	workers := flag.Int("workers", 50, "number of concurrent workers")
	timeoutSec := flag.Int("timeout", 3, "timeout per resolver (seconds)")

	step := flag.String("step", "basic", "scan step: basic|e2e")

	domain := flag.String("domain", "example.com", "domain to query (basic scan)")
	hcSlip := flag.String("slipstream-health", "", "slipstream health check domain")
	hcDNSTT := flag.String("dnstt-health", "", "dnstt health check domain")

	format := flag.String("format", "plain", "output format: plain|json")
	only := flag.String("only", "", "filter results: ok|fail (optional)")
	flag.Parse()

	if *workers <= 0 || *timeoutSec <= 0 {
		fmt.Fprintln(os.Stderr, "invalid workers or timeout")
		os.Exit(2)
	}

	timeout := time.Duration(*timeoutSec) * time.Second

	ips, err := scanner.LoadResolvers(*input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// step handling
	var domainToQuery string

	switch *step {
	case "basic":
		domainToQuery = *domain
	case "e2e":
		if *hcSlip != "" {
			domainToQuery = *hcSlip
		} else if *hcDNSTT != "" {
			domainToQuery = *hcDNSTT
		} else {
			fmt.Fprintln(os.Stderr, "e2e step requires --slipstream-health or --dnstt-health")
			os.Exit(2)
		}
	default:
		fmt.Fprintln(os.Stderr, "invalid step:", *step)
		os.Exit(2)
	}

	results := scanner.RunPool(ips, *workers, timeout, domainToQuery)

	// filter
	if *only == "ok" || *only == "fail" {
		wantOK := *only == "ok"
		filtered := make([]scanner.Result, 0)
		for _, r := range results {
			if r.OK == wantOK {
				filtered = append(filtered, r)
			}
		}
		results = filtered
	}

	// output
	if *format == "json" {
		if err := json.NewEncoder(os.Stdout).Encode(results); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	} else {
		for _, r := range results {
			if r.OK {
				fmt.Println("OK  ", r.IP)
			} else {
				fmt.Println("FAIL", r.IP)
			}
		}
	}

	// exit codes
	for _, r := range results {
		if r.OK {
			os.Exit(0)
		}
	}
	os.Exit(1)
}
