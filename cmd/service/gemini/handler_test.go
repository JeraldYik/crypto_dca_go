package gemini

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/jeraldyik/crypto_dca_go/cmd/config"
	"github.com/stretchr/testify/assert"
)

func TestApi_GetQuoteIncrementAndTickSize(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	type args struct {
		ticker string
	}
	tests := []struct {
		name    string
		setup   func() func()
		args    args
		want    int
		want1   int
		wantErr bool
	}{
		{
			name: "ok",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusOK, `{
					"tick_size": 1E-8,
					"quote_increment": 0.01
				}`)
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(TickerDetailsURI, "btcsgd"), responder)
				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker: "BTC",
			},
			want:  2,
			want1: 8,
		},
		{
			name: "error",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusInternalServerError, ``)
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(TickerDetailsURI, "btcsgd"), responder)
				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker: "BTC",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Api{
				url: "",
			}
			teardown := tt.setup()
			got, got1, err := api.GetQuoteIncrementAndTickSize(tt.args.ticker)
			if (err != nil) != tt.wantErr {
				t.Errorf("Api.GetQuoteIncrementAndTickSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Api.GetQuoteIncrementAndTickSize() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Api.GetQuoteIncrementAndTickSize() got1 = %v, want %v", got1, tt.want1)
			}
			teardown()
		})
	}
}

func TestApi_GetTickerBestBidPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	type args struct {
		ticker string
	}
	tests := []struct {
		name    string
		setup   func() func()
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "ok",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusOK, `{
					"bid": "9345.70"
				}`)
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(TickerV2URI, "btcsgd"), responder)
				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker: "BTC",
			},
			want: 9345.7,
		},
		{
			name: "error",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusInternalServerError, ``)
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(TickerV2URI, "btcsgd"), responder)
				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker: "BTC",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Api{
				url: "",
			}
			teardown := tt.setup()
			got, err := api.GetTickerBestBidPrice(tt.args.ticker)
			if (err != nil) != tt.wantErr {
				t.Errorf("Api.GetTickerBestBidPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Api.GetTickerBestBidPrice() = %v, want %v", got, tt.want)
			}
			teardown()
		})
	}
}

func TestApi_CreateOrder(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	config.TestInit(nil, nil)

	type args struct {
		ticker         string
		bestBid        float64
		tickSize       int
		quoteIncrement int
	}
	tests := []struct {
		name    string
		setup   func() func()
		args    args
		want    *Order
		wantErr bool
	}{
		{
			name: "ok",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811", 
						"avg_execution_price": "3632.8508430064554",
						"is_live": false, 
						"is_cancelled": false, 
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`)
				httpmock.RegisterResponder(http.MethodPost, NewOrderURI, responder)
				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				ticker:         "BTC",
				bestBid:        3632.85,
				tickSize:       8,
				quoteIncrement: 2,
			},
			want: &Order{
				OrderID:           "106817811",
				AvgExecutionPrice: 3632.8508430064554,
				IsLive:            false,
				IsCancelled:       false,
				ExecutedAmount:    3.7567928949,
				ClientOrderID:     "20190110-4738721",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Api{
				url: "",
			}
			teardown := tt.setup()
			got, err := api.CreateOrder(tt.args.ticker, tt.args.bestBid, tt.args.tickSize, tt.args.quoteIncrement)
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

func TestApi_GetOrderStatus(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	config.TestInit(nil, nil)

	type args struct {
		orderID string
	}
	tests := []struct {
		name    string
		setup   func() func()
		args    args
		want    *Order
		wantErr bool
	}{
		{
			name: "ok",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811", 
						"avg_execution_price": "3632.8508430064554",
						"is_live": false, 
						"is_cancelled": false, 
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`)
				httpmock.RegisterResponder(http.MethodPost, OrderStatusURI, responder)
				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				orderID: "106817811",
			},
			want: &Order{
				OrderID:           "106817811",
				AvgExecutionPrice: 3632.8508430064554,
				IsLive:            false,
				IsCancelled:       false,
				ExecutedAmount:    3.7567928949,
				ClientOrderID:     "20190110-4738721",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Api{
				url: "",
			}
			teardown := tt.setup()
			got, err := api.GetOrderStatus(tt.args.orderID)
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

func TestApi_CancelOrder(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	config.TestInit(nil, nil)

	type args struct {
		orderID string
	}
	tests := []struct {
		name    string
		setup   func() func()
		args    args
		want    *Order
		wantErr bool
	}{
		{
			name: "ok",
			setup: func() func() {
				responder := httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811", 
						"avg_execution_price": "3632.8508430064554",
						"is_live": false, 
						"is_cancelled": true, 
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`)
				httpmock.RegisterResponder(http.MethodPost, CancelOrderURI, responder)

				return func() {
					httpmock.Reset()
				}
			},
			args: args{
				orderID: "106817811",
			},
			want: &Order{
				OrderID:           "106817811",
				AvgExecutionPrice: 3632.8508430064554,
				IsLive:            false,
				IsCancelled:       true,
				ExecutedAmount:    3.7567928949,
				ClientOrderID:     "20190110-4738721",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &Api{
				url: "",
			}
			teardown := tt.setup()
			got, err := api.CancelOrder(tt.args.orderID)
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
