package cmd

import (
	"context"

	"github.com/jeraldyik/crypto_dca_go/internal/logger"
)

func Run(ctx context.Context) error {
	location := "cmd.Run"
	logger.Info(location, "Running main script...")

	postOrderDetails := handleOrder(ctx)
	logger.Info(location, "postOrderDetails: %v", postOrderDetails)

	// update google sheets cells
	if err := batchUpdate(postOrderDetails); err != nil {
		logger.Error(location, "Batch update google sheets", err)
		return err
	}
	logger.Info(location, "Batch update google sheets successful")

	// insert into db
	if err := bulkInsertIntoDB(postOrderDetails); err != nil {
		logger.Error(location, "Batch insert into db", err)
		return err
	}
	logger.Info(location, "Batch insert into db successful")

	logger.Info(location, "Successfully completed. Tearing down...")

	return nil
}
