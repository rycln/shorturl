package middleware

import (
	"net/http"
	"time"

	"github.com/rycln/shorturl/internal/logger"
	"go.uber.org/zap"
)

// responseData holds metrics about the HTTP response.
type responseData struct {
	status int
	size   int
}

// loggingResponseWriter wraps http.ResponseWriter to track response metrics.
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// Write intercepts and measures response body writes.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader intercepts and records the status code.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// Logger is middleware that logs HTTP request/response details.
//
// Logs the following information for each request:
// - Request URL and method
// - Response status code
// - Response body size (bytes)
// - Request processing duration
//
// Usage:
//
//	r := chi.NewRouter()
//	r.Use(middleware.Logger)
//
// Example log output (using zap):
//
//	INFO    Req/Res Log  {"url": "/api/data", "method": "GET", "status": 200,
//	        "duration": "12.5ms", "size": 142}
func Logger(h http.Handler) http.Handler {
	log := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Log.Info("Req/Res Log",
			zap.String("url", r.RequestURI),
			zap.String("method", r.Method),
			zap.Int("status", responseData.status),
			zap.Duration("duration", duration),
			zap.Int("size", responseData.size),
		)
	}

	return http.HandlerFunc(log)
}
