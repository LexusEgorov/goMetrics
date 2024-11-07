package saver

import (
	"net/http"
	"strconv"

	"github.com/LexusEgorov/goMetrics/internal/dohSimpson"
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

func (s saver) SaveOld(mName, mType, mValue string) *dohSimpson.Error {
	if mName == "" || mValue == "" {
		return dohSimpson.NewDoh(http.StatusNotFound, "metric not found")
	}

	switch mType {
	case "gauge":
		value, err := strconv.ParseFloat(mValue, 64)

		if err != nil {
			return dohSimpson.NewDoh(http.StatusBadRequest, err.Error())
		}

		s.storage.AddGauge(mName, float64(value))
	case "counter":
		value, err := strconv.ParseInt(mValue, 0, 64)

		if err != nil {
			return dohSimpson.NewDoh(http.StatusBadRequest, err.Error())
		}

		s.storage.AddCounter(mName, int64(value))
	default:
		return dohSimpson.NewDoh(http.StatusBadRequest, "unknown metric Type")
	}

	return nil
}

func (s saver) Save(m models.Metric) (*models.Metric, *dohSimpson.Error) {
	mName := m.ID
	mType := m.MType
	mValue := m.Value
	mDelta := m.Delta

	if mName == "" {
		return nil, dohSimpson.NewDoh(http.StatusNotFound, "empty metric ID")
	}

	if mValue == nil && mDelta == nil {
		return nil, dohSimpson.NewDoh(http.StatusBadRequest, "empty metric value or delta")
	}

	switch mType {
	case "gauge":
		s.storage.AddGauge(mName, float64(*mValue))
	case "counter":
		s.storage.AddCounter(mName, int64(*mDelta))
	default:
		return nil, dohSimpson.NewDoh(http.StatusBadRequest, "unknown metric Type")
	}

	return &m, nil
}

func NewSaver(storage Storager) transport.Saver {
	return saver{
		storage: storage,
	}
}
