package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LexusEgorov/goMetrics/internal/transport"
	"github.com/go-chi/chi/v5"
)

func setupRouter() *chi.Mux {
	transportLayer := transport.NewTransport()
	r := chi.NewRouter()
	r.Get("/", transportLayer.GetMetrics)
	r.Get("/value/{metricType}/{metricName}", transportLayer.GetMetric)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", transportLayer.UpdateMetric)

	return r
}

func Test_main(t *testing.T) {
	r := setupRouter()

	tests := []struct {
		method   string
		url      string
		expected int
	}{
		{"GET", "/", http.StatusOK},
	}

	for _, test := range tests {
		req := httptest.NewRequest(test.method, test.url, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != test.expected {
			t.Errorf("Expected status %d, got %d for %s %s", test.expected, rr.Code, test.method, test.url)
		}
	}
}
