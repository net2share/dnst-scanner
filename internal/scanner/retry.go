package scanner

import (
	"errors"
	"math/rand"
	"time"

	"github.com/miekg/dns"
)

const (
	maxAttempts = 3
	baseDelay   = 150 * time.Millisecond
)

func ResolveWithRetry(ip, domain string, timeout time.Duration) (*dns.Msg, error) {
	var lastErr error

	for i := 0; i < maxAttempts; i++ {
		msg, err := resolveOnce(ip, domain, timeout)
		if err == nil {
			return msg, nil
		}
		lastErr = err

		jitter := time.Duration(rand.Int63n(int64(baseDelay)))
		time.Sleep(baseDelay + jitter)
	}

	return nil, lastErr
}

func resolveOnce(ip, domain string, timeout time.Duration) (*dns.Msg, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	c := new(dns.Client)
	c.Timeout = timeout

	in, _, err := c.Exchange(m, ip+":53")
	if err != nil {
		return nil, err
	}
	if in == nil || len(in.Answer) == 0 {
		return nil, errors.New("no answer")
	}

	return in, nil
}
