package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/LexusEgorov/goMetrics/internal/config"
	"github.com/LexusEgorov/goMetrics/internal/services/collectmetric"
)

type MetricsCollector interface {
	collectMetrics()
	sendMetrics()
	Start(stopChan chan struct{})
}

func main() {
	agentVars := config.GetAgent()

	run(agentVars.Host, agentVars.ReportInterval, agentVars.PollInterval)
}

func run(host string, reportInterval, pollInterval int) {
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
