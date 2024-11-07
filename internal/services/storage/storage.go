package storage

import (
	"github.com/LexusEgorov/goMetrics/internal/models"
	"github.com/LexusEgorov/goMetrics/internal/services/reader"
	"github.com/LexusEgorov/goMetrics/internal/services/saver"
)

type memStorage struct {
	data map[string]models.Metric
}

func (m *memStorage) AddGauge(key string, value float64) {
	m.data[key] = models.Metric{
		ID:    key,
		MType: "gauge",
		Value: &value,
	}
}

func (m *memStorage) AddCounter(key string, value int64) {
	if metric, ok := m.data[key]; ok {
		delta := value + *m.data[key].Delta
		metric.Delta = &delta
		m.data[key] = metric
	} else {
		m.data[key] = models.Metric{
			ID:    key,
			MType: "counter",
			Delta: &value,
		}
	}
}

func (m memStorage) GetGauge(key string) (float64, bool) {
	metric, isFound := m.data[key]

	if metric.Value == nil {
		return 0, false
	}

	return *metric.Value, isFound
}

func (m memStorage) GetCounter(key string) (int64, bool) {
	metric, isFound := m.data[key]

	if metric.Delta == nil {
		return 0, false
	}

	return *metric.Delta, isFound
}

func (m memStorage) GetAll() map[string]models.Metric {
	return m.data
}

func NewStorage() (saver.Storager, reader.Storager) {
	storage := &memStorage{
		data: make(map[string]models.Metric),
	}

	return storage, storage
}
