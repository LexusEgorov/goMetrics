package main

import (
	"github.com/LexusEgorov/goMetrics/internal/config"
	"github.com/LexusEgorov/goMetrics/internal/runners"
)

type MetricsCollector interface {
	collectMetrics()
	sendMetrics()
	Start(stopChan chan struct{})
}

func main() {
	agent := runners.NewAgent()
	agentVars := config.NewAgent()

	agent.Run(agentVars)
}
