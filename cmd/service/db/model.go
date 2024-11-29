package db

import "time"

type Order struct {
	Ticker            string    `gorm:"column:ticker"`
	CreatedForDay     time.Time `gorm:"column:createdForDay"`
	FiatDepositInSGD  float64   `gorm:"column:fiatDepositInSgd"`  // legacy issue: could also be in other fiat curreny (i.e. USD)
	PricePerCoinInSGD float64   `gorm:"column:pricePerCoinInSgd"` // legacy issue: could also be in other fiat curreny (i.e. USD)
	CoinAmount        float64   `gorm:"column:coinAmount"`
	CreatedAt         time.Time `gorm:"column:createdAt"`
	UpdatedAt         time.Time `gorm:"column:updatedAt"`
}

func (Order) TableName() string {
	return "Orders"
}
