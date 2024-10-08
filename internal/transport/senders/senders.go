package senders

import "fmt"

func SendMetric(metricName, metricType, metricValue string) {
	fmt.Printf("METRIC: %s | TYPE: %s | VALUE: %s\n", metricName, metricType, metricValue)
}
