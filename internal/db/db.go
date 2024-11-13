package db

import (
	"database/sql"
	"fmt"
)

type DB struct {
	db *sql.DB
}

func (d *DB) Connect(host string) {
	ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, `root`, `root`, `metrics`)

	var err error
	d.db, err = sql.Open("pgx", ps)

	if err != nil {
		panic(err)
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

func NewDB(host string) *DB {
	db := &DB{}

	db.Connect(host)
	return db
}
