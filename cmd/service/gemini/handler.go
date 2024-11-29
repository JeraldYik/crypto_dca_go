package gemini

import (
	"log"

	"github.com/jeraldyik/crypto_dca_go/cmd/util"
)

func (api *Api) GetQuoteIncrementAndTickSize(ticker string) (int, int, error) {
	tickerData, err := api.tickerDetails(ticker)
	if err != nil {
		log.Printf("[gemini.GetTickSize] err: %+v\n", err)
		return 0, 0, err
	}
	return util.NumDecimalPlaces(tickerData.QuoteIncrement), util.NumDecimalPlaces(tickerData.TickSize), nil
}

func (api *Api) GetTickerBestBidPrice(ticker string) (float64, error) {
	tickerActivity, err := api.tickerV2(ticker)
	if err != nil {
		log.Printf("[gemini.GetTickerBestBidPrice] err: %+v\n", err)
		return 0, err
	}
	return tickerActivity.Bid, nil
}

func (api *Api) CreateOrder(ticker string, bestBid float64, quoteIncrement, tickSize int) (*Order, error) {
	orderPriceStr, orderAmountStr := formCreateOrderReq(ticker, bestBid, quoteIncrement, tickSize)
	order, err := api.newOrder(ticker, orderPriceStr, orderAmountStr)
	if err != nil {
		log.Printf("[gemini.CreateOrder] err: %+v\n", err)
		return nil, err
	}
	return order, nil
}

func (api *Api) GetOrderStatus(orderID string) (*Order, error) {
	order, err := api.orderStatus(orderID)
	if err != nil {
		log.Printf("[gemini.GetOrderStatus] err: %+v\n", err)
		return nil, err
	}
	return order, nil
}

func (api *Api) CancelOrder(orderID string) (*Order, error) {
	order, err := api.cancelOrder(orderID)
	if err != nil {
		log.Printf("[gemini.cancelOrder] err: %+v\n", err)
		return nil, err
	}
	return order, nil
}
