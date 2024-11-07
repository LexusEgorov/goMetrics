package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/LexusEgorov/goMetrics/internal/dohSimpson"
	"github.com/LexusEgorov/goMetrics/internal/middleware"
	"github.com/LexusEgorov/goMetrics/internal/models"
)

type Saver interface {
	SaveOld(mName, mType, value string) *dohSimpson.Error
	Save(m models.Metric) (*models.Metric, *dohSimpson.Error)
}

type Reader interface {
	Read(key, mType string) (*models.Metric, *dohSimpson.Error)
	ReadAll() map[string]models.Metric
}

type transportLayer struct {
	reader Reader
	saver  Saver
}

func (t transportLayer) UpdateMetricOld(w http.ResponseWriter, r *http.Request) {
	mName := r.PathValue("metricName")
	mType := r.PathValue("metricType")
	mValue := r.PathValue("metricValue")

	saveError := t.saver.SaveOld(mName, mType, mValue)

	if saveError != nil {
		w.WriteHeader(saveError.Code)
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

	var currentMetric models.Metric
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

	savedMetric, saveError := t.saver.Save(currentMetric)

	if saveError != nil {
		w.WriteHeader(saveError.Code)
		return
	}

	response, err := json.Marshal(savedMetric)

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
		metric, isFound = t.reader.GetGauge(mName)
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

	var currentMetric models.Metric
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
	Metrics map[string]models.Metric
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

	currentMetric := models.Metric{
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

func NewTransport(saver Saver, logger *zap.SugaredLogger) *transportLayer {
	r := chi.NewRouter()
	r.Use()
	r.Get("/", middleware.WithLogging(http.HandlerFunc(transportLayer.GetMetrics), sugar))
	r.Get("/value/{metricType}/{metricName}", middleware.WithLogging(http.HandlerFunc(transportLayer.GetMetricOld), sugar))
	r.Post("/value/", middleware.WithLogging(http.HandlerFunc(transportLayer.GetMetric), sugar))
	r.Post("/update/{metricType}/{metricName}/{metricValue}", middleware.WithLogging(http.HandlerFunc(transportLayer.UpdateMetricOld), sugar))
	r.Post("/update/", middleware.WithLogging(http.HandlerFunc(transportLayer.UpdateMetric), sugar))

	return &transportLayer{
		saver: saver,
	}
}
