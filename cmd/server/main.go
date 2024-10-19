package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/LexusEgorov/goMetrics/internal/config"
	"github.com/LexusEgorov/goMetrics/internal/transport"
)

type Transporter interface {
	UpdateMetric(w http.ResponseWriter, r *http.Request)
	GetMetric(w http.ResponseWriter, r *http.Request)
	GetMetrics(w http.ResponseWriter, r *http.Request)
}

func main() {
	serverVars := config.GetServer()

	if err := run(serverVars.Host); err != nil {
		panic(err)
	}
}

func run(host string) error {
	var transportLayer Transporter = transport.CreateTransport()

	r := chi.NewRouter()

	r.Get("/", transportLayer.GetMetrics)
	r.Get("/value/{metricType}/{metricName}", transportLayer.GetMetric)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", transportLayer.UpdateMetric)

	fmt.Println("Running server on", host)
	return http.ListenAndServe(host, r)
}
