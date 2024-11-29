package db

import (
	"errors"
	"log"
)

func (o *OrderDB) BulkInsert(rows []*Order) error {
	result := o.db.Create(rows)
	if result.Error != nil {
		log.Printf("[db.BulkInsert] Failed to insert rows, err: %+v\n", result.Error)
		return result.Error
	} else if result.RowsAffected != int64(len(rows)) {
		log.Printf("[db.BulkInsert] Failed to insert correct number of rows. got = %v, expected = %v\n", result.RowsAffected, len(rows))
		return errors.New("db_insert_mismatched_rows_count")
	}

	log.Printf("[db.BulkInsert] Successfully inserted %v rows\n", len(rows))
	return nil
}
