package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/LexusEgorov/goMetrics/internal/config"
	"github.com/LexusEgorov/goMetrics/internal/runners"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func defaultIfEmpty(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}

func init() {
	fmt.Printf("Build version: %s\n", defaultIfEmpty(buildVersion))
	fmt.Printf("Build date: %s\n", defaultIfEmpty(buildDate))
	fmt.Printf("Build commit: %s\n", defaultIfEmpty(buildCommit))
}

func main() {
	server := runners.NewServer()
	serverVars := config.NewServer()

	stopChan := make(chan struct{})
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-signalChan
		close(stopChan)
	}()

	if err := server.Run(serverVars, stopChan); err != nil {
		panic(err)
	}
}
