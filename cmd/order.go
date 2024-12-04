package cmd

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/jeraldyik/crypto_dca_go/cmd/config"
	"github.com/jeraldyik/crypto_dca_go/cmd/service/gemini"
	"github.com/jeraldyik/crypto_dca_go/cmd/util"
	"github.com/jeraldyik/crypto_dca_go/internal/logger"
)

type PostOrder struct {
	ActualFiatDeposit float64
	AvgExecutionPrice float64
	ExecutedAmount    float64
}

type PostOrderDetails struct {
	m  *treemap.Map
	mu sync.Mutex
}

// Entry point for creating & fulfilling orders
func handleOrder(ctx context.Context) *treemap.Map {
	location := "cmd.handlerOrder"
	c := config.Get()
	postOrderDetails := &PostOrderDetails{
		m: treemap.NewWithStringComparator(),
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(c.CryptoTickers))
	for ticker := range c.CryptoTickers {
		go func(ctx context.Context, wg *sync.WaitGroup, ticker string, postOrderMap *PostOrderDetails) {
			defer wg.Done()
			util.RecoverAndGraceFullyExit()
			// Check switch
			if c.OrderMetadata.DailyFiatAmount[ticker] <= 0 {
				logger.Warn(location, "Purchase for ticker '%s' is turned off", ticker)
				return
			}

			// Sandbox environment guard check
			if c.IsSandboxEnv {
				addToPostOrderDetails(postOrderDetails, ticker, nil)
				return
			}

			handlerCexApiCalls(ctx, ticker, postOrderMap)
		}(ctx, wg, ticker, postOrderDetails)
	}
	wg.Wait()

	return postOrderDetails.m
}

// Entry point for goroutine - Level 1
func handlerCexApiCalls(ctx context.Context, ticker string, postOrderDetails *PostOrderDetails) {
	location := "handler.handlerCexApiCalls"
	geminiClient := gemini.GetClient()

	// Get Symbol details
	results, err := gemini.RetryWrapper(ctx, fmt.Sprintf("GetQuoteIncrementAndTickSize - %v", ticker), geminiClient.GetQuoteIncrementAndTickSize, ticker)
	if err != nil {
		logger.Error(location, "[handler.handlerCexApiCalls] Error getting symbol details", err)
		return
	}
	quoteIncrement, tickSize := int(results[0].Int()), int(results[1].Int())

	orderOpenThenCancelWindowCounter := 0
	for orderOpenThenCancelWindowCounter < config.OrderOpenThenCancelWindowCount {
		orderOpenThenCancelWindowCounter++

		order, err := handlerCexApiCallsOrderOpenThenCancel(ctx, ticker, quoteIncrement, tickSize)
		if err != nil {
			continue
		}
		if order != nil {
			addToPostOrderDetails(postOrderDetails, ticker, order)
			return
		}
	}

	logger.Warn(location, "Ticker '%s' failed to have a fulfilled order", ticker)
}

// Level 2
func handlerCexApiCallsOrderOpenThenCancel(ctx context.Context, ticker string, quoteIncrement, tickSize int) (*gemini.Order, error) {
	location := "handler.handlerCexApiCallsOrderOpenThenCancel"
	geminiClient := gemini.GetClient()

	// Get ticker best bid price
	results, err := gemini.RetryWrapper(ctx, fmt.Sprintf("GetTickerBestBidPrice - %v", ticker), geminiClient.GetTickerBestBidPrice, ticker)
	if err != nil {
		logger.Error(location, "'%s' Error getting best bid price", err, ticker)
		return nil, err
	}
	bestBid := results[0].Float()

	// TODO: to monitor on situation on http error and no order created
	// Create order - not retrying to prevent side effects
	order, err := geminiClient.CreateOrder(ticker, bestBid, quoteIncrement, tickSize)
	if err != nil {
		logger.Error(location, "'%s' Error creating order", err, ticker)
		// Search order in case of order already exists
		results, err = gemini.RetryWrapper(ctx, fmt.Sprintf("MatchActiveOrders - %v, ticker), geminiClient.MatchActiveOrders", ticker), geminiClient.MatchActiveOrders, ticker)
		if err != nil {
			logger.Error(location, "'%s' Error getting and matching active orders", err, ticker)
			return nil, err
		}
		order = results[0].Interface().(*gemini.Order)
	}

	// If order is cancelled, re-create order - not retrying to prevent side effects
	recreatingOrderCount := 0
	for order.IsCancelled && recreatingOrderCount < gemini.MaxRetryCount {
		logger.Warn(location, "'%s' Order is cancelled, re-creating order", ticker)
		recreatingOrderCount++

		order, err = geminiClient.CreateOrder(ticker, bestBid, quoteIncrement, tickSize)
		if err != nil {
			logger.Error(location, "Error creating order", err)
			// Search order in case of order already exists
			results, err = gemini.RetryWrapper(ctx, fmt.Sprintf("MatchActiveOrders - %v, ticker), geminiClient.MatchActiveOrders", ticker), geminiClient.MatchActiveOrders, ticker)
			if err != nil {
				logger.Error(location, "'%s' Error getting and matching active orders", err, ticker)
				return nil, err
			}
			order = results[0].Interface().(*gemini.Order)
		}
	}

	// If order is somehow still cancelled after retrying - return error
	if order.IsCancelled {
		err := errors.New("order is cancelled")
		logger.Error(location, "'%s' Order is still cancelled after retrying", err, ticker)
		return nil, err
	}

	// If order fulfilled - return order
	if !order.IsLive {
		logger.Info(location, "'%s' Order is fulfilled", ticker)
		return order, nil
	}

	// Check if order is fulfilled - query every minute for an hour
	// Make sure that order is not cancelled - if cancelled, return
	orderOpenQueryStatusWindowCounter := 0
	for orderOpenQueryStatusWindowCounter < config.OrderOpenQueryStatusWindowCount {
		orderOpenQueryStatusWindowCounter++
		if !util.IsTestFlow(ctx) {
			logger.Info(location, "'%s' waiting for 1 min", ticker)
			time.Sleep(1 * time.Minute)
		}

		order, isCancelled, err := handlerCexApiCallsOrderOpenQueryStatus(ctx, ticker, order)
		if err != nil {
			break // to cancel order
		}
		if isCancelled {
			return nil, errors.New("order is cancelled")
		}
		if order != nil {
			return order, nil
		}
	}

	// Cancel order here, retry creating new order in the next iteration of the loop
	results, err = gemini.RetryWrapper(ctx, fmt.Sprintf("CancelOrder - %v", ticker), geminiClient.CancelOrder, order.OrderID)
	if err != nil || !results[0].Interface().(*gemini.Order).IsCancelled {
		logger.Error(location, "'%s' Failed to cancel order: %+v", err, ticker)
		return nil, err
	}

	// Order is not filled and successfully cancelled
	logger.Warn(location, "'%s' Order is not filled and successfully cancelled", ticker)
	return nil, nil
}

// Level 3
//
// bool: order is cancelled
func handlerCexApiCallsOrderOpenQueryStatus(ctx context.Context, ticker string, order *gemini.Order) (*gemini.Order, bool, error) {
	location := "handler.handlerCexApiCallsOrderOpenQueryStatus"
	geminiClient := gemini.GetClient()

	results, err := gemini.RetryWrapper(ctx, fmt.Sprintf("GetOrderStatus - %v", ticker), geminiClient.GetOrderStatus, order.OrderID)
	if err != nil {
		logger.Error(location, "'%s' Get order status failed", err, ticker)
		return nil, false, err
	}
	queryOrder := results[0].Interface().(*gemini.Order)

	// If order is cancelled - return empty
	if queryOrder.IsCancelled {
		logger.Warn(location, "'%s' Order is cancelled", ticker)
		return nil, true, nil
	}

	// If order fulfilled - return order
	if !queryOrder.IsLive {
		logger.Info(location, "'%s' Order is fulfilled", ticker)
		return queryOrder, false, nil
	}

	// Order is not fulfilled - to continue querying
	logger.Warn(location, "'%s' Order is not fulfilled yet", ticker)
	return nil, false, nil
}

func addToPostOrderDetails(postOrderDetails *PostOrderDetails, ticker string, order *gemini.Order) {
	postOrderDetails.mu.Lock()
	if order != nil {
		// Prod
		postOrderDetails.m.Put(ticker, formPostOrderData(order))
	} else {
		// Sandbox
		c := config.Get()
		postOrderDetails.m.Put(ticker, sandboxPostOrderData(c.OrderMetadata.DailyFiatAmount[ticker]))
	}
	postOrderDetails.mu.Unlock()
}

func formPostOrderData(order *gemini.Order) PostOrder {
	return PostOrder{
		ActualFiatDeposit: order.AvgExecutionPrice * order.ExecutedAmount * (1 + gemini.MakerTradingFee),
		AvgExecutionPrice: order.AvgExecutionPrice,
		ExecutedAmount:    order.ExecutedAmount,
	}
}

func sandboxPostOrderData(dailyFiatAmount float64) PostOrder {
	return PostOrder{
		ActualFiatDeposit: dailyFiatAmount * (1 + gemini.MakerTradingFee),
		AvgExecutionPrice: 1000,
		ExecutedAmount:    1,
	}
}
