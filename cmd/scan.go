package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/net2share/dnst-scanner/internal/scanner"
)

func runScan(args []string) {
	args = normalizeArgs(args)

	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	input := fs.String("input", "resolvers.txt", "resolver IP list file")
	tunnel := fs.String("tunnel-domain", "", "NS subdomain to test tunnel reachability (required)")
	workers := fs.Int("workers", 50, "number of concurrent workers")
	timeoutSec := fs.Int("timeout", 3, "timeout per resolver (seconds)")
	format := fs.String("format", "json", "output format: plain|json")

	if err := fs.Parse(args); err != nil {
		usageScan()
		os.Exit(2)
	}

	if *tunnel == "" {
		fmt.Fprintln(os.Stderr, "--tunnel-domain is required")
		usageScan()
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

	switch *format {
	case "json":
		if err := json.NewEncoder(os.Stdout).Encode(results); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	case "plain":
		for _, r := range results {
			fmt.Printf("%s  %s (ping=%v)\n", r.Classification, r.IP, r.PingOK)
		}
	default:
		fmt.Fprintln(os.Stderr, "invalid format:", *format)
		os.Exit(2)
	}

	for _, r := range results {
		if r.Classification == scanner.ClassClean {
			os.Exit(0)
		}
	}
	os.Exit(1)
}

func usageScan() {
	fmt.Println(`
Usage:
  dnst-scanner scan --tunnel-domain <domain> [flags]

Flags:
  --input <file>           Resolver IP list (default: resolvers.txt)
  --tunnel-domain <domain> REQUIRED – tunnel test subdomain
  --workers <n>            Number of concurrent workers (default: 50)
  --timeout <sec>          Timeout per resolver (default: 3)
  --format <plain|json>    Output format (default: json)
`)
}

// normalizeArgs converts "--flag value" to "--flag=value"
func normalizeArgs(args []string) []string {
	out := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") && !strings.Contains(arg, "=") {
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				out = append(out, arg+"="+args[i+1])
				i++
				continue
			}
		}
		out = append(out, arg)
	}
	return out
}
