package db

import (
	"fmt"

	"github.com/jeraldyik/crypto_dca_go/cmd/config"
	"github.com/jeraldyik/crypto_dca_go/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:generate mockgen -source=/cmd/service/db/main.go -destination=/mocks/mock_OrderRepository.go -package=mocks
type OrderRepository interface {
	BulkInsert(rows []*Order) error
	GetDB() *gorm.DB
}

type OrderDB struct {
	db *gorm.DB
}

var orderDB OrderRepository

func (o *OrderDB) GetDB() *gorm.DB {
	return o.db
}

func MustInit() {
	location := "db.MustInit"
	c := config.Get().Db
	conn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", c.Host, c.Username, c.Password, c.Name)

	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		logger.Panic(location, "Failed to initialise database, err: %+v", err)
	}

	Set(&OrderDB{db: db})
}

func Close() error {
	location := "db.Close"
	sqlDB, err := orderDB.GetDB().DB()
	if err != nil {
		logger.Error(location, "Error in getting underlying db", err)
		return err
	}
	sqlDB.Close()
	return nil
}

func Get() OrderRepository {
	return orderDB
}

func Set(o OrderRepository) {
	orderDB = o
}
