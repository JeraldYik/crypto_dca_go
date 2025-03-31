package db

import "time"

type Order struct {
	Ticker            string    `json:"ticker"`
	CreatedForDay     time.Time `json:"createdForDay"`
	FiatDepositInSGD  float64   `json:"fiatDepositInSgd"`  // legacy issue: could also be in other fiat curreny (i.e. USD)
	PricePerCoinInSGD float64   `json:"pricePerCoinInSgd"` // legacy issue: could also be in other fiat curreny (i.e. USD)
	CoinAmount        float64   `json:"coinAmount"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func (Order) TableName() string {
	return "Orders"
}
