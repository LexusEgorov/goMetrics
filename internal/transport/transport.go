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
	AddGauge(key string, value float64)
	AddCounter(key string, value int64)
	GetGauge(key string) (float64, bool)
	GetCounter(key string) (int64, bool)
	GetAll() map[string]storage.Metric
}

type transportLayer struct {
	storage Storager
}

func (t transportLayer) UpdateMetricOld(w http.ResponseWriter, r *http.Request) {
	mName := r.PathValue("metricName")
	mType := r.PathValue("metricType")
	mValue := r.PathValue("metricValue")

	fmt.Printf("Name: %s\nType: %s\nValue: %s\n ============\n", mName, mType, mValue)

	if mName == "" || mValue == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch mType {
	case "gauge":
		value, err := strconv.ParseFloat(mValue, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		t.storage.AddGauge(mName, float64(value))
	case "counter":
		value, err := strconv.ParseInt(mValue, 0, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		t.storage.AddCounter(mName, int64(value))
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t transportLayer) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
		t.storage.AddGauge(mName, float64(*mValue))
	case "counter":
		t.storage.AddCounter(mName, int64(*mDelta))
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(currentMetric)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	w.WriteHeader(http.StatusOK)
}

func (t transportLayer) GetMetricOld(w http.ResponseWriter, r *http.Request) {
	mName := r.PathValue("metricName")
	mType := r.PathValue("metricType")

	var metric interface{}
	var isFound bool

	switch mType {
	case "gauge":
		metric, isFound = t.storage.GetGauge(mName)
	case "counter":
		metric, isFound = t.storage.GetCounter(mName)
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !isFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write([]byte(fmt.Sprint(metric)))
}

func (t transportLayer) GetMetric(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != "application/json" {
		t.GetMetricOld(w, r)
		return
	}

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

	switch currentMetric.MType {
	case "gauge":
		metric, isFound := t.storage.GetGauge(currentMetric.ID)

		if isFound {
			value := float64(metric)
			currentMetric.Value = &value
		}
	case "counter":
		metric, isFound := t.storage.GetCounter(currentMetric.ID)
		if isFound {
			value := int64(metric)
			currentMetric.Delta = &value
		}
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if currentMetric.Delta == nil && currentMetric.Value == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := json.Marshal(currentMetric)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	w.WriteHeader(http.StatusOK)
}

type PageData struct {
	Title   string
	Header  string
	Metrics map[string]storage.Metric
}

func (t transportLayer) GetMetrics(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != "application/json" {
		t.GetMetricsOld(w, r)
		return
	}

	metrics := t.storage.GetAll()
	response, err := json.Marshal(metrics)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	w.WriteHeader(http.StatusOK)
}

func (t transportLayer) GetMetricsOld(w http.ResponseWriter, r *http.Request) {
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
					{{range .Metrics}}
						<li>ID: {{.ID}}, Type: {{.MType}}, Delta: {{if .Delta}}{{.Delta}}{{else}}N/A{{end}}, Value: {{if .Value}}{{.Value}}{{else}}N/A{{end}}</li>
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
	//url := fmt.Sprintf("http://%s/update", host)

	client := resty.New()

	currentMetric := Metric{
		ID:    metricName,
		MType: metricType,
	}

	if metricType == "gauge" {
		gaugeValue, err := strconv.ParseFloat(metricValue, 64)

		if err != nil {
			fmt.Printf("VALUE: %s, ERROR: %s\n", metricValue, err)
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

func NewTransport() *transportLayer {
	return &transportLayer{
		storage: storage.NewStorage(),
	}
}
