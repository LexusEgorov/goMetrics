package collectmetric

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"runtime"
	"time"

	"github.com/LexusEgorov/goMetrics/internal/services/storage"
	"github.com/LexusEgorov/goMetrics/internal/transport"
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

type Transporter interface {
	SendMetric(host, metricName, metricType, metricValue string)
}

type agentIntervals struct {
	collect int
	send    int
}

type Storager interface {
	AddGauge(key string, value float64)
	AddCounter(key string, value int64)
	GetGauge(key string) (float64, bool)
	GetCounter(key string) (int64, bool)
	GetAll() map[string]storage.Metric
}

type metricAgent struct {
	storage   Storager
	pollCount int64
	host      string
	intervals agentIntervals
}

func (agent *metricAgent) collectMetrics() {
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
					agent.storage.AddGauge(metricName, v)
				case uint64:
					agent.storage.AddGauge(metricName, float64(v))
				case uint32:
					agent.storage.AddGauge(metricName, float64(v))
				case uint16:
					agent.storage.AddGauge(metricName, float64(v))
				case uint8:
					agent.storage.AddGauge(metricName, float64(v))
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
		agent.storage.AddGauge("RandomValue", randomValue)

		fmt.Println("Collect finished")
		time.Sleep(time.Duration(agent.intervals.collect) * time.Second)
	}
}

func (agent metricAgent) sendMetrics() {
	var transportLayer Transporter = transport.NewTransport()

	for {
		time.Sleep(time.Duration(agent.intervals.send) * time.Second)
		fmt.Println("Sending started")
		for k, metric := range agent.storage.GetAll() {
			switch metric.MType {
			case "gauge":
				transportLayer.SendMetric(agent.host, string(k), metric.MType, fmt.Sprint(metric.Value))
			case "counter":
				transportLayer.SendMetric(agent.host, string(k), metric.MType, fmt.Sprint(metric.Delta))
			default:
				fmt.Printf("Unknown metric's type: %T\n", metric.MType)
			}
		}

		fmt.Println("Sending finished")
	}
}

func (agent metricAgent) Start(stopChan chan struct{}) {
	fmt.Println("Agent started")
	fmt.Printf("Host: %s\n", agent.host)
	fmt.Printf("ReportInterval: %d\n", agent.intervals.send)
	fmt.Printf("PollInterval: %d\n", agent.intervals.collect)

	go agent.collectMetrics()
	go agent.sendMetrics()

	<-stopChan
	fmt.Println("Agent finished")
}

func NewAgent(host string, reportInterval, pollInterval int) *metricAgent {
	return &metricAgent{
		storage:   storage.NewStorage(),
		pollCount: 0,
		host:      host,
		intervals: agentIntervals{
			collect: pollInterval,
			send:    reportInterval,
		},
	}
}
