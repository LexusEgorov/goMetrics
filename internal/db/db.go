package db

import (
	"database/sql"
	"net/http"

	"github.com/LexusEgorov/goMetrics/internal/dohsimpson"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	db *sql.DB
}

func (d *DB) Connect(host string) {
	var err error
	d.db, err = sql.Open("pgx", host)

	if err != nil {
		dohsimpson.NewDoh(http.StatusInternalServerError, err.Error())
		return
	}

	err = d.db.Ping()

	if err != nil {
		dohsimpson.NewDoh(http.StatusInternalServerError, err.Error())
		d.db = nil
		return
	}

	defer d.db.Close()
}

func (d DB) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

func (d DB) Check() bool {
	return d.db != nil
	// if d.db == nil {
	// 	return false
	// }

	// err := d.db.Ping()

	// return bool(err == nil)
}

func NewDB(host string) DB {
	db := DB{}

	db.Connect(host)
	return db
}
