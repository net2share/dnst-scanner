# DNS Tunnel Resolver Scanner (dnst-scanner)

A tool to scan and identify recursive DNS resolvers compatible with DNS tunneling. Provides end-to-end validation for finding resolver IPs that can establish DNS tunnels like Slipstream and DNSTT.

## Features

- Fetch raw resolver IP list from [ir-resolvers](https://github.com/net2share/ir-resolvers)
- Two-step scanning process:
  1. **Basic scan**: Test if resolvers respond to standard DNS queries
  2. **E2E validation**: Verify resolvers can establish actual DNS tunnels
- Concurrent scanning with configurable parallelism
- Output working resolvers in various formats (plain list, JSON)
- Standalone CLI tool, orchestrated by [dnstc](https://github.com/net2share/dnstc)

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           dnst-scanner                                      │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                        Step 1: Basic Scan                           │    │
│  │                                                                     │    │
│  │   ir-resolvers ──► Raw IP List ──► DNS Query Test ──► Responding   │    │
│  │   (GitHub)          (10k+ IPs)     (A record query)    Resolvers   │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                     Step 2: E2E Validation                          │    │
│  │                                                                     │    │
│  │   Responding ──► Tunnel Test ──► Health Check ──► Tunnel-Capable   │    │
│  │   Resolvers      (via resolver)   Endpoint        Resolvers        │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
                                     ▼
                        ┌───────────────────────┐
                        │  Server-Side (dnstm)  │
                        │                       │
                        │  Health Check Instances│
                        │  • Slipstream endpoint │
                        │  • DNSTT endpoint      │
                        └───────────────────────┘
```

## Scanning Process

### Step 1: Basic Scan

Tests each resolver IP with a simple DNS query:
- Send A record query for a known domain
- Check if resolver responds within timeout
- Filter out non-responding or invalid resolvers

**Input**: Raw IP list from ir-resolvers (~10k+ IPs)
**Output**: List of responding recursive resolvers

### Step 2: E2E Validation

Tests responding resolvers with actual tunnel protocols:
- Send DNS query through resolver to health check endpoint on server
- Health check domains hosted on dnstm (e.g., `hc-s.example.com` for Slipstream)
- Verify complete tunnel path works: client → resolver → server → response

**Input**: Responding resolvers from Step 1
**Output**: Tunnel-capable resolvers ready for use

## Usage

```bash
# Basic scan only
dnst-scanner scan --step basic --output resolvers.txt

# Full scan with E2E validation
dnst-scanner scan --step e2e \
  --slipstream-health hc-s.example.com \
  --dnstt-health hc-d.example.com \
  --output working-resolvers.json

# Scan with custom resolver list
dnst-scanner scan --input custom-ips.txt --step e2e

# Scan with concurrency control
dnst-scanner scan --workers 100 --timeout 5s
```

## Configuration

| Option | Description | Default |
|--------|-------------|---------|
| `--input` | Custom resolver IP list file | Fetch from ir-resolvers |
| `--step` | Scan step: `basic` or `e2e` | `e2e` |
| `--workers` | Number of concurrent workers | 50 |
| `--timeout` | Timeout per resolver | 3s |
| `--output` | Output file path | stdout |
| `--format` | Output format: `plain` or `json` | `plain` |
| `--slipstream-health` | Slipstream health check domain | - |
| `--dnstt-health` | DNSTT health check domain | - |

## Integration with dnstc

dnstc orchestrates dnst-scanner as a subprocess:
- dnstc runs dnst-scanner with appropriate flags
- Scanner outputs JSON to stdout
- dnstc parses results and updates resolver pool
- Scheduled periodic runs keep resolver list fresh

```bash
# Example: dnstc runs scanner and captures JSON output
dnst-scanner scan --format json --slipstream-health hc-s.example.com
```

## Requirements

- Windows, macOS, or Linux
- Network access to target resolvers
- Server with health check endpoints configured (for E2E validation)

## Related Projects

- [dnstc](https://github.com/net2share/dnstc) - DNS tunnel client (uses this scanner)
- [dnstm](https://github.com/net2share/dnstm) - DNS tunnel server (hosts health check endpoints)
- [ir-resolvers](https://github.com/net2share/ir-resolvers) - Raw resolver IP list
- [go-corelib](https://github.com/net2share/go-corelib) - Shared Go library
