package scanner

import (
	"context"
	"time"

	"github.com/miekg/dns"
)

func QueryA(resolver, domain string, timeout time.Duration) bool {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true

	c := new(dns.Client)
	c.Net = "udp"
	c.Timeout = timeout

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r, _, err := c.ExchangeContext(ctx, m, resolver+":53")
	if err != nil || r == nil {
		return false
	}
	if r.Rcode != dns.RcodeSuccess {
		return false
	}
	return len(r.Answer) > 0
}
