package scanner

import (
	"net"
	"time"
)

func PingCheck(ip string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("udp", ip+":53", timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
