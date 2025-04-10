// Пакет "хранитель". Работает с хранилищем.
package keeper

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/LexusEgorov/goMetrics/internal/dohsimpson"
	"github.com/LexusEgorov/goMetrics/internal/models"
	"github.com/LexusEgorov/goMetrics/internal/transport"
)

// Интерфейс "Хранилище" для работы с метриками.
type Storager interface {
	AddGauge(key string, value float64)
	AddCounter(key string, value int64)
	GetGauge(key string) (float64, bool)
	GetCounter(key string) (int64, bool)
	GetAll() map[string]models.Metric
	MassSave(metrics []models.Metric) ([]models.Metric, error)
	Check() bool
	Close()
}

type keeper struct {
	storage Storager
}

// Метод для получения метрики из хранилища.
func (k keeper) Read(key, mType string) (*models.Metric, *dohsimpson.Error) {
	currentMetric := models.Metric{
		ID:    key,
		MType: mType,
	}

	var isFound = false

	switch mType {
	case "gauge":
		gaugeValue, found := k.storage.GetGauge(key)

		currentMetric.Value = &gaugeValue
		isFound = found
	case "counter":
		counterValue, found := k.storage.GetCounter(key)

		currentMetric.Delta = &counterValue
		isFound = found
	default:
		return nil, dohsimpson.NewDoh(http.StatusNotFound, fmt.Sprintf("reader: wrong mType (%s)", mType))
	}

	if !isFound {
		return nil, dohsimpson.NewDoh(http.StatusNotFound, fmt.Sprintf("reader: metric not found: %s (%s)", key, mType))
	}

	return &currentMetric, nil
}

// Метод для получения всех метрик из хранилища.
func (k keeper) ReadAll() map[string]models.Metric {
	return k.storage.GetAll()
}

// Метод для сохранения метрик старого типа.
func (k keeper) SaveOld(mName, mType, mValue string) *dohsimpson.Error {
	if mName == "" || mValue == "" {
		return dohsimpson.NewDoh(http.StatusNotFound, "metric not found (empty data) (saver)")
	}

	switch mType {
	case "gauge":
		value, err := strconv.ParseFloat(mValue, 64)

		if err != nil {
			return dohsimpson.NewDoh(http.StatusBadRequest, err.Error())
		}

		k.storage.AddGauge(mName, float64(value))
	case "counter":
		value, err := strconv.ParseInt(mValue, 0, 64)

		if err != nil {
			return dohsimpson.NewDoh(http.StatusBadRequest, err.Error())
		}

		k.storage.AddCounter(mName, int64(value))
	default:
		return dohsimpson.NewDoh(http.StatusBadRequest, fmt.Sprintf("unknown metric Type (%s) (saver)", mType))
	}

	return nil
}

// Метод для сохранения метрик нового типа.
func (k keeper) Save(m models.Metric) (*models.Metric, *dohsimpson.Error) {
	mName := m.ID
	mType := m.MType
	mValue := m.Value
	mDelta := m.Delta

	if mName == "" {
		return nil, dohsimpson.NewDoh(http.StatusNotFound, "saver: empty metric ID")
	}

	if mValue == nil && mDelta == nil {
		return nil, dohsimpson.NewDoh(http.StatusBadRequest, "saver: empty metric value and delta")
	}

	switch mType {
	case "gauge":
		k.storage.AddGauge(mName, float64(*mValue))
	case "counter":
		k.storage.AddCounter(mName, int64(*mDelta))
	default:
		return nil, dohsimpson.NewDoh(http.StatusBadRequest, fmt.Sprintf("saver: unknown metric Type (%s)", mType))
	}

	return &m, nil
}

// Метод для массового сохранения метрик.
func (k keeper) SaveBatch(m []models.Metric) ([]models.Metric, *dohsimpson.Error) {
	savedMetrics, err := k.storage.MassSave(m)

	if err != nil {
		return nil, dohsimpson.NewDoh(http.StatusInternalServerError, err.Error())
	}

	return savedMetrics, nil
}

// Метод для проверки работоспособности хранилища.
func (k keeper) Check() bool {
	return k.storage.Check()
}

// Конструктор. Подключает хранилище к "хранителю".
func NewKeeper(storage Storager) transport.Keeper {
	return keeper{
		storage: storage,
	}
}
