package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type metricName string

type Metric struct {
	gauge   float64
	counter int64
}

func (m *Metric) updateGauge(gaugeValue float64) {
	m.gauge = gaugeValue
}

func (m *Metric) updateCounter(counterValue int64) {
	m.counter += counterValue
}

type MemStorage struct {
	metrics map[metricName]Metric
}

func (m *MemStorage) addMetric(metricName metricName) Metric {
	createdMetric := Metric{
		gauge:   0,
		counter: 0,
	}

	m.metrics[metricName] = createdMetric

	return createdMetric
}

func (m MemStorage) getMetric(metricName metricName) *Metric {
	metrics := m.metrics

	if len(metrics) == 0 {
		createdMetric := m.addMetric(metricName)
		return &createdMetric
	}

	foundMetric, isFound := m.metrics[metricName]

	if !isFound {
		foundMetric = m.addMetric(metricName)
	}

	return &foundMetric
}

func (m *MemStorage) setMetric(mName metricName, metric Metric) {
	m.metrics[mName] = metric
}

var storage = MemStorage{
	metrics: make(map[metricName]Metric),
}

func updateMetric(w http.ResponseWriter, r *http.Request) {
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

	currentMetric := storage.getMetric(metricName(mName))

	switch mType {
	case "gauge":
		value, err := strconv.ParseFloat(mValue, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		currentMetric.updateGauge(value)
		storage.setMetric(metricName(mName), *currentMetric)

	case "counter":
		value, err := strconv.ParseInt(mValue, 0, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Println(value)
		currentMetric.updateCounter(value)
		storage.setMetric(metricName(mName), *currentMetric)

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(storage.metrics)
	w.WriteHeader(http.StatusOK)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(`/update/{metricType}/{metricName}/{metricValue}`, updateMetric)

	http.ListenAndServe(":8080", mux)
}
