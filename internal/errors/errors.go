package errors

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func IsServerRetriable(err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == "28P01" {
			return true
		}
	}

	return false
}

func IsClientRetriable(code int) bool {
	return code == 0
}
