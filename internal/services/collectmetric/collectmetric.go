// Пакет для сбора метрик.
package collectmetric

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"

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
	stopChan   chan struct{}
	wg         sync.WaitGroup
}

func (agent *metricAgent) collectMetrics() {
	defer agent.wg.Done()
	for {
		select {
		case _, ok := <-agent.stopChan:
			if ok {
				fmt.Println("collect stop")
				return
			}

			fmt.Println("collect stop (already)")
			return
		default:
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
}

func (agent *metricAgent) sendMetrics(isLast bool) {
	transportClient := transport.NewClient()
	semaphore := make(chan struct{}, agent.config.RateLimit)

	for {
		select {
		case _, ok := <-agent.stopChan:
			if isLast {
				continue
			}

			if ok {
				fmt.Println("send stop")
				return
			}

			fmt.Println("send stop (already)")
			return
		default:
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
}

func (agent *metricAgent) collectAdditionals() {
	for {
		select {
		case _, ok := <-agent.stopChan:
			if ok {
				fmt.Println("add stop")
				return
			}

			fmt.Println("add stop (already)")
			return
		default:
			v, _ := mem.VirtualMemory()
			c, _ := cpu.Counts(false)

			cpuCount := float64(c)

			agent.metricChan <- models.Metric{
				ID:    "CPUutilization1",
				MType: "gauge",
				Value: &cpuCount,
			}

			total := float64(v.Total)
			agent.metricChan <- models.Metric{
				ID:    "TotalMemory",
				MType: "gauge",
				Value: &total,
			}

			free := float64(v.Free)
			agent.metricChan <- models.Metric{
				ID:    "FreeMemory",
				MType: "gauge",
				Value: &free,
			}

			time.Sleep(time.Duration(agent.config.PollInterval) * time.Second)
		}
	}
}

// Метод, который запускает сбор метрик.
func (agent *metricAgent) Start(stopChan chan struct{}) {
	fmt.Println("Agent started")
	fmt.Printf("Host: %s\n", agent.config.Host)
	fmt.Printf("ReportInterval: %d\n", agent.config.ReportInterval)
	fmt.Printf("PollInterval: %d\n", agent.config.PollInterval)
	fmt.Printf("RateLimit: %d\n", agent.config.RateLimit)
	fmt.Printf("Key: %s\n", agent.config.Key)

	agent.stopChan = stopChan

	go agent.collectMetrics()
	go agent.sendMetrics(false)
	go agent.collectAdditionals()

	shutdown := false

	for !shutdown {
		select {
		case <-stopChan:
			shutdown = true
			agent.sendMetrics(true)
		default:
		}
	}

	close(agent.metricChan)
	agent.wg.Wait()
	fmt.Println("Agent finished")
}

// Конструктор агента.
func NewAgent(config config.Agent, signer middleware.Signer) *metricAgent {
	return &metricAgent{
		config:     config,
		pollCount:  0,
		signer:     signer,
		metricChan: make(chan models.Metric, 100),
	}
}
