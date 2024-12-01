package gemini

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jeraldyik/crypto_dca_go/cmd/util"
	"github.com/jeraldyik/crypto_dca_go/internal/logger"
)

// Expect all functions to have error as last return value
func RetryWrapper(ctx context.Context, fnName string, fn any, args ...any) ([]reflect.Value, error) {
	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	if fnType.Kind() != reflect.Func {
		return nil, fmt.Errorf("expected a function, got %s", fnType.Kind())
	}

	if len(args) != fnType.NumIn() {
		return nil, fmt.Errorf("argument count mismatch: expected %d, got %d", fnType.NumIn(), len(args))
	}
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	var err error
	var ok bool
	for i := 0; i < MaxRetryCount; i++ {
		results := fnValue.Call(in)

		errorResult := results[len(results)-1]
		if errorResult.IsNil() {
			return results[:len(results)-1], nil
		}
		err, ok = errorResult.Interface().(error)
		if !ok {
			err := fmt.Errorf("invalid_return_type")
			logger.Error(fnName, "Last return type is not error, is %v instead", err, errorResult.Type().String())
			return nil, err
		}
		if err == nil {
			return results[:len(results)-1], nil
		}

		logger.Error(fnName, "Has some error, retrying in 5 seconds", err)
		if !util.IsTestFlow(ctx) {
			time.Sleep(5 * time.Second)
		}
	}

	logger.Error(fnName, "Retried %v times, returning error\n", err, MaxRetryCount)
	return nil, err
}

// To hardcode this, since metadata does not change
func AppendTickerWithQuoteCurrency(ticker string) string {
	if ticker == BTC || ticker == ETH {
		return strings.ToLower(ticker + SGD)
	}
	return strings.ToLower(ticker + USD)
}
