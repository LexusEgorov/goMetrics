package middleware

import (
	"net/http"
	"time"

	"github.com/LexusEgorov/goMetrics/internal/encoding"
	"go.uber.org/zap"
)

func WithLogging(h http.Handler, s *zap.SugaredLogger) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
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
	}

	return http.HandlerFunc(logFn)
}

func WithDecoding(h http.Handler) http.HandlerFunc {
	encoding := encoding.NewEncoding()

	decodeFn := func(w http.ResponseWriter, r *http.Request) {
		encodingHeader := r.Header.Get("Content-Encoding")

		switch encodingHeader {
		case "gzip":

		case "compress":

		case "deflate":

		case "br":

		default:
		}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(decodeFn)
}

func WithEncoding(h http.Handler) http.HandlerFunc {
	encodeFn := func(w http.ResponseWriter, r *http.Request) {
		encodingHeader := r.Header.Get("Accept-Encoding")

		switch encodingHeader {
		case "gzip":

		case "compress":

		case "deflate":

		case "br":

		default:
		}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(encodeFn)
}
