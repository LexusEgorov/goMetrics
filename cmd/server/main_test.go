package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LexusEgorov/goMetrics/internal/transport"
	"github.com/go-chi/chi/v5"
)

func setupRouter() *chi.Mux {
	transportLayer := transport.CreateTransport()
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
		{"GET", "/value/type/name", http.StatusNotFound},
		{"POST", "/update/gauge/metricName/1", http.StatusOK},
		{"POST", "/update/undefined/metricName/1", http.StatusBadRequest},
		{"POST", "/update/counter/counterMetric/1", http.StatusOK},
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
