package util

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestNumDecimalPlaces(t *testing.T) {
	type args struct {
		v float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "2_dp",
			args: args{v: 0.01},
			want: 2,
		},
		{
			name: "8_dp",
			args: args{v: 1e-8},
			want: 8,
		},
		{
			name: "positive",
			args: args{v: 12.34},
			want: 2,
		},
		{
			name: "no_dp",
			args: args{v: 12},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumDecimalPlaces(tt.args.v); got != tt.want {
				t.Errorf("NumDecimalPlaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertFloatToPrecString(t *testing.T) {
	type Float interface {
		~float64 | decimal.Decimal
	}
	type args[T Float] struct {
		v    T
		prec int
	}
	testsFloat64 := []struct {
		name string
		args args[float64]
		want string
	}{
		{
			name: "pos_prec_round_up",
			args: args[float64]{v: 12.3456, prec: 2},
			want: "12.35",
		},
		{
			name: "pos_prec_round_down",
			args: args[float64]{v: 12.3446, prec: 2},
			want: "12.34",
		},
		{
			name: "pos_prec_exceed",
			args: args[float64]{v: 12.3456, prec: 8},
			want: "12.34560000",
		},
		{
			name: "neg_prec",
			args: args[float64]{v: 12.3446, prec: -1},
			want: "12.3446",
		},
		{
			name: "zero_prec",
			args: args[float64]{v: 12.3446, prec: 0},
			want: "12",
		},
	}
	testsDecimal := []struct {
		name string
		args args[decimal.Decimal]
		want string
	}{
		{
			name: "decimal.Decimal",
			args: args[decimal.Decimal]{v: decimal.NewFromFloat(12.3456), prec: 2},
			want: "12.35",
		},
	}
	for _, tt := range testsFloat64 {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertFloatToPrecString(tt.args.v, tt.args.prec); got != tt.want {
				t.Errorf("ConvertFloatToPrecString() float64 = %v, want %v", got, tt.want)
			}
		})
	}
	for _, tt := range testsDecimal {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertFloatToPrecString(tt.args.v, tt.args.prec); got != tt.want {
				t.Errorf("ConvertFloatToPrecString() decimal.Decimal = %v, want %v", got, tt.want)
			}
		})
	}
}
