package scanner

import (
	"context"
	"time"

	"github.com/miekg/dns"
)

func QueryAWithResponse(resolver, domain string, timeout time.Duration) (*dns.Msg, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true

	c := new(dns.Client)
	c.Net = "udp"
	c.Timeout = timeout

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r, _, err := c.ExchangeContext(ctx, m, resolver+":53")
	if err != nil {
		return nil, err
	}
	return r, nil
}
