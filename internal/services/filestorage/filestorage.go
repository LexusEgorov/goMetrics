// Файловое хранилище. Один из вариантов хранилищ.
// Работает аналогично хранилищу в оперативной памяти, но имеет возможность
// восстановить данные из файла и записывать данные в файл.
package filestorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/LexusEgorov/goMetrics/internal/dohsimpson"
	"github.com/LexusEgorov/goMetrics/internal/keeper"
	"github.com/LexusEgorov/goMetrics/internal/models"
	"github.com/LexusEgorov/goMetrics/internal/services/storage"
)

type fileStorage struct {
	path     string
	file     *os.File
	interval int
	storage  keeper.Storager
}

// Метод для массового сохранения метрик.
func (fs fileStorage) MassSave(metrics []models.Metric) ([]models.Metric, error) {
	savedMetrics := make([]models.Metric, len(metrics))

	for i, metric := range metrics {
		switch metric.MType {
		case "gauge":
			fs.storage.AddGauge(metric.ID, *metric.Value)
			savedMetrics[i] = metric
		case "counter":
			oldValue, _ := fs.storage.GetCounter(metric.ID)

			fs.storage.AddCounter(metric.ID, *metric.Delta)

			newValue := *metric.Delta + oldValue

			metric.Delta = &newValue
			savedMetrics[i] = metric
		}
	}

	fs.save(fs.storage.GetAll())
	return savedMetrics, nil
}

func (fs fileStorage) save(metrics map[string]models.Metric) {
	jsonedMetrics, err := json.Marshal(metrics)

	if err != nil {
		fmt.Println(err)
		return
	}

	if err = fs.file.Truncate(0); err != nil {
		fmt.Println(err)
		return
	}

	if _, err = fs.file.Seek(0, 0); err != nil {
		fmt.Println(err)
		return
	}

	fs.file.Write(jsonedMetrics)
}

func (fs fileStorage) runSave(interval int) {
	for {
		fs.save(fs.storage.GetAll())
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func (fs fileStorage) read() map[string]models.Metric {
	reader := bufio.NewReader(fs.file)
	metrics, err := io.ReadAll(reader)

	if err != nil {
		fmt.Println(err)
		return make(map[string]models.Metric)
	}

	if len(metrics) == 0 {
		return make(map[string]models.Metric)
	}

	var parsedMetrics map[string]models.Metric

	err = json.Unmarshal(metrics, &parsedMetrics)

	if err != nil {
		fmt.Println(err)
		return make(map[string]models.Metric)
	}

	return parsedMetrics
}

// Метод для закрытия работы с файлом.
func (fs fileStorage) Close() {
	fs.file.Close()
}

// Метод для сохранения метрики типа "gauge".
func (fs fileStorage) AddGauge(key string, value float64) {
	fs.storage.AddGauge(key, value)

	if fs.interval == 0 {
		fs.save(fs.storage.GetAll())
	}
}

// Метод для сохранения метрики типа "counter".
func (fs fileStorage) AddCounter(key string, value int64) {
	fs.storage.AddCounter(key, value)

	if fs.interval == 0 {
		fs.save(fs.storage.GetAll())
	}
}

// Метод для получения значения метрики типа "gauge" по ключу.
func (fs fileStorage) GetGauge(key string) (float64, bool) {
	return fs.storage.GetGauge(key)
}

// Метод для получения значения метрики типа "counter" по ключу.
func (fs fileStorage) GetCounter(key string) (int64, bool) {
	return fs.storage.GetCounter(key)
}

// Метод для получения всех метрик.
func (fs fileStorage) GetAll() map[string]models.Metric {
	return fs.storage.GetAll()
}

// Метод заглушка. Всегда возвращает true.
func (fs fileStorage) Check() bool {
	return true
}

// Конструктор.
func NewFileStorage(filepath string, saveInterval int, isRestore bool) keeper.Storager {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		dohsimpson.NewDoh(0, err.Error())
		return nil
	}

	fileStorage := fileStorage{
		path:     filepath,
		file:     file,
		interval: saveInterval,
	}

	initMetrics := make(map[string]models.Metric)

	if isRestore {
		initMetrics = fileStorage.read()
	}

	storage := storage.NewStorage(initMetrics)

	fileStorage.storage = storage

	if saveInterval != 0 {
		go fileStorage.runSave(saveInterval)
	}

	return fileStorage
}
