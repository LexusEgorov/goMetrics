package filestorage

import (
	"encoding/json"
	"os"
	"time"

	"github.com/LexusEgorov/goMetrics/internal/dohsimpson"
	"github.com/LexusEgorov/goMetrics/internal/models"
	"github.com/LexusEgorov/goMetrics/internal/services/saver"
)

type FileWriter struct {
	path string
	file *os.File
}

func (f FileWriter) Save(metrics map[string]models.Metric) {
	jsonedMetrics, err := json.Marshal(metrics)

	if err != nil {
		return
	}

	if err = f.file.Truncate(0); err != nil {
		return
	}

	f.file.Write(jsonedMetrics)
}

func (f FileWriter) RunSave(storage saver.Storager, interval int) {
	for {
		f.Save(storage.GetAll())
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func NewFileWriter(filepath string) saver.FileWriter {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		dohsimpson.NewDoh(0, err.Error())
		return nil
	}
	defer file.Close()

	return &FileWriter{
		path: filepath,
		file: file,
	}
}

type fileReader struct {
	path string
	file *os.File
}

func (f fileReader) Read() map[string]models.Metric {
	var metrics []byte
	_, err := f.file.Read(metrics)

	if err != nil {
		return make(map[string]models.Metric)
	}

	var parsedMetrics map[string]models.Metric

	err = json.Unmarshal(metrics, &parsedMetrics)

	if err != nil {
		return make(map[string]models.Metric)
	}

	return parsedMetrics
}

func NewFileReader(filepath string) *fileReader {
	file, err := os.OpenFile(filepath, os.O_RDONLY|os.O_CREATE, 0666)

	if err != nil {
		return nil
	}
	defer file.Close()

	return &fileReader{
		path: filepath,
		file: file,
	}
}
