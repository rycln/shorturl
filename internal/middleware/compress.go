package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// compressWriter wraps http.ResponseWriter to provide gzip compression.
//
// Automatically sets "Content-Encoding: gzip" header for successful responses.
// Implements http.ResponseWriter, http.Flusher and io.Closer.
type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// newCompressWriter creates new compression writer.
func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header returns the header map from the original response writer.
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write compresses and writes data to the underlying connection.
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader sets status code and conditionally enables compression.
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close flushes compressed data and releases resources.
//
// Must be called to ensure all data is properly written.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader wraps io.ReadCloser to provide gzip decompression.
//
// Automatically handles both compressed and uncompressed request bodies.
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// newCompressReader creates new decompression reader.
func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read decompresses data from the underlying request body.
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close releases all compression-related resources.
//
// Ensures both the gzip reader and original body are properly closed.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// Compress is middleware that handles GZIP compression for HTTP responses
// and decompression for requests.
//
// Features:
// - Response compression for clients that support gzip (Accept-Encoding)
// - Request decompression for gzipped request bodies (Content-Encoding)
func Compress(h http.Handler) http.Handler {
	comp := func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	}

	return http.HandlerFunc(comp)
}
