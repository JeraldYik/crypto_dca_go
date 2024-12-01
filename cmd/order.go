package cmd

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/jeraldyik/crypto_dca_go/cmd/config"
	"github.com/jeraldyik/crypto_dca_go/cmd/service/gemini"
	"github.com/jeraldyik/crypto_dca_go/cmd/util"
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
				log.Printf("[cmd.handlerOrder] Purchase for ticker '%s' is turned off\n", ticker)
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
	geminiClient := gemini.GetClient()

	// Get Symbol details
	results, err := gemini.RetryWrapper(ctx, "GetQuoteIncrementAndTickSize", geminiClient.GetQuoteIncrementAndTickSize, ticker)
	if err != nil {
		log.Printf("[handler.handlerCexApiCalls] Error getting symbol details, err: %+v\n", err)
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

	log.Printf("[handler.handlerCexApiCalls] Ticker '%s' failed to have a fulfilled order\n", ticker)
}

// Level 2
func handlerCexApiCallsOrderOpenThenCancel(ctx context.Context, ticker string, quoteIncrement, tickSize int) (*gemini.Order, error) {
	geminiClient := gemini.GetClient()

	// Get ticker best bid price
	results, err := gemini.RetryWrapper(ctx, "GetTickerBestBidPrice", geminiClient.GetTickerBestBidPrice, ticker)
	if err != nil {
		log.Printf("[handler.handlerCexApiCallsOrderOpenThenCancel] '%s' Error getting best bid price, err: %+v\n", ticker, err)
		return nil, err
	}
	bestBid := results[0].Float()

	// Create order
	results, err = gemini.RetryWrapper(ctx, "CreateOrder", geminiClient.CreateOrder, ticker, bestBid, quoteIncrement, tickSize)
	if err != nil {
		log.Printf("[handler.handlerCexApiCallsOrderOpenThenCancel] '%s' Error creating order, err: %+v\n", ticker, err)
		return nil, err
	}
	order := results[0].Interface().(*gemini.Order)

	// If order is cancelled, re-create order
	recreatingOrderCount := 0
	for order.IsCancelled && recreatingOrderCount < gemini.MaxRetryCount {
		log.Printf("[handler.handlerCexApiCallsOrderOpenThenCancel] '%s' Order is cancelled, re-creating order\n", ticker)
		recreatingOrderCount++

		results, err = gemini.RetryWrapper(ctx, "CreateOrder", geminiClient.CreateOrder, ticker, bestBid, quoteIncrement, tickSize)
		if err != nil {
			log.Printf("[handler.handlerCexApiCallsOrderOpenThenCancel] Error creating order, err: %+v\n", err)
			return nil, err
		}
		order = results[0].Interface().(*gemini.Order)
	}

	// If order is somehow still cancelled after retrying - return error
	if order.IsCancelled {
		log.Printf("[handler.handlerCexApiCallsOrderOpenThenCancel] '%s' Order is still cancelled after retrying\n", ticker)
		return nil, errors.New("order is cancelled")
	}

	// If order fulfilled - return order
	if !order.IsLive {
		log.Printf("[handler.handlerCexApiCallsOrderOpenThenCancel] '%s' Order is fulfilled\n", ticker)
		return order, nil
	}

	// Check if order is fulfilled - query every minute for an hour
	// Make sure that order is not cancelled - if cancelled, return
	orderOpenQueryStatusWindowCounter := 0
	for orderOpenQueryStatusWindowCounter < config.OrderOpenQueryStatusWindowCount {
		orderOpenQueryStatusWindowCounter++
		if !util.IsTestFlow(ctx) {
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
	results, err = gemini.RetryWrapper(ctx, "CancelOrder", geminiClient.CancelOrder, order.OrderID)
	if err != nil || !results[0].Interface().(*gemini.Order).IsCancelled {
		log.Printf("[handler.handlerCexApiCallsOrderOpenThenCancel] '%s' Failed to cancel order: %+v\n", ticker, err)
		return nil, err
	}

	// Order is not filled and successfully cancelled
	log.Printf("[handler.handlerCexApiCallsOrderOpenThenCancel] '%s' Order is not filled and successfully cancelled\n", ticker)
	return nil, nil
}

// Level 3
//
// bool: order is cancelled
func handlerCexApiCallsOrderOpenQueryStatus(ctx context.Context, ticker string, order *gemini.Order) (*gemini.Order, bool, error) {
	geminiClient := gemini.GetClient()

	results, err := gemini.RetryWrapper(ctx, "GetOrderStatus", geminiClient.GetOrderStatus, order.OrderID)
	if err != nil {
		log.Printf("[handler.handlerCexApiCallsOrderOpenQueryStatus] '%s' Get order status failed, err: %+v\n", ticker, err)
		return nil, false, err
	}
	queryOrder := results[0].Interface().(*gemini.Order)

	// If order is cancelled - return empty
	if queryOrder.IsCancelled {
		log.Printf("[handler.handlerCexApiCallsOrderOpenQueryStatus] '%s' Order is cancelled\n", ticker)
		return nil, true, nil
	}

	// If order fulfilled - return order
	if !queryOrder.IsLive {
		log.Printf("[handler.handlerCexApiCallsOrderOpenQueryStatus] '%s' Order is fulfilled\n", ticker)
		return queryOrder, false, nil
	}

	// Order is not fulfilled - to continue querying
	log.Printf("[handler.handlerCexApiCallsOrderOpenQueryStatus] '%s' Order is not fulfilled yet\n", ticker)
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
