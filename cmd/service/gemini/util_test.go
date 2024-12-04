package gemini

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jeraldyik/crypto_dca_go/cmd/config"
	"github.com/jeraldyik/crypto_dca_go/cmd/util"
	"github.com/stretchr/testify/assert"
)

func TestRetryWrapper(t *testing.T) {
	ctx := util.TestContext()
	generalValidate := func(t *testing.T, got []reflect.Value, err error, want []reflect.Value, wantErr bool) {
		if (err != nil) != wantErr {
			t.Errorf("RetryWrapper() error = %v, wantErr %v", err, wantErr)
			return
		}
		if len(got) != len(want) {
			t.Errorf("RetryWrapper() len = %v, want %v", len(got), len(want))
		}
		for i := range got {
			if !reflect.DeepEqual(got[i].Interface(), want[i].Interface()) {
				t.Errorf("RetryWrapper() got[%d] = %v, want %v", i, util.SafeJsonDump(got[i]), util.SafeJsonDump(want[i]))
			}
		}
	}

	t.Run("no_retry_1_arg_1_ret", func(t *testing.T) {
		fnName := "no_retry_1_arg_1_ret"
		fn := func(arg1 string) (string, error) {
			return "ret1", nil
		}
		args := []any{"arg1"}
		want := []reflect.Value{reflect.ValueOf("ret1")}
		wantErr := false
		got, err := RetryWrapper(ctx, fnName, fn, args...)
		generalValidate(t, got, err, want, wantErr)
		assert.Equal(t, "ret1", got[0].String())
	})

	t.Run("no_retry_2_arg_2_ret", func(t *testing.T) {
		fnName := "no_retry_2_arg_2_ret"
		fn := func(arg1 string, arg2 int) (string, int, error) {
			return "ret1", 10, nil
		}
		args := []any{"arg1", 10}
		want := []reflect.Value{reflect.ValueOf("ret1"), reflect.ValueOf(10)}
		wantErr := false
		got, err := RetryWrapper(ctx, fnName, fn, args...)
		generalValidate(t, got, err, want, wantErr)
		assert.Equal(t, got[0].String(), "ret1")
		assert.Equal(t, 10, int(got[1].Int()))
	})

	t.Run("no_retry_arg_ret_custom_type", func(t *testing.T) {
		fnName := "no_retry_arg_ret_custom_type"
		fn := func(o *Order) (*Order, error) {
			o.ClientOrderID = "test_client_order_id"
			return o, nil
		}
		args := []any{&Order{OrderID: "test_order_id"}}
		want := []reflect.Value{reflect.ValueOf(&Order{
			OrderID:       "test_order_id",
			ClientOrderID: "test_client_order_id",
		})}
		wantErr := false
		got, err := RetryWrapper(ctx, fnName, fn, args...)
		generalValidate(t, got, err, want, wantErr)
		assert.Equal(t, &Order{
			OrderID:       "test_order_id",
			ClientOrderID: "test_client_order_id",
		}, got[0].Interface().(*Order))
	})

	t.Run("retry_max_times_then_succeed", func(t *testing.T) {
		retryCount := 0
		fnName := "retry_max_times_then_succeed"
		fn := func(o *Order) (*Order, error) {
			if retryCount < MaxRetryCount-1 {
				retryCount++
				return nil, errors.New("error")
			}
			return o, nil
		}
		args := []any{&Order{}}
		want := []reflect.Value{reflect.ValueOf(&Order{})}
		wantErr := false
		got, err := RetryWrapper(ctx, fnName, fn, args...)
		generalValidate(t, got, err, want, wantErr)
		assert.Equal(t, &Order{}, got[0].Interface().(*Order))
	})

	t.Run("retry_max_times_then_fail", func(t *testing.T) {
		fnName := "retry_max_times_then_fail"
		fn := func(o *Order) (*Order, error) {
			return nil, errors.New("error")
		}
		args := []any{&Order{}}
		var want []reflect.Value
		wantErr := true
		got, err := RetryWrapper(ctx, fnName, fn, args...)
		generalValidate(t, got, err, want, wantErr)
	})
}

func Test_formCreateOrderReq(t *testing.T) {
	config.TestInit(nil, nil)

	type args struct {
		ticker         string
		bestBid        float64
		quoteIncrement int
		tickSize       int
	}
	tests := []struct {
		name  string
		args  args
		want  string // orderPrice: bestBid * orderMetadata.OrderPriceToBidPriceRatio (quoteIncrement dp)
		want1 string // orderAmount: orderMetadata.DailyFiatAmount[ticker / orderPrice (tickSize dp)
	}{
		{
			name: "ok_BTC",
			args: args{
				ticker:         "BTC",
				bestBid:        3632.85,
				quoteIncrement: 2,
				tickSize:       8,
			},
			want:  "3629.22",
			want1: "0.00027554",
		},
		{
			name: "ok_ETH",
			args: args{
				ticker:         "ETH",
				bestBid:        125.51,
				quoteIncrement: 2,
				tickSize:       6,
			},
			want:  "125.38",
			want1: "0.015951",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := formCreateOrderReq(tt.args.ticker, tt.args.bestBid, tt.args.quoteIncrement, tt.args.tickSize)
			if got != tt.want {
				t.Errorf("formCreateOrderReq() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("formCreateOrderReq() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
