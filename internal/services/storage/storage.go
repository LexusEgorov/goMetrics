package storage

import (
	"sync"

	"github.com/LexusEgorov/goMetrics/internal/keeper"
	"github.com/LexusEgorov/goMetrics/internal/models"
)

type memStorage struct {
	mu   sync.Mutex
	data map[string]models.Metric
}

func (m *memStorage) MassSave(metrics []models.Metric) ([]models.Metric, error) {
	savedMetrics := make([]models.Metric, len(metrics))

	for i, metric := range metrics {
		switch metric.MType {
		case "gauge":
			m.AddGauge(metric.ID, *metric.Value)
			savedMetrics[i] = metric
		case "counter":
			oldValue, _ := m.GetCounter(metric.ID)

			m.AddCounter(metric.ID, int64(*metric.Delta))

			newValue := *metric.Delta + oldValue

			metric.Delta = &newValue
			savedMetrics[i] = metric
		}
	}

	return savedMetrics, nil
}

func (m *memStorage) AddGauge(key string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = models.Metric{
		ID:    key,
		MType: "gauge",
		Value: &value,
	}
}

func (m *memStorage) AddCounter(key string, value int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

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

func (m *memStorage) GetGauge(key string) (float64, bool) {
	metric, isFound := m.data[key]

	if metric.Value == nil {
		return 0, false
	}

	return *metric.Value, isFound
}

func (m *memStorage) GetCounter(key string) (int64, bool) {
	metric, isFound := m.data[key]

	if metric.Delta == nil {
		return 0, false
	}

	return *metric.Delta, isFound
}

func (m *memStorage) GetAll() map[string]models.Metric {
	return m.data
}

func (m *memStorage) Check() bool {
	return true
}

func (m *memStorage) Close() {}

func NewStorage(metrics map[string]models.Metric) keeper.Storager {
	storage := &memStorage{
		mu:   sync.Mutex{},
		data: metrics,
	}

	return storage
}
