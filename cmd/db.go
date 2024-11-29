package cmd

import (
	"strings"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/jeraldyik/crypto_dca_go/cmd/config"
	"github.com/jeraldyik/crypto_dca_go/cmd/service/db"
	"github.com/jeraldyik/crypto_dca_go/cmd/service/gemini"
)

func formRows(postOrders *treemap.Map) []*db.Order {
	orders := make([]*db.Order, postOrders.Size())
	i := 0
	it := postOrders.Iterator()
	for it.Next() {
		ticker, postOrder := it.Key().(string), it.Value().(PostOrder)
		orders[i] = &db.Order{
			Ticker:            strings.ToLower(gemini.AppendTickerWithQuoteCurrency(ticker)),
			CreatedForDay:     config.GetTime().GetTodayDate(),
			FiatDepositInSGD:  postOrder.ActualFiatDeposit,
			PricePerCoinInSGD: postOrder.AvgExecutionPrice,
			CoinAmount:        postOrder.ExecutedAmount,
		}
		i++
	}
	return orders
}

func bulkInsertIntoDB(postOrders *treemap.Map) error {
	orders := formRows(postOrders)
	if err := db.Get().BulkInsert(orders); err != nil {
		return err
	}
	return nil
}
