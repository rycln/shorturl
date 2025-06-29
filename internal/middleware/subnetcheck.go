package middleware

import (
	"net"
	"net/http"
)

// TrustedSubnetMW checks if the request's X-Real-IP is in a trusted subnet.
// If trustedSubnet is empty, all requests are denied.
func TrustedSubnet(trustedSubnet string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if trustedSubnet == "" {
				http.Error(w, "Access forbidden (no trusted subnet configured)", http.StatusForbidden)
				return
			}

			realIP := r.Header.Get("X-Real-IP")
			if realIP == "" {
				http.Error(w, "X-Real-IP header required", http.StatusForbidden)
				return
			}

			_, subnet, err := net.ParseCIDR(trustedSubnet)
			if err != nil {
				http.Error(w, "Invalid trusted subnet configuration", http.StatusInternalServerError)
				return
			}

			clientIP := net.ParseIP(realIP)
			if clientIP == nil || !subnet.Contains(clientIP) {
				http.Error(w, "Access forbidden (IP not in trusted subnet)", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
