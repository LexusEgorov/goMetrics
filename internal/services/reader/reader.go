package reader

import (
	"fmt"
	"net/http"

	"github.com/LexusEgorov/goMetrics/internal/dohsimpson"
	"github.com/LexusEgorov/goMetrics/internal/models"
	"github.com/LexusEgorov/goMetrics/internal/transport"
)

type Storager interface {
	GetGauge(key string) (float64, bool)
	GetCounter(key string) (int64, bool)
	GetAll() map[string]models.Metric
}

type reader struct {
	storage Storager
}

func (r reader) Read(key, mType string) (*models.Metric, *dohsimpson.Error) {
	currentMetric := models.Metric{
		ID:    key,
		MType: mType,
	}

	var isFound = false

	switch mType {
	case "gauge":
		gaugeValue, found := r.storage.GetGauge(key)

		currentMetric.Value = &gaugeValue
		isFound = found
	case "counter":
		counterValue, found := r.storage.GetCounter(key)

		currentMetric.Delta = &counterValue
		isFound = found
	default:
		return nil, dohsimpson.NewDoh(http.StatusNotFound, fmt.Sprintf("wrong mType (%s)", mType))
	}

	if !isFound {
		return nil, dohsimpson.NewDoh(http.StatusNotFound, fmt.Sprintf("metric not found: %s (%s)", key, mType))
	}

	return &currentMetric, nil
}

func (r reader) ReadAll() map[string]models.Metric {
	return r.storage.GetAll()
}

func NewReader(storage Storager) transport.Reader {
	return reader{
		storage: storage,
	}
}
