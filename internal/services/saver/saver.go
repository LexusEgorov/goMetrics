package saver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/LexusEgorov/goMetrics/internal/dohsimpson"
	"github.com/LexusEgorov/goMetrics/internal/models"
	"github.com/LexusEgorov/goMetrics/internal/transport"
)

type Storager interface {
	AddGauge(key string, value float64)
	AddCounter(key string, value int64)
}

type saver struct {
	storage Storager
}

func (s saver) SaveOld(mName, mType, mValue string) *dohsimpson.Error {
	if mName == "" || mValue == "" {
		return dohsimpson.NewDoh(http.StatusNotFound, "metric not found (empty data) (saver)")
	}

	switch mType {
	case "gauge":
		value, err := strconv.ParseFloat(mValue, 64)

		if err != nil {
			return dohsimpson.NewDoh(http.StatusBadRequest, err.Error())
		}

		s.storage.AddGauge(mName, float64(value))
	case "counter":
		value, err := strconv.ParseInt(mValue, 0, 64)

		if err != nil {
			return dohsimpson.NewDoh(http.StatusBadRequest, err.Error())
		}

		s.storage.AddCounter(mName, int64(value))
	default:
		return dohsimpson.NewDoh(http.StatusBadRequest, fmt.Sprintf("unknown metric Type (%s) (saver)", mType))
	}

	return nil
}

func (s saver) Save(m models.Metric) (*models.Metric, *dohsimpson.Error) {
	mName := m.ID
	mType := m.MType
	mValue := m.Value
	mDelta := m.Delta

	if mName == "" {
		return nil, dohsimpson.NewDoh(http.StatusNotFound, "saver: empty metric ID")
	}

	if mValue == nil && mDelta == nil {
		return nil, dohsimpson.NewDoh(http.StatusBadRequest, "saver: empty metric value and delta")
	}

	switch mType {
	case "gauge":
		s.storage.AddGauge(mName, float64(*mValue))
	case "counter":
		s.storage.AddCounter(mName, int64(*mDelta))
	default:
		return nil, dohsimpson.NewDoh(http.StatusBadRequest, fmt.Sprintf("saver: unknown metric Type (%s)", mType))
	}

	return &m, nil
}

func NewSaver(storage Storager) transport.Saver {
	return saver{
		storage: storage,
	}
}
