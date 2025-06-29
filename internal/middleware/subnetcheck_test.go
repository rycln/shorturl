package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrustedSubnetMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		trustedSubnet  string
		xRealIP        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "No trusted subnet configured",
			trustedSubnet:  "",
			xRealIP:        "192.168.1.1",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "Access forbidden (no trusted subnet configured)\n",
		},
		{
			name:           "Missing X-Real-IP header",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "X-Real-IP header required\n",
		},
		{
			name:           "Invalid trusted subnet configuration",
			trustedSubnet:  "invalid_subnet",
			xRealIP:        "192.168.1.1",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Invalid trusted subnet configuration\n",
		},
		{
			name:           "IP not in trusted subnet",
			trustedSubnet:  "10.0.0.0/8",
			xRealIP:        "192.168.1.1",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "Access forbidden (IP not in trusted subnet)\n",
		},
		{
			name:           "Valid IP in trusted subnet",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "192.168.1.100",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "Invalid IP format",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "invalid_ip",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "Access forbidden (IP not in trusted subnet)\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			middleware := TrustedSubnet(tt.trustedSubnet)
			handlerToTest := middleware(nextHandler)

			req := httptest.NewRequest("GET", "/", nil)
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}

			rec := httptest.NewRecorder()

			handlerToTest.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if rec.Body.String() != tt.expectedBody {
				t.Errorf("expected body '%s', got '%s'", tt.expectedBody, rec.Body.String())
			}
		})
	}
}
