package cmd

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"os"
	"strings"
)

const defaultResolversURL = "https://raw.githubusercontent.com/net2share/ir-resolvers/main/resolvers.txt"

func loadResolvers() ([]string, error) {
	// 1) ENV: Path
	if p := os.Getenv("DNST_SCANNER_RESOLVERS_PATH"); p != "" {
		return loadFromFile(p)
	}

	// 2) ENV: URL
	if u := os.Getenv("DNST_SCANNER_RESOLVERS_URL"); u != "" {
		return loadFromURL(u)
	}

	// 3) Default URL
	return loadFromURL(defaultResolversURL)
}

func loadFromFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parseIPs(f)
}

func loadFromURL(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch resolvers")
	}

	return parseIPs(resp.Body)
}

func parseIPs(r interface{}) ([]string, error) {
	var s *bufio.Scanner

	switch v := r.(type) {
	case *os.File:
		s = bufio.NewScanner(v)
	case *strings.Reader:
		s = bufio.NewScanner(v)
	case *http.Response:
		s = bufio.NewScanner(v.Body)
	default:
		// generic reader
		s = bufio.NewScanner(r.(interface {
			Read([]byte) (int, error)
		}))
	}

	var ips []string
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if net.ParseIP(line) == nil {
			continue
		}
		ips = append(ips, line)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return ips, nil
}
