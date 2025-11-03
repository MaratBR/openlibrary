package olhttp

import (
	"net"
	"net/http"
	"strings"
)

func GetIP(r *http.Request) net.IP {
	forwardedFor := r.Header.Get("x-forwarded-for")
	if forwardedFor == "" {
		return net.ParseIP(r.RemoteAddr)
	} else {
		ips := strings.Split(forwardedFor, ",")
		ip := strings.Trim(ips[0], " ")
		return net.ParseIP(ip)
	}
}
