package db

import (
	"errors"

	"github.com/jeraldyik/crypto_dca_go/internal/logger"
)

func (o *OrderDB) BulkInsert(rows []*Order) error {
	location := "db.BulkInsert"
	result := o.db.Create(rows)
	if result.Error != nil {
		logger.Error(location, "Failed to insert rows", result.Error)
		return result.Error
	} else if result.RowsAffected != int64(len(rows)) {
		err := errors.New("db_insert_mismatched_rows_count")
		logger.Error(location, "Failed to insert correct number of rows. got = %v, expected = %v", err, result.RowsAffected, len(rows))
		return err
	}

	logger.Info(location, "Successfully inserted %v rows", len(rows))
	return nil
}
