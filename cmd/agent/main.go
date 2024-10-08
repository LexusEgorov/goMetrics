package main

import "github.com/LexusEgorov/goMetrics/internal/services/collectmetric"

func main() {
	collectmetric.CollectMetrics()
	collectmetric.SendMetrics()
}
