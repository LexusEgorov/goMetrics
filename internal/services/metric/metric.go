package metric

type Metric struct {
	gauge   float64
	counter int64
}

func (m *Metric) UpdateGauge(gaugeValue float64) {
	m.gauge = gaugeValue
}

func (m *Metric) UpdateCounter(counterValue int64) {
	m.counter += counterValue
}

func CreateMetric() Metric {
	return Metric{
		gauge:   0,
		counter: 0,
	}
}
