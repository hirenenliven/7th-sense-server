package middleware

import (
	"net"
	"net/http"
	"net/netip"
	"strings"
)

// forwardedForHeader lists the addresses a request passed through, the client first.
const forwardedForHeader = "X-Forwarded-For"

// ClientIP returns the IP address of the browser that sent request, or "" when it cannot be read.
// It is the TCP peer's address, unless that peer is one of trustedProxies (TRUSTED_PROXIES): then
// X-Forwarded-For is read from the right, skipping trusted proxies, so a client cannot forge its
// address by sending the header itself.
func ClientIP(request *http.Request, trustedProxies []netip.Prefix) string {
	peer, ok := parseIP(request.RemoteAddr)
	if !ok {
		return ""
	}
	if !isTrustedProxy(peer, trustedProxies) {
		return peer.String()
	}
	forwarded := strings.Split(strings.Join(request.Header.Values(forwardedForHeader), ","), ",")
	client := peer
	for index := len(forwarded) - 1; index >= 0; index-- {
		address, ok := parseIP(forwarded[index])
		if !ok {
			break
		}
		client = address
		if !isTrustedProxy(address, trustedProxies) {
			break
		}
	}
	return client.String()
}

// parseIP reads an IP address with or without a port ("10.0.0.1:5000", "[::1]:80", " 10.0.0.1").
func parseIP(value string) (netip.Addr, bool) {
	value = strings.TrimSpace(value)
	if host, _, err := net.SplitHostPort(value); err == nil {
		value = host
	}
	address, err := netip.ParseAddr(value)
	if err != nil {
		return netip.Addr{}, false
	}
	return address.Unmap(), true
}

// isTrustedProxy reports whether address is in one of trustedProxies.
func isTrustedProxy(address netip.Addr, trustedProxies []netip.Prefix) bool {
	for _, prefix := range trustedProxies {
		if prefix.Contains(address) {
			return true
		}
	}
	return false
}
