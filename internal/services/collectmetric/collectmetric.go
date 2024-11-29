package collectmetric

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/LexusEgorov/goMetrics/internal/config"
	"github.com/LexusEgorov/goMetrics/internal/middleware"
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
	config     config.Agent
	signer     middleware.Signer
	pollCount  int64
	metricChan chan models.Metric
	wg         sync.WaitGroup
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

			agent.metricChan <- currentMetric
		}

		agent.metricChan <- models.Metric{
			ID:    "PollCount",
			MType: "counter",
			Delta: &agent.pollCount,
		}

		randomValue := rand.Float64()

		agent.metricChan <- models.Metric{
			ID:    "RandomValue",
			MType: "gauge",
			Value: &randomValue,
		}

		time.Sleep(time.Duration(agent.config.PollInterval) * time.Second)
	}
}

func (agent *metricAgent) sendMetrics() {
	transportClient := transport.NewClient()
	semaphore := make(chan struct{}, agent.config.RateLimit)

	for {
		time.Sleep(time.Duration(agent.config.ReportInterval) * time.Second)
		fmt.Println("Sending started")
		for metric := range agent.metricChan {
			semaphore <- struct{}{}
			agent.wg.Add(1)

			go func(metric models.Metric) {
				defer func() {
					<-semaphore
					agent.wg.Done()
				}()

				switch metric.MType {
				case "gauge":
					transportClient.SendMetric(agent.config.Host, string(metric.ID), metric.MType, fmt.Sprint(*metric.Value), agent.signer)
				case "counter":
					transportClient.SendMetric(agent.config.Host, string(metric.ID), metric.MType, fmt.Sprint(*metric.Delta), agent.signer)
				default:
					fmt.Printf("Unknown metric's type: %T\n", metric.MType)
				}
			}(metric)
		}

		agent.wg.Wait()
		fmt.Println("Sending finished")
	}
}

func (agent *metricAgent) Start(stopChan chan struct{}) {
	fmt.Println("Agent started")
	fmt.Printf("Host: %s\n", agent.config.Host)
	fmt.Printf("ReportInterval: %d\n", agent.config.ReportInterval)
	fmt.Printf("PollInterval: %d\n", agent.config.PollInterval)
	fmt.Printf("Key: %s\n", agent.config.Key)

	go agent.collectMetrics()
	go agent.sendMetrics()

	<-stopChan
	close(agent.metricChan)
	agent.wg.Wait()
	fmt.Println("Agent finished")
}

func NewAgent(config config.Agent, signer middleware.Signer) *metricAgent {
	return &metricAgent{
		config:    config,
		pollCount: 0,
		signer:    signer,
	}
}
