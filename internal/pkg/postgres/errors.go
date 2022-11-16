package postgres

import (
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

const (
	UniqueViolationError = 1
)

func ProcessError(err error) int {
	pgErr, ok := err.(*pq.Error)
	if !ok {
		return 0
	}

	if pgErr.Code == pgerrcode.UniqueViolation {
		return UniqueViolationError
	}

	return 0
}
