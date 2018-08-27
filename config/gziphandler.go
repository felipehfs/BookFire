// Package config make the setup of the all app
package config

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

var zippers = sync.Pool{New: func() interface{} {
	return gzip.NewWriter(nil)
}}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GzipHandler compacts the response of the server
func GzipHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gz := zippers.Get().(*gzip.Writer)

		gz.Reset(w)
		defer zippers.Put(gz)
		defer gz.Close()
		next(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	}
}
