package cmd

import (
	"context"
	"log"
)

func Run(ctx context.Context) error {
	log.Println("[cmd.Run] Running main script...")

	postOrderDetails := handleOrder(ctx)
	log.Printf("[cmd.Run] postOrderDetails: %+v\n", postOrderDetails)

	// update google sheets cells
	if err := batchUpdate(postOrderDetails); err != nil {
		log.Printf("[cmd.Run] Batch update google sheets err: %+v\n", err)
		return err
	}
	log.Println("[cmd.Run] Batch update google sheets successful")

	// insert into db
	if err := bulkInsertIntoDB(postOrderDetails); err != nil {
		log.Printf("[cmd.Run] Batch insert into db err: %+v\n", err)
		return err
	}
	log.Println("[cmd.Run] Batch insert into db successful")

	log.Println("[cmd.Run] Successfully completed. Tearing down...")

	return nil
}
