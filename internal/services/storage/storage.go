package storage

import (
	"fmt"
	"strconv"
)

type MetricName string
type Gauge float64
type Counter int64

func (g Gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}

func (c Counter) String() string {
	return strconv.Itoa(int(c))
}

type memStorage struct {
	data map[MetricName]interface{}
}

func (m *memStorage) AddGauge(key MetricName, value Gauge) {
	m.data[key] = value
}

func (m *memStorage) AddCounter(key MetricName, value Counter) {
	if existing, ok := m.data[key]; ok {
		if counterValue, ok := existing.(Counter); ok {
			m.data[key] = counterValue + value
		} else {
			fmt.Printf("Err: value for key: %s isn't Counter", key)
		}
	} else {
		m.data[key] = value
	}
}

func (m memStorage) GetGauge(key MetricName) (Gauge, bool) {
	value, isFound := m.data[key].(Gauge)

	return value, isFound
}

func (m memStorage) GetCounter(key MetricName) (Counter, bool) {
	value, isFound := m.data[key].(Counter)

	return value, isFound
}

func (m memStorage) GetAll() map[MetricName]interface{} {
	return m.data
}

func CreateStorage() *memStorage {
	return &memStorage{
		data: make(map[MetricName]interface{}),
	}
}
