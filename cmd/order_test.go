package cmd

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/jarcoal/httpmock"
	"github.com/jeraldyik/crypto_dca_go/cmd/config"
	"github.com/jeraldyik/crypto_dca_go/cmd/service/gemini"
	"github.com/jeraldyik/crypto_dca_go/cmd/util"
	"github.com/stretchr/testify/assert"
)

func Test_handleOrder(t *testing.T) {
	ctx := util.TestContext()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	config.TestInit(nil, nil)
	gemini.MustInitClient()

	t.Run("ok_prod", func(t *testing.T) {
		config.TestInit(&config.ConfigUpdateable{
			IsSandboxEnv: util.PtrOf(false),
		}, nil)
		defer httpmock.Reset()
		responder := httpmock.NewStringResponder(http.StatusOK, `{
			"tick_size": 1E-8,
			"quote_increment": 0.01
		}`)
		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerDetailsURI, "btcsgd"), responder)
		responder = httpmock.NewStringResponder(http.StatusOK, `{
			"tick_size": 1E-6,
			"quote_increment": 0.01
		}`)
		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerDetailsURI, "ethsgd"), responder)

		responder = httpmock.NewStringResponder(http.StatusOK, `{
			"bid": "9345.70"
		}`)
		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerV2URI, "btcsgd"), responder)
		responder = httpmock.NewStringResponder(http.StatusOK, `{
			"bid": "9345.70"
		}`)
		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerV2URI, "ethsgd"), responder)

		responder = httpmock.NewStringResponder(http.StatusOK, `{
				"order_id": "106817811", 
				"avg_execution_price": "3632.8508430064554",
				"is_live": false, 
				"is_cancelled": false, 
				"executed_amount": "3.7567928949",
				"client_order_id": "20190110-4738721"
		}`)
		httpmock.RegisterResponder(http.MethodPost, gemini.NewOrderURI, responder)

		postOrderDetails := handleOrder(ctx)
		time.Sleep(time.Duration(len(config.Get().CryptoTickers)) * time.Second) // wait for goroutines to finish

		postOrders := treemap.NewWithStringComparator()
		postOrders.Put("BTC", PostOrder{
			ActualFiatDeposit: 13675.163971708602,
			AvgExecutionPrice: 3632.8508430064553,
			ExecutedAmount:    3.7567928949,
		})
		postOrders.Put("ETH", PostOrder{
			ActualFiatDeposit: 13675.163971708602,
			AvgExecutionPrice: 3632.8508430064553,
			ExecutedAmount:    3.7567928949,
		})
		assert.Equal(t, util.SafeJsonDump(postOrders), util.SafeJsonDump(postOrderDetails))
	})

	t.Run("ok_sandbox", func(t *testing.T) {
		config.TestInit(nil, nil)
		postOrderDetails := handleOrder(ctx)
		time.Sleep(time.Duration(len(config.Get().CryptoTickers)) * time.Second) // wait for goroutines to finish

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
		assert.Equal(t, util.SafeJsonDump(postOrders), util.SafeJsonDump(postOrderDetails))
	})

	t.Run("ok_sandbox_ignore_ETH", func(t *testing.T) {
		config.TestInit(&config.ConfigUpdateable{
			DailyFiatAmount: map[string]float64{
				"BTC": 1,
				"ETH": 0,
			},
		}, nil)
		postOrderDetails := handleOrder(ctx)
		time.Sleep(time.Duration(len(config.Get().CryptoTickers)) * time.Second) // wait for goroutines to finish

		postOrders := treemap.NewWithStringComparator()
		postOrders.Put("BTC", PostOrder{
			ActualFiatDeposit: 1.002,
			AvgExecutionPrice: 1000,
			ExecutedAmount:    1,
		})
		assert.Equal(t, util.SafeJsonDump(postOrders), util.SafeJsonDump(postOrderDetails))
	})
}

func Test_handlerCexApiCalls(t *testing.T) {
	ctx := util.TestContext()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	config.TestInit(nil, nil)
	gemini.MustInitClient()

	t.Run("ok", func(t *testing.T) {
		defer httpmock.Reset()
		responder := httpmock.NewStringResponder(http.StatusOK, `{
			"tick_size": 1E-8,
			"quote_increment": 0.01
		}`)
		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerDetailsURI, "btcsgd"), responder)

		responder = httpmock.NewStringResponder(http.StatusOK, `{
			"bid": "9345.70"
		}`)
		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerV2URI, "btcsgd"), responder)

		responder = httpmock.NewStringResponder(http.StatusOK, `{
				"order_id": "106817811", 
				"avg_execution_price": "3632.8508430064554",
				"is_live": false, 
				"is_cancelled": false, 
				"executed_amount": "3.7567928949",
				"client_order_id": "20190110-4738721"
		}`)
		httpmock.RegisterResponder(http.MethodPost, gemini.NewOrderURI, responder)

		postOrderMap := &PostOrderDetails{
			m: treemap.NewWithStringComparator(),
		}
		postOrders := treemap.NewWithStringComparator()
		postOrders.Put("BTC", PostOrder{
			ActualFiatDeposit: 13675.163971708602,
			AvgExecutionPrice: 3632.8508430064553,
			ExecutedAmount:    3.7567928949,
		})
		handlerCexApiCalls(ctx, "BTC", postOrderMap)
		assert.Equal(t, util.SafeJsonDump(postOrderMap), util.SafeJsonDump(&PostOrderDetails{
			m: postOrders,
		}))
	})

	t.Run("error_GetQuoteIncrementAndTickSize", func(t *testing.T) {
		defer httpmock.Reset()
		responder := httpmock.NewStringResponder(http.StatusInternalServerError, ``)
		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerDetailsURI, "btcsgd"), responder)

		postOrderMap := &PostOrderDetails{
			m: treemap.NewWithStringComparator(),
		}
		handlerCexApiCalls(ctx, "BTC", postOrderMap)
		assert.Equal(t, util.SafeJsonDump(postOrderMap), util.SafeJsonDump(&PostOrderDetails{
			m: treemap.NewWithStringComparator(),
		}))
	})

	t.Run("error_handlerCexApiCallsOrderOpenThenCancel", func(t *testing.T) {
		defer httpmock.Reset()
		responder := httpmock.NewStringResponder(http.StatusOK, `{
			"tick_size": 1E-8,
			"quote_increment": 0.01
		}`)
		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerDetailsURI, "btcsgd"), responder)

		responder = httpmock.NewStringResponder(http.StatusInternalServerError, ``)
		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerV2URI, "btcsgd"), responder)

		postOrderMap := &PostOrderDetails{
			m: treemap.NewWithStringComparator(),
		}
		handlerCexApiCalls(ctx, "BTC", postOrderMap)
		assert.Equal(t, util.SafeJsonDump(postOrderMap), util.SafeJsonDump(&PostOrderDetails{
			m: treemap.NewWithStringComparator(),
		}))
	})

	t.Run("error_no_fulfilled_order", func(t *testing.T) {
		defer httpmock.Reset()
		responder := httpmock.NewStringResponder(http.StatusOK, `{
			"tick_size": 1E-8,
			"quote_increment": 0.01
		}`)
		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerDetailsURI, "btcsgd"), responder)

		responder = httpmock.NewStringResponder(http.StatusOK, `{
			"bid": "9345.70"
		}`)
		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerV2URI, "btcsgd"), responder)

		responder = httpmock.NewStringResponder(http.StatusOK, `{
				"order_id": "106817811", 
				"avg_execution_price": "3632.8508430064554",
				"is_live": true, 
				"is_cancelled": false, 
				"executed_amount": "3.7567928949",
				"client_order_id": "20190110-4738721"
		}`)
		httpmock.RegisterResponder(http.MethodPost, gemini.NewOrderURI, responder)

		responder = httpmock.NewStringResponder(http.StatusInternalServerError, ``)
		httpmock.RegisterResponder(http.MethodPost, gemini.OrderStatusURI, responder)

		responder = httpmock.NewStringResponder(http.StatusOK, `{
				"order_id": "106817811", 
				"avg_execution_price": "3632.8508430064554",
				"is_live": true, 
				"is_cancelled": true, 
				"executed_amount": "3.7567928949",
				"client_order_id": "20190110-4738721"
		}`)
		httpmock.RegisterResponder(http.MethodPost, gemini.CancelOrderURI, responder)

		postOrderMap := &PostOrderDetails{
			m: treemap.NewWithStringComparator(),
		}
		handlerCexApiCalls(ctx, "BTC", postOrderMap)
		assert.Equal(t, util.SafeJsonDump(postOrderMap), util.SafeJsonDump(&PostOrderDetails{
			m: treemap.NewWithStringComparator(),
		}))
	})
}

func Test_handlerCexApiCallsOrderOpenThenCancel(t *testing.T) {
	ctx := util.TestContext()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	config.TestInit(nil, nil)
	gemini.MustInitClient()

	type args struct {
		ticker         string
		quoteIncrement int
		tickSize       int
	}
	tests := []struct {
		name    string
		setup   func() func()
		args    args
		want    *gemini.Order
		wantErr bool
	}{
		{
			name: "ok_is_fulfilled_upon_creation",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusOK, `{
					"bid": "9345.70"
				}`)
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerV2URI, "btcsgd"), responder)

				responder = httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811", 
						"avg_execution_price": "3632.8508430064554",
						"is_live": false, 
						"is_cancelled": false, 
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`)
				httpmock.RegisterResponder(http.MethodPost, gemini.NewOrderURI, responder)

				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker:         "BTC",
				quoteIncrement: 2,
				tickSize:       8,
			},
			want: &gemini.Order{
				OrderID:           "106817811",
				AvgExecutionPrice: 3632.8508430064554,
				IsLive:            false,
				IsCancelled:       false,
				ExecutedAmount:    3.7567928949,
				ClientOrderID:     "20190110-4738721",
			},
			wantErr: false,
		},
		{
			name: "ok_is_fulfilled_in_query",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusOK, `{
					"bid": "9345.70"
				}`)
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerV2URI, "btcsgd"), responder)

				responder = httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811", 
						"avg_execution_price": "3632.8508430064554",
						"is_live": false, 
						"is_cancelled": false, 
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`)
				httpmock.RegisterResponder(http.MethodPost, gemini.NewOrderURI, responder)

				responder = httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811", 
						"avg_execution_price": "3632.8508430064554",
						"is_live": false, 
						"is_cancelled": false, 
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`)
				httpmock.RegisterResponder(http.MethodPost, gemini.OrderStatusURI, responder)

				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker:         "BTC",
				quoteIncrement: 2,
				tickSize:       8,
			},
			want: &gemini.Order{
				OrderID:           "106817811",
				AvgExecutionPrice: 3632.8508430064554,
				IsLive:            false,
				IsCancelled:       false,
				ExecutedAmount:    3.7567928949,
				ClientOrderID:     "20190110-4738721",
			},
			wantErr: false,
		},
		{
			name: "ok_is_cancelled_upon_creation_recreating_success_fulfilled",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusOK, `{
					"bid": "9345.70"
				}`)
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerV2URI, "btcsgd"), responder)

				responder = httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811", 
						"avg_execution_price": "3632.8508430064554",
						"is_live": false, 
						"is_cancelled": true, 
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`).Times(1)
				httpmock.RegisterResponder(http.MethodPost, gemini.NewOrderURI, responder)

				responder = httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811", 
						"avg_execution_price": "3632.8508430064554",
						"is_live": false, 
						"is_cancelled": false, 
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`).Times(1)
				httpmock.RegisterResponder(http.MethodPost, gemini.NewOrderURI, responder)

				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker:         "BTC",
				quoteIncrement: 2,
				tickSize:       8,
			},
			want: &gemini.Order{
				OrderID:           "106817811",
				AvgExecutionPrice: 3632.8508430064554,
				IsLive:            false,
				IsCancelled:       false,
				ExecutedAmount:    3.7567928949,
				ClientOrderID:     "20190110-4738721",
			},
			wantErr: false,
		},
		{
			name: "error_GetTickerBestBidPrice",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusInternalServerError, ``)
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerV2URI, "btcsgd"), responder)

				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker:         "BTC",
				quoteIncrement: 2,
				tickSize:       8,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error_CreateOrder",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusOK, `{
					"bid": "9345.70"
				}`)
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerV2URI, "btcsgd"), responder)

				responder = httpmock.NewStringResponder(http.StatusInternalServerError, ``)
				httpmock.RegisterResponder(http.MethodPost, gemini.NewOrderURI, responder)

				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker:         "BTC",
				quoteIncrement: 2,
				tickSize:       8,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error_is_cancelled_upon_creation_recreating_failure",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusOK, `{
					"bid": "9345.70"
				}`)
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerV2URI, "btcsgd"), responder)

				responder = httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811", 
						"avg_execution_price": "3632.8508430064554",
						"is_live": false, 
						"is_cancelled": true, 
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`).Times(1)
				httpmock.RegisterResponder(http.MethodPost, gemini.NewOrderURI, responder)

				responder = httpmock.NewStringResponder(http.StatusInternalServerError, ``)
				httpmock.RegisterResponder(http.MethodPost, gemini.NewOrderURI, responder)

				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker:         "BTC",
				quoteIncrement: 2,
				tickSize:       8,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error_is_cancelled_upon_creation_recreating_still_cancelled",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusOK, `{
					"bid": "9345.70"
				}`)
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerV2URI, "btcsgd"), responder)

				responder = httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811", 
						"avg_execution_price": "3632.8508430064554",
						"is_live": false, 
						"is_cancelled": true, 
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`)
				httpmock.RegisterResponder(http.MethodPost, gemini.NewOrderURI, responder)

				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker:         "BTC",
				quoteIncrement: 2,
				tickSize:       8,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error_in_query_ok_in_cancel",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusOK, `{
					"bid": "9345.70"
				}`)
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerV2URI, "btcsgd"), responder)

				responder = httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811", 
						"avg_execution_price": "3632.8508430064554",
						"is_live": true, 
						"is_cancelled": false, 
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`)
				httpmock.RegisterResponder(http.MethodPost, gemini.NewOrderURI, responder)

				responder = httpmock.NewStringResponder(http.StatusInternalServerError, ``)
				httpmock.RegisterResponder(http.MethodPost, gemini.OrderStatusURI, responder)

				responder = httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811", 
						"avg_execution_price": "3632.8508430064554",
						"is_live": false, 
						"is_cancelled": true, 
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`)
				httpmock.RegisterResponder(http.MethodPost, gemini.CancelOrderURI, responder)

				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker:         "BTC",
				quoteIncrement: 2,
				tickSize:       8,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "error_is_cancelled_in_query",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusOK, `{
					"bid": "9345.70"
				}`)
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerV2URI, "btcsgd"), responder)

				responder = httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811", 
						"avg_execution_price": "3632.8508430064554",
						"is_live": true, 
						"is_cancelled": false, 
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`)
				httpmock.RegisterResponder(http.MethodPost, gemini.NewOrderURI, responder)

				responder = httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811", 
						"avg_execution_price": "3632.8508430064554",
						"is_live": false, 
						"is_cancelled": true, 
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`)
				httpmock.RegisterResponder(http.MethodPost, gemini.OrderStatusURI, responder)

				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker:         "BTC",
				quoteIncrement: 2,
				tickSize:       8,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error_in_query_error_in_cancel",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusOK, `{
					"bid": "9345.70"
				}`)
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerV2URI, "btcsgd"), responder)

				responder = httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811", 
						"avg_execution_price": "3632.8508430064554",
						"is_live": true, 
						"is_cancelled": false, 
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`)
				httpmock.RegisterResponder(http.MethodPost, gemini.NewOrderURI, responder)

				responder = httpmock.NewStringResponder(http.StatusInternalServerError, ``)
				httpmock.RegisterResponder(http.MethodPost, gemini.OrderStatusURI, responder)

				responder = httpmock.NewStringResponder(http.StatusInternalServerError, ``)
				httpmock.RegisterResponder(http.MethodPost, gemini.CancelOrderURI, responder)

				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker:         "BTC",
				quoteIncrement: 2,
				tickSize:       8,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := tt.setup()
			got, err := handlerCexApiCallsOrderOpenThenCancel(ctx, tt.args.ticker, tt.args.quoteIncrement, tt.args.tickSize)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
			teardown()
		})
	}
}

func Test_handlerCexApiCallsOrderOpenQueryStatus(t *testing.T) {
	ctx := util.TestContext()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	config.TestInit(nil, nil)
	gemini.MustInitClient()

	type args struct {
		ticker string
		order  *gemini.Order
	}
	tests := []struct {
		name    string
		setup   func() func()
		args    args
		want    *gemini.Order
		want1   bool
		wantErr bool
	}{
		{
			name: "ok_is_fulfilled",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811", 
						"avg_execution_price": "3632.8508430064554",
						"is_live": false, 
						"is_cancelled": false, 
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`)
				httpmock.RegisterResponder(http.MethodPost, gemini.OrderStatusURI, responder)
				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker: "BTC",
				order: &gemini.Order{
					OrderID: "106817811",
				},
			},
			want: &gemini.Order{
				OrderID:           "106817811",
				AvgExecutionPrice: 3632.8508430064554,
				IsLive:            false,
				IsCancelled:       false,
				ExecutedAmount:    3.7567928949,
				ClientOrderID:     "20190110-4738721",
			},
			want1:   false,
			wantErr: false,
		},
		{
			name: "ok_is_cancelled",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811", 
						"avg_execution_price": "3632.8508430064554",
						"is_live": false, 
						"is_cancelled": true, 
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`)
				httpmock.RegisterResponder(http.MethodPost, gemini.OrderStatusURI, responder)
				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker: "BTC",
				order: &gemini.Order{
					OrderID: "106817811",
				},
			},
			want:    nil,
			want1:   true,
			wantErr: false,
		},
		{
			name: "ok_is_not_fulfilled_not_cancelled",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811", 
						"avg_execution_price": "3632.8508430064554",
						"is_live": true, 
						"is_cancelled": false, 
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`)
				httpmock.RegisterResponder(http.MethodPost, gemini.OrderStatusURI, responder)
				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker: "BTC",
				order: &gemini.Order{
					OrderID: "106817811",
				},
			},
			want:    nil,
			want1:   false,
			wantErr: false,
		},
		{
			name: "error_get_order_status",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusInternalServerError, ``)
				httpmock.RegisterResponder(http.MethodPost, gemini.OrderStatusURI, responder)
				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker: "BTC",
				order: &gemini.Order{
					OrderID: "106817811",
				},
			},
			want:    nil,
			want1:   false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := tt.setup()
			got, got1, err := handlerCexApiCallsOrderOpenQueryStatus(ctx, tt.args.ticker, tt.args.order)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			teardown()
		})
	}
}
