package main

import (
	"fmt"
	"net/http"

	"github.com/LexusEgorov/goMetrics/internal/config"
	"github.com/LexusEgorov/goMetrics/internal/transport"
	"github.com/go-chi/chi/v5"
)

func main() {
	serverVars := config.GetServerVars()

	if err := run(serverVars.Host); err != nil {
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
