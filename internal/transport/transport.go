package transport

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/LexusEgorov/goMetrics/internal/services/storage"
)

type Transporter interface {
	UpdateMetric(w http.ResponseWriter, r *http.Request)
	SendMetric(metricName, metricType, metricValue string)
}

type TransportLayer struct {
	storage storage.Storager
}

func (t TransportLayer) UpdateMetric(w http.ResponseWriter, r *http.Request) {
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

		t.storage.AddGauge(storage.MetricName(mName), storage.Gauge(value))
	case "counter":
		value, err := strconv.ParseInt(mValue, 0, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		t.storage.AddCounter(storage.MetricName(mName), storage.Counter(value))
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t TransportLayer) SendMetric(metricName, metricType, metricValue string) {
	url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%s", metricType, metricName, metricValue)
	response, err := http.Post(url, "text/plain", nil)

	if err != nil {
		fmt.Printf("ERR: %s\n", err)
	} else {
		fmt.Printf("RESPONSE: %s\n", response.Status)
	}

	defer response.Body.Close()
}

func CreateTransport() Transporter {
	return TransportLayer{
		storage: storage.CreateStorage(),
	}
}
