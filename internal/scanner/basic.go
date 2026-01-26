package scanner

import (
	"net"
	"time"
)

func TestResolver(ip string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("udp", ip+":53", timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
