package main

import (
	"net/http"
	"strings"
)

func compressMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		if contentEncoding := r.Header.Get("Content-Encoding"); strings.Contains(contentEncoding, "gzip") {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}

		if acceptEncoding := r.Header.Get("Accept-Encoding"); strings.Contains(acceptEncoding, "gzip") {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		// передаём управление хендлеру
		h.ServeHTTP(ow, r)
	})
}
