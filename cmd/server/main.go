package main

import (
	"net/http"

	"github.com/LexusEgorov/goMetrics/internal/transport/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(`/update/{metricType}/{metricName}/{metricValue}`, handlers.UpdateMetric)

	http.ListenAndServe(":8080", mux)
}
