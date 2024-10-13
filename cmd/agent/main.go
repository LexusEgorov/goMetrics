package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/LexusEgorov/goMetrics/internal/services/collectmetric"
	"github.com/LexusEgorov/goMetrics/internal/services/flags"
)

func main() {
	agentFlags := flags.GetAgentFlags()

	run(agentFlags.Host, agentFlags.ReportInterval, agentFlags.PollInterval)
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
