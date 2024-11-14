package filestorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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

func (f FileWriter) RunSave(storage saver.Storager, interval int) {
	for {
		f.Save(storage.GetAll())
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func (f FileWriter) Close() {
	defer f.file.Close()
}

func NewFileWriter(filepath string) saver.FileWriter {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		dohsimpson.NewDoh(0, err.Error())
		return nil
	}

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

func (f fileReader) Close() {
	defer f.file.Close()
}

func NewFileReader(filepath string) *fileReader {
	file, err := os.OpenFile(filepath, os.O_RDONLY|os.O_CREATE, 0666)

	if err != nil {
		return nil
	}

	return &fileReader{
		path: filepath,
		file: file,
	}
}
