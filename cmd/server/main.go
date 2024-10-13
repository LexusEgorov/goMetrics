package main

import (
	"net/http"

	"github.com/LexusEgorov/goMetrics/internal/transport"
	"github.com/go-chi/chi/v5"
)

func main() {
	transportLayer := transport.CreateTransport()

	r := chi.NewRouter()

	r.Get("/", transportLayer.GetMetrics)
	r.Get("/value/{metricType}/{metricName}", transportLayer.GetMetric)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", transportLayer.UpdateMetric)

	http.ListenAndServe(":8080", r)
}
