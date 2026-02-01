package scanner

import (
	"net"
	"testing"

	"github.com/miekg/dns"
)

func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		ip     string
		expect bool
	}{
		{"10.1.2.3", true},
		{"192.168.1.1", true},
		{"172.16.0.5", true},
		{"8.8.8.8", false},
		{"1.1.1.1", false},
	}

	for _, tt := range tests {
		ip := net.ParseIP(tt.ip)
		if isPrivateIP(ip) != tt.expect {
			t.Fatalf("ip %s expected %v", tt.ip, tt.expect)
		}
	}
}

func TestAnalyzeResponse_HijackDetected(t *testing.T) {
	msg := new(dns.Msg)
	msg.Answer = append(msg.Answer, &dns.A{
		Hdr: dns.RR_Header{Name: "facebook.com.", Rrtype: dns.TypeA},
		A:   net.ParseIP("10.0.0.1"),
	})

	res := analyzeResponse("facebook.com", msg)
	if !res.Hijacked {
		t.Fatal("expected hijack to be detected")
	}
}
