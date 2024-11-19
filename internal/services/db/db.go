package db

import (
	"database/sql"
	"net/http"

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
		id SERIAL PRIMARY KEY,
		mtype VARCHAR(50) NOT NULL,
		delta BIGINT,
		value DOUBLE PRECISION
	);`

	_, err = d.db.Exec(createTableSQL)

	if err != nil {
		dohsimpson.NewDoh(http.StatusInternalServerError, err.Error())
		d.db = nil
		return
	}

	defer d.db.Close()
}

func (d DB) close() {
	if d.db != nil {
		d.db.Close()
	}
}

func (d DB) Check() bool {
	return d.db != nil
}

func (d DB) AddCounter(key string, value int64) {
	panic("unimplemented")
}

func (d DB) AddGauge(key string, value float64) {
	panic("unimplemented")
}

func (d DB) GetAll() map[string]models.Metric {
	panic("unimplemented")
}

func (d DB) GetCounter(key string) (int64, bool) {
	panic("unimplemented")
}

func (d DB) GetGauge(key string) (float64, bool) {
	panic("unimplemented")
}

func NewDB(host string) keeper.Storager {
	db := DB{}

	db.connect(host)

	defer db.close()
	return db
}
