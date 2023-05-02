package compression

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

func WithGzip(next http.HandlerFunc) http.HandlerFunc {
	compFn := func(w http.ResponseWriter, r *http.Request) {
		WithUnpackingGZIP(WithPackingGZIP(next))
	}
	return compFn
}

func WithUnpackingGZIP(next http.Handler) http.HandlerFunc {
	compFn := func(w http.ResponseWriter, r *http.Request) {
		var reader io.Reader

		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Printf("cannot create gzip reader: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			reader = gz
			defer gz.Close()
		} else {
			reader = r.Body
		}
		r.Body = io.NopCloser(reader)
		next.ServeHTTP(w, r)
	}
	return compFn
}

func WithPackingGZIP(next http.Handler) http.HandlerFunc {
	compFn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			log.Printf("cannot create newWriterLevel: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	}
	return compFn
}
