package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/jeraldyik/crypto_dca_go/internal/logger"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func RecoverAndGraceFullyExit() {
	location := "RecoverAndGraceFullyExit"
	if r := recover(); r != nil {
		// sentry.CaptureErr(fmt.Errorf("PANIC: %v", r))
		errStr := fmt.Sprintf("%v", r)
		logger.Fatal(location, errStr, errors.New(errStr))
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
	location := "config.ConvertFloatToPrecString"
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
		errStr := "Unsupported type for ConvertFloatToPrecString"
		logger.Panic(location, errStr, errors.New(errStr))
	}

	return strconv.FormatFloat(f, 'f', prec, 64)
}

func PtrOf[T any](v T) *T {
	return &v
}

// For unit tests only
func RecoverAndGraceFullyExitTestHelper(t *testing.T, wantErr string) {
	if r := recover(); r != nil {
		var err error
		switch x := r.(type) {
		case string:
			err = errors.New(x)
		case error:
			err = x
		case *logrus.Entry:
			err = fmt.Errorf("%+v", x)
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
