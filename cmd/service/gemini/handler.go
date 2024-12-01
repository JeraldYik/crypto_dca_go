package gemini

import (
	"github.com/jeraldyik/crypto_dca_go/cmd/util"
	"github.com/jeraldyik/crypto_dca_go/internal/logger"
)

func (api *Api) GetQuoteIncrementAndTickSize(ticker string) (int, int, error) {
	location := "gemini.GetTickSize"
	tickerData, err := api.tickerDetails(ticker)
	if err != nil {
		logger.Error(location, "", err)
		return 0, 0, err
	}
	return util.NumDecimalPlaces(tickerData.QuoteIncrement), util.NumDecimalPlaces(tickerData.TickSize), nil
}

func (api *Api) GetTickerBestBidPrice(ticker string) (float64, error) {
	location := "gemini.GetTickerBestBidPrice"
	tickerActivity, err := api.tickerV2(ticker)
	if err != nil {
		logger.Error(location, "", err)
		return 0, err
	}
	return tickerActivity.Bid, nil
}

func (api *Api) CreateOrder(ticker string, bestBid float64, quoteIncrement, tickSize int) (*Order, error) {
	location := "gemini.CreateOrder"
	orderPriceStr, orderAmountStr := formCreateOrderReq(ticker, bestBid, quoteIncrement, tickSize)
	order, err := api.newOrder(ticker, orderPriceStr, orderAmountStr)
	if err != nil {
		logger.Error(location, "", err)
		return nil, err
	}
	return order, nil
}

func (api *Api) GetOrderStatus(orderID string) (*Order, error) {
	location := "gemini.GetOrderStatus"
	order, err := api.orderStatus(orderID)
	if err != nil {
		logger.Error(location, "", err)
		return nil, err
	}
	return order, nil
}

func (api *Api) CancelOrder(orderID string) (*Order, error) {
	location := "gemini.cancelOrder"
	order, err := api.cancelOrder(orderID)
	if err != nil {
		logger.Error(location, "", err)
		return nil, err
	}
	return order, nil
}
