package scanner

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

const defaultResolversURL = "https://raw.githubusercontent.com/net2share/ir-resolvers/main/resolvers.txt"

func LoadResolvers(path string) ([]string, error) {
	// 1) ENV: local file override
	if p := os.Getenv("DNST_SCANNER_RESOLVERS_PATH"); p != "" {
		f, err := os.Open(p)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return parseScanner(f)
	}

	// 2) explicit --input
	if path != "" {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return parseScanner(f)
	}

	// 3) ENV URL override or default
	url := defaultResolversURL
	if u := os.Getenv("DNST_SCANNER_RESOLVERS_URL"); u != "" {
		url = u
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to fetch resolvers list")
	}

	return parseScanner(resp.Body)
}

func parseScanner(r io.Reader) ([]string, error) {
	sc := bufio.NewScanner(r)
	out := make([]string, 0)

	for sc.Scan() {
		ip := strings.TrimSpace(sc.Text())
		if ip == "" || strings.HasPrefix(ip, "#") {
			continue
		}
		if !IsValidIP(ip) {
			continue
		}
		out = append(out, ip)
	}

	return out, sc.Err()
}
