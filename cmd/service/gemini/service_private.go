package gemini

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jeraldyik/crypto_dca_go/cmd/config"
	"github.com/jeraldyik/crypto_dca_go/cmd/util"
	"github.com/jeraldyik/crypto_dca_go/internal/logger"
)

// New Order
func (api *Api) newOrder(ticker, price, amount string) (*Order, error) {
	location := "gemini.newOrder"
	now := config.GetTime().NowTimestamp(nil)
	quoteCurrency := AppendTickerWithQuoteCurrency(ticker)
	params := map[string]any{
		"request":         NewOrderURI,
		"nonce":           now,
		"client_order_id": fmt.Sprintf("%v_%v", now, quoteCurrency),
		"symbol":          quoteCurrency,
		"price":           price,
		"amount":          amount,
		"side":            "buy",
		"type":            "exchange limit",
	}

	logger.Info(location, "params:%+v", params)

	order := &Order{}

	body, err := api.request(http.MethodPost, NewOrderURI, params)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, order); err != nil {
		return nil, err
	}

	logger.Info(location, "order: %v", util.SafeJsonDump(order))

	return order, nil
}

// Order Status
func (api *Api) orderStatus(orderID string) (*Order, error) {
	location := "gemini.orderStatus"
	params := map[string]any{
		"request":  OrderStatusURI,
		"nonce":    config.GetTime().NowTimestamp(nil),
		"order_id": orderID,
	}

	logger.Info(location, "params:%+v", params)

	order := &Order{}

	body, err := api.request(http.MethodPost, OrderStatusURI, params)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, order); err != nil {
		return nil, err
	}

	logger.Info(location, "order: %v", util.SafeJsonDump(order))

	return order, nil
}

// Cancel Order
func (api *Api) cancelOrder(orderID string) (*Order, error) {
	location := "gemini.cancelOrder"
	params := map[string]any{
		"request":  CancelOrderURI,
		"nonce":    config.GetTime().NowTimestamp(nil),
		"order_id": orderID,
	}

	logger.Info(location, "params:%+v", params)

	order := &Order{}

	body, err := api.request(http.MethodPost, CancelOrderURI, params)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, order); err != nil {
		return nil, err
	}

	logger.Info(location, "order: %v", util.SafeJsonDump(order))

	return order, nil
}
