package scanner

import (
	"net"

	"github.com/miekg/dns"
)

func analyzeResponse(domain string, msg *dns.Msg) DomainResult {
	res := DomainResult{
		Domain:   domain,
		Resolved: false,
		Hijacked: false,
	}

	if msg == nil || len(msg.Answer) == 0 {
		return res
	}

	res.Resolved = true

	for _, ans := range msg.Answer {
		if a, ok := ans.(*dns.A); ok {
			ip := a.A
			if isPrivateIP(ip) {
				res.Hijacked = true
				return res
			}
		}
	}

	return res
}

func isPrivateIP(ip net.IP) bool {
	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	for _, cidr := range privateRanges {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}
	return false
}
