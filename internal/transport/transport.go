package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/LexusEgorov/goMetrics/internal/dohsimpson"
	"github.com/LexusEgorov/goMetrics/internal/middleware"
	"github.com/LexusEgorov/goMetrics/internal/models"
)

type Saver interface {
	SaveOld(mName, mType, value string) *dohsimpson.Error
	Save(m models.Metric) (*models.Metric, *dohsimpson.Error)
}

type Reader interface {
	Read(key, mType string) (*models.Metric, *dohsimpson.Error)
	ReadAll() map[string]models.Metric
}

type transportServer struct {
	Router *chi.Mux
	reader Reader
	saver  Saver
}

func (t transportServer) UpdateMetricOld(w http.ResponseWriter, r *http.Request) {
	mName := r.PathValue("metricName")
	mType := r.PathValue("metricType")
	mValue := r.PathValue("metricValue")

	saveError := t.saver.SaveOld(mName, mType, mValue)

	if saveError != nil {
		w.WriteHeader(saveError.Code)
		return
	}
}

func (t transportServer) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var currentMetric models.Metric
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(body, &currentMetric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	savedMetric, saveError := t.saver.Save(currentMetric)

	if saveError != nil {
		w.WriteHeader(saveError.Code)
		return
	}

	updatedMetric, readError := t.reader.Read(savedMetric.ID, savedMetric.MType)

	if readError != nil {
		w.WriteHeader(readError.Code)
		return
	}

	response, err := json.Marshal(updatedMetric)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (t transportServer) GetMetricOld(w http.ResponseWriter, r *http.Request) {
	mName := r.PathValue("metricName")
	mType := r.PathValue("metricType")

	currentMetric := models.Metric{
		ID:    mName,
		MType: mType,
	}

	foundMetric, readError := t.reader.Read(currentMetric.ID, currentMetric.MType)

	if readError != nil {
		w.WriteHeader(readError.Code)
		return
	}

	switch currentMetric.MType {
	case "gauge":
		w.Write([]byte(fmt.Sprint(*foundMetric.Value)))
		return
	case "counter":
		w.Write([]byte(fmt.Sprint(*foundMetric.Delta)))
		return
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (t transportServer) GetMetric(w http.ResponseWriter, r *http.Request) {
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

	foundMetric, readError := t.reader.Read(currentMetric.ID, currentMetric.MType)

	if readError != nil {
		w.WriteHeader(readError.Code)
		return
	}

	response, err := json.Marshal(foundMetric)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

type PageData struct {
	Title   string
	Header  string
	Metrics map[string]models.Metric
}

func (t transportServer) GetMetrics(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != "application/json" {
		t.GetMetricsOld(w, r)
		return
	}

	metrics := t.reader.ReadAll()
	response, err := json.Marshal(metrics)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (t transportServer) GetMetricsOld(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		Title:   "Metrics",
		Header:  "Metrics list: ",
		Metrics: t.reader.ReadAll(),
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

func NewServer(saver Saver, reader Reader, router *chi.Mux, logger *zap.SugaredLogger) *transportServer {
	transportServer := transportServer{
		Router: router,
		reader: reader,
		saver:  saver,
	}

	router.Use(middleware.WithLogging(logger))
	router.Use(middleware.WithDecoding)
	router.Use(middleware.WithEncoding)

	router.Post("/update/", http.HandlerFunc(transportServer.UpdateMetric))
	router.Get("/", http.HandlerFunc(transportServer.GetMetrics))
	router.Post("/value/", http.HandlerFunc(transportServer.GetMetric))
	router.Get("/value/{metricType}/{metricName}", http.HandlerFunc(transportServer.GetMetricOld))
	router.Post("/update/{metricType}/{metricName}/{metricValue}", http.HandlerFunc(transportServer.UpdateMetricOld))

	return &transportServer
}

type transportClient struct{}

func (t transportClient) SendMetric(host, metricName, metricType, metricValue string) {
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

func NewClient() *transportClient {
	return &transportClient{}
}
