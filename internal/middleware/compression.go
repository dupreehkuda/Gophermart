package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	g.Header().Del("Content-Length")
	return g.writer.Write([]byte(s))
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	g.Header().Del("content-Length")
	return g.writer.Write(data)
}

func (m middleware) CheckCompression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer func(gz *gzip.Reader) {
			err := gz.Close()
			if err != nil {
				log.Printf("Error occured while closing writer: %v", err)
			}
		}(gz)

		r.Body = io.NopCloser(gz)
		w.Header().Del("Content-Length")

		next.ServeHTTP(w, r)
	})
}

func (m middleware) WriteCompressed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer func(gz *gzip.Writer) {
			err := gz.Close()
			if err != nil {
				log.Printf("Error occured while closing writer: %v", err)
			}
		}(gz)

		w = &gzipWriter{w, gz}
		w.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(w, r)
	})
}
