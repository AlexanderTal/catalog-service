package rcpostgres

import (
	"database/sql"
	"errors"

	"github.com/AlexanderTal/catalog-service/internal/app/entity"
)

func RowsAffected(res sql.Result) int64 {
	ra, _ := res.RowsAffected()
	return ra
}

func UpdateErr(res sql.Result, err error) error {
	if err == nil && RowsAffected(res) == 0 {
		return entity.ErrNotFound
	}
	if errors.Is(err, sql.ErrNoRows) {
		return entity.ErrNotFound
	}
	return err
}

func DeleteErr(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}
