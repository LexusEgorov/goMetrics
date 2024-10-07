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

	fmt.Printf("Name: %s\nType: %s\nValue: %s\n", mName, mType, mValue)

	if mName == "" || mValue == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	currentMetric := metricStorage.GetMetric(storage.MetricName(mName))

	switch mType {
	case "gauge":
		value, err := strconv.ParseFloat(mValue, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		currentMetric.UpdateGauge(value)
		metricStorage.SetMetric(storage.MetricName(mName), *currentMetric)

	case "counter":
		value, err := strconv.ParseInt(mValue, 0, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Println(value)
		currentMetric.UpdateCounter(value)
		metricStorage.SetMetric(storage.MetricName(mName), *currentMetric)

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(metricStorage.GetMetrics())
	w.WriteHeader(http.StatusOK)
}
