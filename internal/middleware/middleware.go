package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/LexusEgorov/goMetrics/internal/decoding"
	"github.com/LexusEgorov/goMetrics/internal/dohsimpson"
	"github.com/LexusEgorov/goMetrics/internal/encoding"
)

type ResponseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
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
		if r.Method != "POST" {
			next.ServeHTTP(w, r)
			return
		}

		encodingHeader := r.Header.Get("Content-Encoding")

		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data []byte
		var decodeErr *dohsimpson.Error

		switch encodingHeader {
		case "gzip":
			data, decodeErr = decoder.DecodeGz(buf.Bytes())
		case "deflate":
			data, decodeErr = decoder.DecodeDeflate(buf.Bytes())
		case "br":
			data, decodeErr = decoder.DecodeBr(buf.Bytes())
		default:
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
	encoder := encoding.NewEncoding()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			next.ServeHTTP(w, r)
			return
		}

		rw := &ResponseWriter{ResponseWriter: w}
		encodingHeader := r.Header.Get("Accept-Encoding")

		next.ServeHTTP(rw, r)

		var data []byte
		var encodeErr *dohsimpson.Error

		switch encodingHeader {
		case "gzip":
			data, encodeErr = encoder.EncodeGz(rw.body.Bytes())
			w.Header().Set("Content-Encoding", "gzip")
		case "deflate":
			data, encodeErr = encoder.EncodeDeflate(rw.body.Bytes())
			w.Header().Set("Content-Encoding", "deflate")
		case "br":
			data, encodeErr = encoder.EncodeBr(rw.body.Bytes())
			w.Header().Set("Content-Encoding", "br")
		default:
			fmt.Printf("HEADER: %s\n", encodingHeader)
			next.ServeHTTP(w, r)
			return
		}

		if encodeErr != nil {
			w.WriteHeader(encodeErr.Code)
			return
		}

		w.Write(data)

		next.ServeHTTP(w, r)
	})
}
