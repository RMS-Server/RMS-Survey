package request

import (
	"net"
	"net/http"
	"strings"
)

// TrustedProxies defines IP ranges that are trusted to send X-Forwarded-For headers.
// In production, this should be configured with your reverse proxy IPs.
// Default: trust local addresses (for development behind nginx/docker)
var trustedProxies = []string{
	"127.0.0.0/8",
	"::1/128",
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
}

// isTrustedProxy checks if the remote address is in the trusted proxy list.
func isTrustedProxy(remoteAddr string) bool {
	// Extract IP from remote address (handle host:port format)
	ipStr := remoteAddr
	if idx := strings.LastIndex(remoteAddr, ":"); idx != -1 {
		ipStr = remoteAddr[:idx]
	}
	// Remove brackets from IPv6 addresses
	ipStr = strings.TrimPrefix(strings.TrimSuffix(ipStr, "]"), "[")

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	for _, cidr := range trustedProxies {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

// GetClientIP extracts the client IP address from an HTTP request.
// It checks X-Forwarded-For and X-Real-IP headers ONLY if the request
// comes from a trusted proxy. Otherwise, it uses RemoteAddr directly.
// This prevents IP spoofing attacks where clients set fake headers.
func GetClientIP(r *http.Request) string {
	remoteAddr := r.RemoteAddr

	// Only trust forwarding headers if request comes from a trusted proxy
	if isTrustedProxy(remoteAddr) {
		// Check X-Forwarded-For header (may contain multiple IPs, first is client)
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			ips := strings.Split(xff, ",")
			for _, ip := range ips {
				ip = strings.TrimSpace(ip)
				if ip != "" && net.ParseIP(ip) != nil {
					return ip
				}
			}
		}

		// Check X-Real-IP header
		if xri := r.Header.Get("X-Real-IP"); xri != "" {
			ip := strings.TrimSpace(xri)
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	// Fallback to RemoteAddr
	if idx := strings.LastIndex(remoteAddr, ":"); idx != -1 {
		ip := remoteAddr[:idx]
		// Remove brackets from IPv6 addresses
		ip = strings.TrimPrefix(strings.TrimSuffix(ip, "]"), "[")
		if net.ParseIP(ip) != nil {
			return ip
		}
	}

	// Try to parse the whole RemoteAddr as IP
	ip := remoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	ip = strings.TrimPrefix(strings.TrimSuffix(ip, "]"), "[")
	if net.ParseIP(ip) != nil {
		return ip
	}

	return ""
}
