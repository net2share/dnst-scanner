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
	tunnel := flag.String("tunnel-domain", "", "NS subdomain to test tunnel reachability (required)")
	workers := flag.Int("workers", 50, "number of concurrent workers")
	timeoutSec := flag.Int("timeout", 3, "timeout per resolver (seconds)")
	format := flag.String("format", "json", "output format: plain|json")
	flag.Parse()

	if *tunnel == "" {
		fmt.Fprintln(os.Stderr, "--tunnel-domain is required")
		os.Exit(2)
	}
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

	results := scanner.RunBasicScanPool(ips, *workers, *tunnel, timeout)

	if *format == "json" {
		if err := json.NewEncoder(os.Stdout).Encode(results); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	} else {
		for _, r := range results {
			fmt.Printf("%s  %s\n", r.Classification, r.IP)
		}
	}

	// exit codes: success if at least one clean
	for _, r := range results {
		if r.Classification == scanner.ClassClean {
			os.Exit(0)
		}
	}
	os.Exit(1)
}
