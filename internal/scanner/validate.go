package scanner

import "net"

// IsValidIP validates IPv4 / IPv6 literals
func IsValidIP(ip string) bool {
	parsed := net.ParseIP(ip)
	return parsed != nil
}
