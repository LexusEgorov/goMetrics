package main

import (
	"net/http"

	"github.com/LexusEgorov/goMetrics/internal/transport"
)

func main() {
	transportLayer := transport.CreateTransport()
	mux := http.NewServeMux()

	mux.HandleFunc(`/update/{metricType}/{metricName}/{metricValue}`, transportLayer.UpdateMetric)

	http.ListenAndServe(":8080", mux)
}
