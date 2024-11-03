package runners

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"

	"github.com/LexusEgorov/goMetrics/internal/services/collectmetric"
	"github.com/LexusEgorov/goMetrics/internal/transport"
)

type Transporter interface {
	UpdateMetric(w http.ResponseWriter, r *http.Request)
	GetMetric(w http.ResponseWriter, r *http.Request)
	GetMetrics(w http.ResponseWriter, r *http.Request)
}

type serverRunner struct{}

func (s serverRunner) Run(host string) error {
	var transportLayer Transporter = transport.CreateTransport()

	r := chi.NewRouter()

	r.Get("/", transportLayer.GetMetrics)
	r.Get("/value/{metricType}/{metricName}", transportLayer.GetMetric)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", transportLayer.UpdateMetric)

	fmt.Println("Running server on", host)
	return http.ListenAndServe(host, r)
}

func NewServer() *serverRunner {
	return &serverRunner{}
}

type agentRunner struct{}

func (a agentRunner) Run(host string, reportInterval, pollInterval int) {
	var agent = collectmetric.CreateAgent(host, reportInterval, pollInterval)

	stopChan := make(chan struct{})
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		close(stopChan)
	}()

	agent.Start(stopChan)
}

func NewAgent() *agentRunner {
	return &agentRunner{}
}
