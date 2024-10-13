package main

import (
	"fmt"
	"net/http"

	"github.com/LexusEgorov/goMetrics/internal/services/flags"
	"github.com/LexusEgorov/goMetrics/internal/transport"
	"github.com/go-chi/chi/v5"
)

func main() {
	ServerFlags := flags.GetServerFlags()

	if err := run(ServerFlags.Host); err != nil {
		panic(err)
	}
}

func run(host string) error {
	transportLayer := transport.CreateTransport()

	r := chi.NewRouter()

	r.Get("/", transportLayer.GetMetrics)
	r.Get("/value/{metricType}/{metricName}", transportLayer.GetMetric)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", transportLayer.UpdateMetric)

	fmt.Println("Running server on", host)
	return http.ListenAndServe(host, r)
}
