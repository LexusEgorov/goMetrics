package collectmetric

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"runtime"

	"github.com/LexusEgorov/goMetrics/internal/services/storage"
	"github.com/LexusEgorov/goMetrics/internal/transport/senders"
)

var gaugeMetrics = [...]string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

var agentStorage = storage.CreateStorage()
var pollCount storage.Counter = 0

func CollectMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	for _, metricName := range gaugeMetrics {
		value := reflect.ValueOf(memStats).FieldByName(metricName)
		pollCount++

		if value.IsValid() && value.CanInterface() {
			switch v := value.Interface().(type) {
			case float64:
				agentStorage.AddGauge(storage.MetricName(metricName), storage.Gauge(v))
			default:
				fmt.Printf("Unable to convert metric %s (%s) to a float64\n", metricName, v)
			}
		} else {
			fmt.Printf("Metric %s is not valid or accessible\n", metricName)
			continue
		}

	}

	agentStorage.AddCounter("PollCount", pollCount)
	randomValue := rand.Float64()
	agentStorage.AddGauge("RandomValue", storage.Gauge(randomValue))
}

func SendMetrics() {
	for k, v := range agentStorage.GetAll() {
		switch metric := v.(type) {
		case storage.Gauge:
			senders.SendMetric(string(k), "gauge", metric.ToString())
		case storage.Counter:
			senders.SendMetric(string(k), "counter", metric.ToString())
		default:
			fmt.Printf("Unknown metric's type: %T\n", v)
		}
	}
}
