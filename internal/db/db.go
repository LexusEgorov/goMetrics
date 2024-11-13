package db

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/LexusEgorov/goMetrics/internal/dohsimpson"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	db *sql.DB
}

func (d *DB) Connect(host string) {
	ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, `postgres`, `root`, `metrics`)

	var err error
	d.db, err = sql.Open("pgx", ps)

	if err != nil {
		dohsimpson.NewDoh(http.StatusInternalServerError, err.Error())
	}

	defer d.db.Close()
}

func (d DB) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

func (d DB) Check() bool {
	return bool(d.db != nil)
}

func NewDB(host string) DB {
	db := DB{}

	db.Connect(host)
	return db
}
