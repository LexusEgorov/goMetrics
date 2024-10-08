package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/LexusEgorov/goMetrics/internal/services/collectmetric"
)

func main() {
	agent := collectmetric.CreateAgent()

	stopChan := make(chan struct{})
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		close(stopChan)
	}()

	agent.Start(stopChan)
}
