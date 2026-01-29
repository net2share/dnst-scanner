package scanner

import (
	"time"

	"github.com/miekg/dns"
)

const (
	retryAttempts = 3
	retryBaseDelay = 200 * time.Millisecond
)

func ResolveWithRetry(resolver, domain string, timeout time.Duration) (*dns.Msg, error) {
	var lastErr error

	for attempt := 1; attempt <= retryAttempts; attempt++ {
		msg, err := QueryAWithResponse(resolver, domain, timeout)
		if err == nil && msg != nil && len(msg.Answer) > 0 {
			return msg, nil
		}

		lastErr = err
		time.Sleep(time.Duration(attempt) * retryBaseDelay)
	}

	return nil, lastErr
}
