package main

import (
	"fmt"

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

	if err := server.Run(serverVars); err != nil {
		panic(err)
	}
}
