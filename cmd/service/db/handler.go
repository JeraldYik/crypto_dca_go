package db

import (
	"errors"

	"github.com/jeraldyik/crypto_dca_go/internal/logger"
)

func (o *OrderDB) BulkInsert(rows []*Order) error {
	location := "db.BulkInsert"
	/*
		result := o.db.Create(rows)
		if result.Error != nil {
			logger.Error(location, "Failed to insert rows", result.Error)
			return result.Error
		} else if result.RowsAffected != int64(len(rows)) {
			err := errors.New("db_insert_mismatched_rows_count")
			logger.Error(location, "Failed to insert correct number of rows. got = %v, expected = %v", err, result.RowsAffected, len(rows))
			return err
		}
	*/
	_, num_rows, err := o.db.From("Orders").Insert(rows, false, "", "minimal", "exact").Execute()
	if err != nil {
		logger.Error(location, "Failed to insert rows: %v", err)
		return err
	} else if num_rows != int64(len(rows)) {
		err := errors.New("db_insert_mismatched_rows_count")
		logger.Error(location, "Failed to insert correct number of rows. got = %v, expected = %v", err, num_rows, len(rows))
		return err
	}

	logger.Info(location, "Successfully inserted %v rows", len(rows))
	return nil
}
