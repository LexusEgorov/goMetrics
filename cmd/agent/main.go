package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/LexusEgorov/goMetrics/internal/config"
	"github.com/LexusEgorov/goMetrics/internal/services/collectmetric"
)

func main() {
	agentVars := config.GetAgentVars()

	run(agentVars.Host, agentVars.ReportInterval, agentVars.PollInterval)
}

func run(host string, reportInterval, pollInterval int) {
	agent := collectmetric.CreateAgent(host, reportInterval, pollInterval)

	stopChan := make(chan struct{})
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		close(stopChan)
	}()

	agent.Start(stopChan)
}
