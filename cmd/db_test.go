package cmd

import (
	"testing"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/jeraldyik/crypto_dca_go/cmd/config"
	"github.com/jeraldyik/crypto_dca_go/cmd/service/db"
	"github.com/stretchr/testify/assert"
)

func Test_formRows(t *testing.T) {
	config.TestInit(nil, nil)

	t.Run("ok", func(t *testing.T) {
		postOrders := treemap.NewWithStringComparator()
		postOrders.Put("BTC", PostOrder{
			ActualFiatDeposit: 1.002,
			AvgExecutionPrice: 1000,
			ExecutedAmount:    1,
		})
		postOrders.Put("ETH", PostOrder{
			ActualFiatDeposit: 2.004,
			AvgExecutionPrice: 1000,
			ExecutedAmount:    1,
		})
		got := formRows(postOrders)
		assert.Equal(t, []*db.Order{
			{
				Ticker:            "btcsgd",
				CreatedForDay:     config.GetTime().GetTodayDate(),
				FiatDepositInSGD:  1.002,
				PricePerCoinInSGD: 1000,
				CoinAmount:        1,
				CreatedAt:         config.GetTime().Now(),
				UpdatedAt:         config.GetTime().Now(),
			},
			{
				Ticker:            "ethsgd",
				CreatedForDay:     config.GetTime().GetTodayDate(),
				FiatDepositInSGD:  2.004,
				PricePerCoinInSGD: 1000,
				CoinAmount:        1,
				CreatedAt:         config.GetTime().Now(),
				UpdatedAt:         config.GetTime().Now(),
			},
		}, got)
	})
}
