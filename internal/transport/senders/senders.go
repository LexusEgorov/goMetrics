package senders

import (
	"fmt"
	"net/http"
)

func SendMetric(metricName, metricType, metricValue string) {
	url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", metricType, metricName, metricValue)
	response, err := http.Post(url, "text/plain", nil)

	if err != nil {
		fmt.Printf("ERR: %s\n", err)
	} else {
		fmt.Printf("RESPONSE: %s\n", response.Status)
	}

	// fmt.Printf("METRIC: %s | TYPE: %s | VALUE: %s\n", metricName, metricType, metricValue)
}
