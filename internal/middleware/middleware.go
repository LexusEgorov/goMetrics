package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/LexusEgorov/goMetrics/internal/decoding"
	"github.com/LexusEgorov/goMetrics/internal/dohsimpson"
)

type responseWriter struct {
	http.ResponseWriter
	Body *bytes.Buffer
}

func (rw responseWriter) Write(b []byte) (int, error) {
	return rw.Body.Write(b)
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
	decoder := decoding.NewDecoding()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		encodingHeader := r.Header.Get("Content-Encoding")

		body, err := io.ReadAll(r.Body)

		if len(body) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data []byte
		var decodeErr *dohsimpson.Error

		switch encodingHeader {
		case "gzip":
			data, decodeErr = decoder.DecodeGz(body)
		case "deflate":
			data, decodeErr = decoder.DecodeDeflate(body)
		case "br":
			data, decodeErr = decoder.DecodeBr(body)
		default:
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			next.ServeHTTP(w, r)
			return
		}

		if decodeErr != nil {
			w.WriteHeader(decodeErr.Code)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(data))

		next.ServeHTTP(w, r)
	})
}

func WithEncoding(next http.Handler) http.Handler {
	// encoder := encoding.NewEncoding()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{
			ResponseWriter: w,
			Body:           new(bytes.Buffer),
		}

		encodingHeader := r.Header.Get("Accept-Encoding")

		next.ServeHTTP(rw, r)

		data := rw.Body.Bytes()
		var encodeErr *dohsimpson.Error

		switch encodingHeader {
		case "gzip":
			// data, encodeErr = encoder.EncodeGz(rw.Body.Bytes())
			w.Header().Set("Content-Encoding", "gzip")
		case "deflate":
			// data, encodeErr = encoder.EncodeDeflate(rw.Body.Bytes())
			w.Header().Set("Content-Encoding", "deflate")
		case "br":
			// data, encodeErr = encoder.EncodeBr(rw.Body.Bytes())
			w.Header().Set("Content-Encoding", "br")
		default:
		}

		if encodeErr != nil {
			w.WriteHeader(encodeErr.Code)
			return
		}

		// w.Header().Set("Content-Length", fmt.Sprint(len(data)))
		w.Write(data)
	})
}
