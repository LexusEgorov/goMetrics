package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/go-resty/resty/v2"

	"github.com/LexusEgorov/goMetrics/internal/services/storage"
)

type Metric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type Storager interface {
	AddGauge(key storage.MetricName, value storage.Gauge)
	AddCounter(key storage.MetricName, value storage.Counter)
	GetGauge(key storage.MetricName) (storage.Gauge, bool)
	GetCounter(key storage.MetricName) (storage.Counter, bool)
	GetAll() map[storage.MetricName]interface{}
}

type transportLayer struct {
	storage Storager
}

func (t transportLayer) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	var currentMetric Metric
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &currentMetric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mName := currentMetric.ID
	mType := currentMetric.MType
	mValue := currentMetric.Value
	mDelta := currentMetric.Delta

	if mName == "" || (mValue == nil && mDelta == nil) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch mType {
	case "gauge":
		t.storage.AddGauge(storage.MetricName(mName), storage.Gauge(*mValue))
	case "counter":
		t.storage.AddCounter(storage.MetricName(mName), storage.Counter(*mDelta))
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t transportLayer) GetMetric(w http.ResponseWriter, r *http.Request) {
	mName := r.PathValue("metricName")
	mType := r.PathValue("metricType")

	var metric fmt.Stringer
	var isFound bool

	switch mType {
	case "gauge":
		metric, isFound = t.storage.GetGauge(storage.MetricName(mName))
	case "counter":
		metric, isFound = t.storage.GetCounter(storage.MetricName(mName))
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !isFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write([]byte(metric.String()))
}

type PageData struct {
	Title   string
	Header  string
	Metrics map[storage.MetricName]interface{}
}

func (t transportLayer) GetMetrics(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		Title:   "Metrics",
		Header:  "Metrics list: ",
		Metrics: t.storage.GetAll(),
	}

	page, err := template.New("webpage").
		Parse(`
			<!DOCTYPE html>
				<html lang="ru">
				<head>
					<meta charset="UTF-8">
					<title>{{.Title}}</title>
				</head>
				<body>
					<h1>{{.Header}}</h1>
					<ul>
						{{range $key, $value := .Metrics}}
							<li>{{$key}}: {{$value}}</li>
						{{end}}
					</ul>
				</body>
			</html>
		`)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = page.Execute(w, pageData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t transportLayer) SendMetric(host, metricName, metricType, metricValue string) {
	url := fmt.Sprintf("http://%s/update/%s/%s/%s", host, metricType, metricName, metricValue)

	client := resty.New()

	currentMetric := Metric{
		ID:    metricName,
		MType: metricType,
	}

	if metricType == "gauge" {
		gaugeValue, err := strconv.ParseFloat(metricValue, 64)

		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		currentMetric.Value = &gaugeValue
	} else if metricType == "counter" {
		counterValue, err := strconv.Atoi(metricValue)

		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		deltaValue := int64(counterValue)
		currentMetric.Delta = &deltaValue
	} else {
		fmt.Printf("ERROR: unknown metricType\n")
		return
	}

	body, err := json.Marshal(currentMetric)

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
}

func CreateTransport() *transportLayer {
	return &transportLayer{
		storage: storage.CreateStorage(),
	}
}
