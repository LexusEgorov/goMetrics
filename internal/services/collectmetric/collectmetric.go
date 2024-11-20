package collectmetric

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"runtime"
	"time"

	"github.com/LexusEgorov/goMetrics/internal/config"
	"github.com/LexusEgorov/goMetrics/internal/models"
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

type metricAgent struct {
	config    config.Agent
	keeper    transport.Keeper
	pollCount int64
}

func (agent *metricAgent) collectMetrics() {
	for {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		for _, metricName := range gaugeMetrics {
			agent.pollCount++
			currentMetric := models.Metric{
				ID:    metricName,
				MType: "gauge",
			}

			value := reflect.ValueOf(memStats).FieldByName(metricName)

			if value.IsValid() && value.CanInterface() {
				var floatedValue float64

				switch v := value.Interface().(type) {
				case float64:
					floatedValue = v
				case uint64:
					floatedValue = float64(v)
				case uint32:
					floatedValue = float64(v)
				case uint16:
					floatedValue = float64(v)
				case uint8:
					floatedValue = float64(v)
				default:
					fmt.Printf("Unable to convert metric %s (%s) to a float64\n", metricName, v)
					continue
				}

				currentMetric.Value = &floatedValue
			} else {
				fmt.Printf("Metric %s is not valid or accessible\n", metricName)
				continue
			}

			agent.keeper.Save(currentMetric)
		}

		agent.keeper.Save(models.Metric{
			ID:    "PollCount",
			MType: "counter",
			Delta: &agent.pollCount,
		})

		randomValue := rand.Float64()
		agent.keeper.Save(models.Metric{
			ID:    "RandomValue",
			MType: "gauge",
			Value: &randomValue,
		})

		time.Sleep(time.Duration(agent.config.PollInterval) * time.Second)
	}
}

func (agent metricAgent) sendMetrics() {
	transportClient := transport.NewClient()

	for {
		time.Sleep(time.Duration(agent.config.ReportInterval) * time.Second)
		fmt.Println("Sending started")
		for k, metric := range agent.keeper.ReadAll() {
			switch metric.MType {
			case "gauge":
				transportClient.SendMetric(agent.config.Host, string(k), metric.MType, fmt.Sprint(*metric.Value))
			case "counter":
				transportClient.SendMetric(agent.config.Host, string(k), metric.MType, fmt.Sprint(*metric.Delta))
			default:
				fmt.Printf("Unknown metric's type: %T\n", metric.MType)
			}
		}

		fmt.Println("Sending finished")
	}
}

func (agent metricAgent) Start(stopChan chan struct{}) {
	fmt.Println("Agent started")
	fmt.Printf("Host: %s\n", agent.config.Host)
	fmt.Printf("ReportInterval: %d\n", agent.config.ReportInterval)
	fmt.Printf("PollInterval: %d\n", agent.config.PollInterval)

	go agent.collectMetrics()
	go agent.sendMetrics()

	<-stopChan
	fmt.Println("Agent finished")
}

func NewAgent(init config.Agent, keeper transport.Keeper) *metricAgent {
	return &metricAgent{
		config:    init,
		keeper:    keeper,
		pollCount: 0,
	}
}
