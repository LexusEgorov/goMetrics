package db

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/LexusEgorov/goMetrics/internal/dohsimpson"
	"github.com/LexusEgorov/goMetrics/internal/keeper"
	"github.com/LexusEgorov/goMetrics/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	db *sql.DB
}

func (d *DB) connect(host string) {
	var err error
	d.db, err = sql.Open("pgx", host)

	if err != nil {
		dohsimpson.NewDoh(http.StatusInternalServerError, err.Error())
		return
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS metrics (
		id VARCHAR(255) PRIMARY KEY,
		mtype VARCHAR(50) NOT NULL,
		delta BIGINT,
		value DOUBLE PRECISION
	);`

	_, err = d.db.Exec(createTableSQL)

	if err != nil {
		dohsimpson.NewDoh(http.StatusInternalServerError, "DB (createTable): "+err.Error())
		d.db = nil
		return
	}
}

func (d DB) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

func (d DB) Check() bool {
	return d.db != nil
}

func (d DB) MassSave(metrics []models.Metric) ([]models.Metric, error) {
	savedMetrics := make([]string, len(metrics))

	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	for i, metric := range metrics {
		savedMetrics[i] = metric.ID

		switch metric.MType {
		case "gauge":
			query := `
				INSERT INTO metrics (id, mtype, value) 
				VALUES ($1, 'gauge', $2)
				ON CONFLICT (id) 
				DO UPDATE SET value = EXCLUDED.value;`

			_, err := d.db.Exec(query, metric.ID, metric.Value)

			if err != nil {
				dohsimpson.NewDoh(http.StatusInternalServerError, "DB (Mass:addGauge): "+err.Error())
			}

		case "counter":
			query := `
				INSERT INTO metrics (id, mtype, delta) 
				VALUES ($1, 'counter', $2)
				ON CONFLICT (id) 
				DO UPDATE SET delta = metrics.delta + EXCLUDED.delta;`

			_, err := d.db.Exec(query, metric.ID, metric.Delta)

			if err != nil {
				dohsimpson.NewDoh(http.StatusInternalServerError, "DB (Mass:addCounter): "+err.Error())
			}
		}
	}

	err = tx.Commit()

	if err != nil {
		dohsimpson.NewDoh(http.StatusInternalServerError, "DB (Mass:commit): "+err.Error())
	}

	query := `SELECT * FROM metrics WHERE id IN (&1)`
	rows, err := d.db.Query(query, strings.Join(savedMetrics, ", "))

	if err != nil {
		dohsimpson.NewDoh(http.StatusInternalServerError, "DB (Mass:read): "+err.Error())
		return nil, err
	}

	defer rows.Close()

	resultMetrics := make([]models.Metric, 0)
	for rows.Next() {
		var m models.Metric
		err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)

		if err != nil {
			dohsimpson.NewDoh(http.StatusInternalServerError, "DB (Mass:read row): "+err.Error())
			return nil, err
		}

		resultMetrics = append(resultMetrics, m)
	}

	err = rows.Err()

	if err != nil {
		dohsimpson.NewDoh(http.StatusInternalServerError, "DB (Mass:read rows): "+err.Error())
		return nil, err
	}

	return resultMetrics, nil
}

func (d DB) AddCounter(key string, value int64) {
	query := `
		INSERT INTO metrics (id, mtype, delta) 
		VALUES ($1, 'counter', $2)
		ON CONFLICT (id) 
		DO UPDATE SET delta = metrics.delta + EXCLUDED.delta;`

	_, err := d.db.Exec(query, key, value)

	if err != nil {
		dohsimpson.NewDoh(http.StatusInternalServerError, "DB (addCounter): "+err.Error())
	}
}

func (d DB) AddGauge(key string, value float64) {
	query := `
	INSERT INTO metrics (id, mtype, value) 
	VALUES ($1, 'gauge', $2)
	ON CONFLICT (id) 
	DO UPDATE SET value = EXCLUDED.value;`

	_, err := d.db.Exec(query, key, value)

	if err != nil {
		dohsimpson.NewDoh(http.StatusInternalServerError, "DB (addGauge): "+err.Error())
	}
}

func (d DB) GetAll() map[string]models.Metric {
	metrics := make(map[string]models.Metric)
	query := `SELECT * FROM metrics`

	rows, err := d.db.Query(query)

	if err != nil {
		dohsimpson.NewDoh(http.StatusInternalServerError, "DB (getAll): "+err.Error())
		return nil
	}

	defer rows.Close()

	for rows.Next() {
		var m models.Metric
		err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)

		if err != nil {
			dohsimpson.NewDoh(http.StatusInternalServerError, "DB (getAll row): "+err.Error())
			return nil
		}

		metrics[m.ID] = m
	}

	err = rows.Err()

	if err != nil {
		dohsimpson.NewDoh(http.StatusInternalServerError, "DB (getAll rows): "+err.Error())
		return nil
	}

	return metrics
}

func (d DB) getMetric(key string) (*models.Metric, bool) {
	query := `SELECT * FROM metrics WHERE id = $1`

	row := d.db.QueryRow(query, key)

	var m models.Metric
	err := row.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)

	if err != nil {
		dohsimpson.NewDoh(http.StatusInternalServerError, "DB (getMetric): "+err.Error())
		return nil, false
	}

	return &m, true
}

func (d DB) GetCounter(key string) (int64, bool) {
	metric, isFound := d.getMetric(key)

	if !isFound {
		return 0, false
	}

	return int64(*metric.Delta), true
}

func (d DB) GetGauge(key string) (float64, bool) {
	metric, isFound := d.getMetric(key)

	if !isFound {
		return 0, false
	}

	return float64(*metric.Value), true
}

func NewDB(host string) keeper.Storager {
	db := DB{}

	db.connect(host)

	return db
}
