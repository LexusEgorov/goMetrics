package middleware

import (
	"net/http"
	"time"

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
