package compression

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func WithUnpackingGZIP(h http.HandlerFunc) http.HandlerFunc {
	compFn := func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.Header.Get(`Content-Encoding`), "gzip") {
			unzip(&w, r)
		}

		h.ServeHTTP(w, r)
	}
	return compFn
}

func unzip(w *http.ResponseWriter, r *http.Request) {
	var reader io.Reader

	gz, err := gzip.NewReader(r.Body)
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return
	}
	reader = gz
	body, err := io.ReadAll(reader)
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		return
	}

	err = gz.Close()
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return
	}

	r.Body = io.NopCloser(strings.NewReader(string(body)))
	r.ContentLength = int64(len(string(body)))
}