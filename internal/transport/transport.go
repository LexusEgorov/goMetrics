// Пакет транспортного уровня. Отвечает за обработку http запросов сервером и отправку метрик агентом.
package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	debugMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/LexusEgorov/goMetrics/internal/dohsimpson"
	"github.com/LexusEgorov/goMetrics/internal/errors"
	"github.com/LexusEgorov/goMetrics/internal/middleware"
	"github.com/LexusEgorov/goMetrics/internal/models"
)

// Интерфейс "хранитель". Работает с метриками.
type Keeper interface {
	SaveOld(mName, mType, value string) *dohsimpson.Error
	Save(m models.Metric) (*models.Metric, *dohsimpson.Error)
	SaveBatch(m []models.Metric) ([]models.Metric, *dohsimpson.Error)
	Read(key, mType string) (*models.Metric, *dohsimpson.Error)
	ReadAll() map[string]models.Metric
	Check() bool
}

func SendClosed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusServiceUnavailable)
}

type transportServer struct {
	Router *chi.Mux
	keeper Keeper
	isExit bool
}

type pageData struct {
	Title   string
	Header  string
	Metrics map[string]models.Metric
}

func (t *transportServer) Shutdown() {
	t.isExit = true

	time.Sleep(time.Second * 5)
}

// Обработчик обновления метрики старого типа.
func (t transportServer) UpdateMetricOld(w http.ResponseWriter, r *http.Request) {
	if t.isExit {
		SendClosed(w)
		return
	}

	mName := r.PathValue("metricName")
	mType := r.PathValue("metricType")
	mValue := r.PathValue("metricValue")

	saveError := t.keeper.SaveOld(mName, mType, mValue)

	if saveError != nil {
		w.WriteHeader(saveError.Code)
		return
	}
}

// Обработчик обновления метрики.
func (t transportServer) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	if t.isExit {
		SendClosed(w)
		return
	}
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

	savedMetric, saveError := t.keeper.Save(currentMetric)

	if saveError != nil {
		w.WriteHeader(saveError.Code)
		return
	}

	updatedMetric, readError := t.keeper.Read(savedMetric.ID, savedMetric.MType)

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

// Обработчик массового обновления метрик.
func (t transportServer) UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	if t.isExit {
		SendClosed(w)
		return
	}
	contentType := r.Header.Get("Content-Type")

	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var metrics []models.Metric
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(body, &metrics); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	savedMetrics, saveError := t.keeper.SaveBatch(metrics)

	if saveError != nil {
		w.WriteHeader(saveError.Code)
		return
	}

	response, err := json.Marshal(savedMetrics)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// Обработчик получения метрик старого типа.
func (t transportServer) GetMetricOld(w http.ResponseWriter, r *http.Request) {
	if t.isExit {
		SendClosed(w)
		return
	}
	mName := r.PathValue("metricName")
	mType := r.PathValue("metricType")

	currentMetric := models.Metric{
		ID:    mName,
		MType: mType,
	}

	foundMetric, readError := t.keeper.Read(currentMetric.ID, currentMetric.MType)

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

// Обработчик получения метрики.
func (t transportServer) GetMetric(w http.ResponseWriter, r *http.Request) {
	if t.isExit {
		SendClosed(w)
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

	foundMetric, readError := t.keeper.Read(currentMetric.ID, currentMetric.MType)

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

// Обработчик получения списка метрик.
func (t transportServer) GetMetrics(w http.ResponseWriter, r *http.Request) {
	if t.isExit {
		SendClosed(w)
		return
	}

	contentType := r.Header.Get("Content-Type")

	if contentType != "application/json" {
		t.GetMetricsOld(w, r)
		return
	}

	metrics := t.keeper.ReadAll()
	response, err := json.Marshal(metrics)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// Обработчик получения списка метрик старого типа.
func (t transportServer) GetMetricsOld(w http.ResponseWriter, r *http.Request) {
	if t.isExit {
		SendClosed(w)
		return
	}
	pageData := pageData{
		Title:   "Metrics",
		Header:  "Metrics list: ",
		Metrics: t.keeper.ReadAll(),
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

	w.Header().Set("Content-Type", "text/html")
}

// Обработчик проверки связи с хранилищем.
func (t transportServer) CheckDB(w http.ResponseWriter, r *http.Request) {
	if t.isExit {
		SendClosed(w)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	done := make(chan bool)

	go func() {
		done <- t.keeper.Check()
	}()

	select {
	case success := <-done:
		if success {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
	case <-ctx.Done():
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Конструктор сервера.
func NewServer(keeper Keeper, router *chi.Mux, logger *zap.SugaredLogger, signer middleware.Signer, accept string) *transportServer {
	transportServer := transportServer{
		Router: router,
		keeper: keeper,
	}

	router.Use(middleware.WithLogging(logger))
	router.Use(middleware.WithAccepting(accept))
	router.Use(middleware.WithDecoding)
	router.Use(middleware.WithVerifying(signer))
	router.Use(middleware.WithEncoding)
	router.Use(middleware.WithSigning(signer))

	router.Get("/", http.HandlerFunc(transportServer.GetMetrics))
	router.Get("/ping", http.HandlerFunc(transportServer.CheckDB))
	router.Post("/value/", http.HandlerFunc(transportServer.GetMetric))
	router.Get("/value/{metricType}/{metricName}", http.HandlerFunc(transportServer.GetMetricOld))
	router.Post("/update/", http.HandlerFunc(transportServer.UpdateMetric))
	router.Post("/update/{metricType}/{metricName}/{metricValue}", http.HandlerFunc(transportServer.UpdateMetricOld))
	router.Post("/updates/", http.HandlerFunc(transportServer.UpdateMetrics))

	router.Mount("/debug", debugMiddleware.Profiler())

	return &transportServer
}

type transportClient struct {
	ip string
}

// Метод для отправки метрики.
func (t transportClient) SendMetric(host, metricName, metricType, metricValue string, signer middleware.Signer) {
	const maxRetries = 3

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

	sign := signer.Sign(body)

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	for retriesCount := 0; retriesCount < maxRetries; retriesCount++ {
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("HashSHA256", sign).
			SetHeader("X-Real-IP", t.ip).
			SetBody(body).
			Post(url)

		if err != nil {
			if errors.IsClientRetriable(resp.StatusCode()) {
				sleepDuration := retriesCount*2 + 1

				fmt.Printf("Error. Attempt: %d/%d Retry in %ds.\n", retriesCount+1, maxRetries, sleepDuration)
				time.Sleep(time.Second * time.Duration(sleepDuration))
				continue
			}

			return
		}
	}
}

// Конструктор агента.
func NewClient() *transportClient {
	ip := ""
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range ifaces {
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}

			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.IsPrivate() {
					ip = ipnet.IP.String()
				}
			}
		}

		fmt.Printf("Agent IP: %s\n", ip)
	}

	return &transportClient{
		ip: ip,
	}
}
