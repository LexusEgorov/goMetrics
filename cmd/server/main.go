package main

import (
	"github.com/LexusEgorov/goMetrics/internal/config"
	"github.com/LexusEgorov/goMetrics/internal/runners"
)

func main() {
	server := runners.NewServer()
	serverVars := config.NewServer()

	if err := server.Run(serverVars); err != nil {
		panic(err)
	}
}
