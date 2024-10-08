package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/LexusEgorov/goMetrics/internal/services/storage"
)

var metricStorage = storage.CreateStorage()

func UpdateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mName := r.PathValue("metricName")
	mType := r.PathValue("metricType")
	mValue := r.PathValue("metricValue")

	fmt.Printf("Name: %s\nType: %s\nValue: %s\n ============\n", mName, mType, mValue)

	if mName == "" || mValue == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch mType {
	case "gauge":
		value, err := strconv.ParseFloat(mValue, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		metricStorage.AddGauge(storage.MetricName(mName), storage.Gauge(value))
	case "counter":
		value, err := strconv.ParseInt(mValue, 0, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		metricStorage.AddCounter(storage.MetricName(mName), storage.Counter(value))
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
