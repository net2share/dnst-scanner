# DNS Tunnel Resolver Scanner (dnst-scanner)

A tool to scan and identify recursive DNS resolvers compatible with DNS tunneling. Provides end-to-end validation for finding resolver IPs that can establish DNS tunnels like Slipstream and DNSTT.

## Features

- Fetch raw resolver IP list from [ir-resolvers](https://github.com/net2share/ir-resolvers)
- Two-step scanning process:
  1. **Basic scan**: Ping check, multi-domain DNS testing, resolver classification
  2. **E2E validation** (optional): Test with actual Slipstream/DNSTT tunnel connections
- Multi-domain testing to identify resolver behavior:
  - Normal domains (google.com, microsoft.com)
  - Blocked domains (facebook.com, x.com) to detect censorship/hijacking
  - Custom tunnel domain to verify reachability
- Classification: `clean` (resolves blocked domains) vs `censored` (hijacks to 10.x.x.x)
- Concurrent scanning with configurable parallelism
- Output in JSON or plain text formats
- Standalone CLI tool, orchestrated by [dnstc](https://github.com/net2share/dnstc)

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           dnst-scanner                                      │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                        Step 1: Basic Scan                           │    │
│  │                                                                     │    │
│  │   ir-resolvers ──► Raw IPs ──► Ping ──► DNS Tests ──► Classified   │    │
│  │   (GitHub)                      Check    (multi-domain)  Resolvers │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                    │                                        │
│                                    ▼                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │                 Step 2: E2E Validation (Optional)                   │    │
│  │                                                                     │    │
│  │   Classified  ──► Slipstream/DNSTT ──► Health Check ──► Verified   │    │
│  │   Resolvers       Client Test           Endpoint       Resolvers   │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Scanning Process

### Step 1: Basic Scan

1. **Ping check**: Verify basic connectivity
2. **Normal domain queries**: Test google.com and microsoft.com (baseline)
3. **Blocked domain query**: Test facebook.com/x.com to detect hijacking
4. **Tunnel domain query**: Test custom NS subdomain reachability

**Output**: Classified resolvers with response times and behavior data

### Step 2: E2E Validation (Optional)

Tests resolvers with actual tunnel connections:

- Requires Slipstream/DNSTT client binaries
- Connects through each resolver to health check endpoint
- Verifies complete tunnel path works

**Output**: Resolvers marked with E2E success/failure per protocol

## Usage

```bash
# Basic scan only (default)
dnst-scanner scan --tunnel-domain t.example.com

# Basic scan + E2E validation
dnst-scanner scan --tunnel-domain t.example.com --e2e \
  --slipstream-health hc-s.example.com \
  --slipstream-fingerprint abc123 \
  --dnstt-health hc-d.example.com \
  --dnstt-pubkey xyz789

# Custom resolver list
dnst-scanner scan --input custom-ips.txt --tunnel-domain t.example.com

# JSON output to file
dnst-scanner scan --tunnel-domain t.example.com --format json --output results.json
```

## Configuration

| Option                     | Description                               | Default                 |
| -------------------------- | ----------------------------------------- | ----------------------- |
| `--input`                  | Custom resolver IP list file              | Fetch from ir-resolvers |
| `--tunnel-domain`          | NS subdomain to test tunnel reachability  | Required                |
| `--e2e`                    | Enable E2E validation with actual tunnels | false                   |
| `--slipstream-health`      | Slipstream health check domain (for E2E)  | -                       |
| `--slipstream-fingerprint` | Slipstream TLS fingerprint (for E2E)      | -                       |
| `--dnstt-health`           | DNSTT health check domain (for E2E)       | -                       |
| `--dnstt-pubkey`           | DNSTT public key (for E2E)                | -                       |
| `--workers`                | Number of concurrent workers              | 50                      |
| `--timeout`                | Timeout per resolver                      | 3s                      |
| `--output`                 | Output file path                          | stdout                  |
| `--format`                 | Output format: `plain` or `json`          | `json`                  |

### Environment Variable Overrides

| Variable                       | Description                       |
| ------------------------------ | --------------------------------- |
| `DNST_SCANNER_RESOLVERS_URL`   | Override default ir-resolvers URL |
| `DNST_SCANNER_RESOLVERS_PATH`  | Use local file (skips download)   |
| `DNST_SCANNER_SLIPSTREAM_PATH` | Path to slipstream-client binary  |
| `DNST_SCANNER_DNSTT_PATH`      | Path to dnstt-client binary       |

## Integration with dnstc

dnstc orchestrates dnst-scanner as a subprocess:

- dnstc runs dnst-scanner with appropriate flags
- Scanner outputs JSON to stdout
- dnstc parses results and updates resolver pool
- Scheduled periodic runs keep resolver list fresh

```bash
# Example: dnstc runs scanner and captures JSON output
dnst-scanner scan --tunnel-domain t.example.com --format json
```

## Requirements

- Windows, macOS, or Linux
- Network access to target resolvers
- Server with health check endpoints configured (for E2E validation)
- Slipstream/DNSTT client binaries (for E2E validation)

## Related Projects

- [dnstc](https://github.com/net2share/dnstc) - DNS tunnel client (orchestrates this scanner)
- [dnstm](https://github.com/net2share/dnstm) - DNS tunnel server (hosts health check endpoints)
- [ir-resolvers](https://github.com/net2share/ir-resolvers) - Raw resolver IP list
- [go-corelib](https://github.com/net2share/go-corelib) - Shared Go library

### Slipstream E2E (Platform Note)

Slipstream E2E validation requires the Slipstream client binary.
At the moment, official Slipstream client binaries are only available for Linux.

- Windows: E2E will report a clear error if the binary is not present
- Linux / CI / WSL: Full Slipstream E2E validation is supported

This is expected behavior.
