package collectmetric

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"runtime"
	"time"

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

type MetricsCollector interface {
	collectMetrics()
	sendMetrics()
	Start(stopChan chan struct{})
}

type agentIntervals struct {
	collect int
	send    int
}

type MetricAgent struct {
	storage   storage.MemStorage
	pollCount storage.Counter
	intervals agentIntervals
}

func (agent *MetricAgent) collectMetrics() {
	for {
		fmt.Println("Collect started")
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		for _, metricName := range gaugeMetrics {
			value := reflect.ValueOf(memStats).FieldByName(metricName)
			agent.pollCount++

			if value.IsValid() && value.CanInterface() {
				switch v := value.Interface().(type) {
				case float64:
					agent.storage.AddGauge(storage.MetricName(metricName), storage.Gauge(v))
				case uint64:
					agent.storage.AddGauge(storage.MetricName(metricName), storage.Gauge(float64(v)))
				case uint32:
					agent.storage.AddGauge(storage.MetricName(metricName), storage.Gauge(float64(v)))
				case uint16:
					agent.storage.AddGauge(storage.MetricName(metricName), storage.Gauge(float64(v)))
				case uint8:
					agent.storage.AddGauge(storage.MetricName(metricName), storage.Gauge(float64(v)))
				default:
					fmt.Printf("Unable to convert metric %s (%s) to a float64\n", metricName, v)
				}
			} else {
				fmt.Printf("Metric %s is not valid or accessible\n", metricName)
				continue
			}
		}

		agent.storage.AddCounter("PollCount", agent.pollCount)
		randomValue := rand.Float64()
		agent.storage.AddGauge("RandomValue", storage.Gauge(randomValue))

		fmt.Println("Collect finished")
		time.Sleep(time.Duration(agent.intervals.collect) * time.Second)
	}
}

func (agent MetricAgent) sendMetrics() {
	for {
		time.Sleep(time.Duration(agent.intervals.send) * time.Second)
		fmt.Println("Sending started")
		for k, v := range agent.storage.GetAll() {
			switch metric := v.(type) {
			case storage.Gauge:
				senders.SendMetric(string(k), "gauge", metric.ToString())
			case storage.Counter:
				senders.SendMetric(string(k), "counter", metric.ToString())
			default:
				fmt.Printf("Unknown metric's type: %T\n", v)
			}
		}

		fmt.Println("Sending finished")
	}
}

func (agent MetricAgent) Start(stopChan chan struct{}) {
	fmt.Println("Agent started")
	go agent.collectMetrics()
	go agent.sendMetrics()

	<-stopChan
	fmt.Println("Agent finished")
}

func CreateAgent() MetricsCollector {
	return &MetricAgent{
		storage:   storage.CreateStorage(),
		pollCount: 0,
		intervals: agentIntervals{
			collect: 2,
			send:    10,
		},
	}
}
