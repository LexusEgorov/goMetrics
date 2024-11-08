package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ResponseWriter struct {
	http.ResponseWriter
}

func WithLogging(s *zap.SugaredLogger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			uri := r.RequestURI
			method := r.Method

			h.ServeHTTP(w, r)

			timeDiff := time.Since(start)

			s.Infoln(
				"uri", uri,
				"method", method,
				"duration", timeDiff,
			)
		})
	}
}

func WithDecoding(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodingHeader := r.Header.Get("Content-Encoding")

		switch encodingHeader {
		case "gzip":

		case "compress":

		case "deflate":

		case "br":

		default:
		}

		next.ServeHTTP(w, r)
	})
}

func WithEncoding(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &ResponseWriter{ResponseWriter: w}
		encodingHeader := r.Header.Get("Accept-Encoding")

		next.ServeHTTP(rw, r)

		switch encodingHeader {
		case "gzip":

		case "compress":

		case "deflate":

		case "br":

		default:
		}

		next.ServeHTTP(w, r)
	})
}
