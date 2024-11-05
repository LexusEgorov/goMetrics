package storage

type Metric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type memStorage struct {
	data map[string]Metric
}

func (m *memStorage) AddGauge(key string, value float64) {
	m.data[key] = Metric{
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
		m.data[key] = Metric{
			ID:    key,
			MType: "counter",
			Delta: &value,
		}
	}
}

func (m memStorage) GetGauge(key string) (float64, bool) {
	metric, isFound := m.data[key]

	return *metric.Value, isFound
}

func (m memStorage) GetCounter(key string) (int64, bool) {
	metric, isFound := m.data[key]

	return *metric.Delta, isFound
}

func (m memStorage) GetAll() map[string]Metric {
	return m.data
}

func CreateStorage() *memStorage {
	return &memStorage{
		data: make(map[string]Metric),
	}
}
