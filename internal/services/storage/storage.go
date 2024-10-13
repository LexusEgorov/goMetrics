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

type Storager interface {
	AddGauge(key MetricName, value Gauge)
	AddCounter(key MetricName, value Counter)
	GetGauge(key MetricName) (Gauge, bool)
	GetCounter(key MetricName) (Counter, bool)
	GetAll() map[MetricName]interface{}
}

type MemStorage struct {
	data map[MetricName]interface{}
}

func (m *MemStorage) AddGauge(key MetricName, value Gauge) {
	m.data[key] = value
}

func (m *MemStorage) AddCounter(key MetricName, value Counter) {
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

func (m MemStorage) GetGauge(key MetricName) (Gauge, bool) {
	value, isFound := m.data[key].(Gauge)

	return value, isFound
}

func (m MemStorage) GetCounter(key MetricName) (Counter, bool) {
	value, isFound := m.data[key].(Counter)

	return value, isFound
}

func (m MemStorage) GetAll() map[MetricName]interface{} {
	return m.data
}

func CreateStorage() Storager {
	return &MemStorage{
		data: make(map[MetricName]interface{}),
	}
}
