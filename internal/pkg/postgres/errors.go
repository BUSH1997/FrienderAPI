package postgres

import (
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

const (
	UniqueViolationError = 1
)

func ProcessError(err error) int {
	return UniqueViolationError
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return 0
	}

	if pgErr.Code == pgerrcode.UniqueViolation {
		return UniqueViolationError
	}

	return 0
}
