package gemini

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/jeraldyik/crypto_dca_go/cmd/util"
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
			log.Printf("[%s] last return type is not error, is %v instead\n", fnName, errorResult.Type().String())
			return nil, fmt.Errorf("invalid_return_type")
		}
		if err == nil {
			return results[:len(results)-1], nil
		}

		log.Printf("[%s] has some error, retrying in 5 seconds, err: %+v\n", fnName, err)
		if !util.IsTestFlow(ctx) {
			time.Sleep(5 * time.Second)
		}
	}

	log.Printf("[%s] retried %v times, returning error\n", fnName, MaxRetryCount)
	return nil, err
}

// To hardcode this, since metadata does not change
func AppendTickerWithQuoteCurrency(ticker string) string {
	if ticker == BTC || ticker == ETH {
		return ticker + SGD
	}
	return ticker + USD
}
