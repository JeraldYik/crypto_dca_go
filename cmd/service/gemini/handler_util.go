package gemini

import (
	"github.com/jeraldyik/crypto_dca_go/cmd/config"
	"github.com/jeraldyik/crypto_dca_go/cmd/util"
	"github.com/shopspring/decimal"
)

// return orderPriceStr, orderAmountStr
func formCreateOrderReq(ticker string, bestBid float64, quoteIncrement, tickSize int) (string, string) {
	orderMetadata := config.Get().OrderMetadata
	orderPrice := bestBid * orderMetadata.OrderPriceToBidPriceRatio
	orderAmount := decimal.NewFromFloat(orderMetadata.DailyFiatAmount[ticker]).Div(decimal.NewFromFloat(orderPrice))
	orderPriceStr := util.ConvertFloatToPrecString(orderPrice, quoteIncrement)
	orderAmountStr := util.ConvertFloatToPrecString(orderAmount, tickSize)

	return orderPriceStr, orderAmountStr
}
