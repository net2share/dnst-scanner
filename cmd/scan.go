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

	// basic scan flags
	input := fs.String("input", "resolvers.txt", "resolver IP list file")
	tunnel := fs.String("tunnel-domain", "", "NS subdomain to test tunnel reachability (required)")
	workers := fs.Int("workers", 50, "number of concurrent workers")
	timeoutSec := fs.Int("timeout", 3, "timeout per resolver (seconds)")
	format := fs.String("format", "json", "output format: plain|json")

	// e2e flags (Phase G – wiring only)
	e2e := fs.Bool("e2e", false, "enable E2E tunnel validation (experimental)")
	slipHealth := fs.String("slipstream-health", "", "slipstream health check domain")
	slipFP := fs.String("slipstream-fingerprint", "", "slipstream tls fingerprint")
	dnsttHealth := fs.String("dnstt-health", "", "dnstt health check domain")
	dnsttKey := fs.String("dnstt-pubkey", "", "dnstt public key")

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

	e2eCfg := scanner.E2EConfig{
		Enable:                *e2e,
		SlipstreamHealth:      *slipHealth,
		SlipstreamFingerprint: *slipFP,
		DNSTTHealth:           *dnsttHealth,
		DNSTTPubKey:           *dnsttKey,
	}

	var results []scanner.ScanResult
	if e2eCfg.Enable {
		results = scanner.RunScanPoolWithE2E(
			ips,
			*workers,
			*tunnel,
			timeout,
			e2eCfg,
		)
	} else {
		results = scanner.RunBasicScanPool(
			ips,
			*workers,
			*tunnel,
			timeout,
		)
	}

	switch *format {
	case "json":
		if err := json.NewEncoder(os.Stdout).Encode(results); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	case "plain":
		for _, r := range results {
			fmt.Printf(
				"%s  %s (ping=%v)\n",
				r.Classification,
				r.IP,
				r.PingOK,
			)
		}
	default:
		fmt.Fprintln(os.Stderr, "invalid format:", *format)
		os.Exit(2)
	}

	// exit code: success if at least one clean resolver
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
  --input <file>                 Resolver IP list (default: resolvers.txt)
  --tunnel-domain <domain>       REQUIRED – tunnel test subdomain
  --workers <n>                  Number of concurrent workers (default: 50)
  --timeout <sec>                Timeout per resolver (default: 3)
  --format <plain|json>          Output format (default: json)

E2E (experimental):
  --e2e                          Enable E2E tunnel validation
  --slipstream-health <domain>   Slipstream health check domain
  --slipstream-fingerprint <fp>  Slipstream TLS fingerprint
  --dnstt-health <domain>        DNSTT health check domain
  --dnstt-pubkey <key>           DNSTT public key
`)
}

// normalizeArgs converts "--flag value" into "--flag=value"
// Needed for Windows / PowerShell robustness.
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
