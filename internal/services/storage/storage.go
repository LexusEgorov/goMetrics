package storage

import "github.com/LexusEgorov/goMetrics/internal/services/metric"

type MetricName string

type MemStorage struct {
	metrics map[MetricName]metric.Metric
}

func (m *MemStorage) AddMetric(metricName MetricName) metric.Metric {
	createdMetric := metric.CreateMetric()

	m.metrics[metricName] = createdMetric

	return createdMetric
}

func (m MemStorage) GetMetric(metricName MetricName) *metric.Metric {
	metrics := m.metrics

	if len(metrics) == 0 {
		createdMetric := m.AddMetric(metricName)
		return &createdMetric
	}

	foundMetric, isFound := m.metrics[metricName]

	if !isFound {
		foundMetric = m.AddMetric(metricName)
	}

	return &foundMetric
}

func (m *MemStorage) SetMetric(mName MetricName, metric metric.Metric) {
	m.metrics[mName] = metric
}

func (m *MemStorage) GetMetrics() string {
	return "123"
}

func CreateStorage() MemStorage {
	return MemStorage{
		metrics: make(map[MetricName]metric.Metric),
	}
}
