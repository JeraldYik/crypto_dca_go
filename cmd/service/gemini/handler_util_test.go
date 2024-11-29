package gemini

import (
	"testing"

	"github.com/jeraldyik/crypto_dca_go/cmd/config"
)

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
