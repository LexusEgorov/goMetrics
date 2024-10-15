package transport

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/go-resty/resty/v2"

	"github.com/LexusEgorov/goMetrics/internal/services/storage"
)

type Transporter interface {
	UpdateMetric(w http.ResponseWriter, r *http.Request)
	GetMetric(w http.ResponseWriter, r *http.Request)
	GetMetrics(w http.ResponseWriter, r *http.Request)
	SendMetric(host, metricName, metricType, metricValue string)
}

type TransportLayer struct {
	storage storage.Storager
}

func (t TransportLayer) UpdateMetric(w http.ResponseWriter, r *http.Request) {
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

		t.storage.AddGauge(storage.MetricName(mName), storage.Gauge(value))
	case "counter":
		value, err := strconv.ParseInt(mValue, 0, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		t.storage.AddCounter(storage.MetricName(mName), storage.Counter(value))
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t TransportLayer) GetMetric(w http.ResponseWriter, r *http.Request) {
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

func (t TransportLayer) GetMetrics(w http.ResponseWriter, r *http.Request) {
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

func (t TransportLayer) SendMetric(host, metricName, metricType, metricValue string) {
	url := fmt.Sprintf("http://%s/update/%s/%s/%s", host, metricType, metricName, metricValue)

	client := resty.New()

	_, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)

	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}
}

func CreateTransport() Transporter {
	return TransportLayer{
		storage: storage.CreateStorage(),
	}
}
