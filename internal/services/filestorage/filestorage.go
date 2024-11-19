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

func (f fileStorage) save(metrics map[string]models.Metric) {
	jsonedMetrics, err := json.Marshal(metrics)

	if err != nil {
		fmt.Println(err)
		return
	}

	if err = f.file.Truncate(0); err != nil {
		fmt.Println(err)
		return
	}

	if _, err = f.file.Seek(0, 0); err != nil {
		fmt.Println(err)
		return
	}

	f.file.Write(jsonedMetrics)
}

func (f fileStorage) runSave(interval int) {
	for {
		f.save(f.storage.GetAll())
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func (f fileStorage) read() map[string]models.Metric {
	reader := bufio.NewReader(f.file)
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

func (f fileStorage) Close() {
	f.file.Close()
}

func (fs fileStorage) AddGauge(key string, value float64) {
	fs.storage.AddGauge(key, value)

	if fs.interval == 0 {
		fs.save(fs.storage.GetAll())
	}
}

func (fs fileStorage) AddCounter(key string, value int64) {
	fs.storage.AddCounter(key, value)

	if fs.interval == 0 {
		fs.save(fs.storage.GetAll())
	}
}

func (fs fileStorage) GetGauge(key string) (float64, bool) {
	return fs.storage.GetGauge(key)
}

func (fs fileStorage) GetCounter(key string) (int64, bool) {
	return fs.storage.GetCounter(key)
}

func (fs fileStorage) GetAll() map[string]models.Metric {
	return fs.storage.GetAll()
}

func (f fileStorage) Check() bool {
	return true
}

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
