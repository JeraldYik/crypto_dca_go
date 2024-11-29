package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func RecoverAndGraceFullyExit() {
	if r := recover(); r != nil {
		// sentry.CaptureErr(fmt.Errorf("PANIC: %v", r))
		log.Fatal(r)
	}
}

func SafeJsonDump(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

// For unit tests only
func TestContext() context.Context {
	return context.WithValue(context.Background(), contextTestKey, true)
}

func IsTestFlow(ctx context.Context) bool {
	return ctx.Value(contextTestKey) != nil
}

func NumDecimalPlaces(v float64) int {
	return -int(decimal.NewFromFloat(v).Exponent())
}

func ConvertFloatToPrecString[T float64 | decimal.Decimal](val T, prec int) string {
	if prec < 0 {
		return fmt.Sprintf("%v", val)
	}
	var f float64
	switch v := any(val).(type) {
	case float64:
		f = v
	case decimal.Decimal:
		f, _ = v.Float64()
	default:
		log.Fatal("[config.ConvertFloatToPrecString] Unsupported type for ConvertFloatToPrecString")
	}

	return strconv.FormatFloat(f, 'f', prec, 64)
}

func PtrOf[T any](v T) *T {
	return &v
}

// Only used in unit tests
func RecoverAndGraceFullyExitTestHelper(t *testing.T, wantErr string) {
	if r := recover(); r != nil {
		var err error
		switch x := r.(type) {
		case string:
			err = errors.New(x)
		case error:
			err = x
		default:
			t.Fatal("Recovered panic of invalid type")
		}
		if wantErr != "" {
			assert.Error(t, err, wantErr)
		} else {
			assert.NoError(t, err)
		}
	}
}
