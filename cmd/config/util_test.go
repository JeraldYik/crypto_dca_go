package config

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/jeraldyik/crypto_dca_go/cmd/util"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/sheets/v4"
)

func Test_mustBeDefined(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		mustBeDefined("key", "val", true)
		defer util.RecoverAndGraceFullyExitTestHelper(t, "")
	})

	t.Run("panic", func(t *testing.T) {
		defer util.RecoverAndGraceFullyExitTestHelper(t, "Missing environment variable 'key'")
		mustBeDefined("key", "", false)
	})
}

func Test_mustRetrieveConfigFromEnv(t *testing.T) {
	testEnvKey := "TEST_ENV_KEY"
	testEnvValue := "TEST_ENV_VALUE"
	t.Run("ok", func(t *testing.T) {
		defer util.RecoverAndGraceFullyExitTestHelper(t, "")
		defer os.Clearenv()
		os.Setenv(testEnvKey, testEnvValue)
		val := mustRetrieveConfigFromEnv(envKey(testEnvKey))
		assert.Equal(t, testEnvValue, val)
	})

	t.Run("panic", func(t *testing.T) {
		defer util.RecoverAndGraceFullyExitTestHelper(t, "Missing environment variable 'TEST_ENV_KEY'")
		defer os.Setenv(testEnvKey, "")
		mustRetrieveConfigFromEnv(envKey(testEnvKey))
		os.Clearenv()
	})
}

func Test_mustTransformArrayStringToArray(t *testing.T) {
	t.Run("ok_with_square_brackets", func(t *testing.T) {
		defer util.RecoverAndGraceFullyExitTestHelper(t, "")
		val := mustTransformArrayStringToArray("[a,b,c]")
		assert.Equal(t, []string{"a", "b", "c"}, val)
	})
	t.Run("ok_without_square_brackets", func(t *testing.T) {
		defer util.RecoverAndGraceFullyExitTestHelper(t, "")
		val := mustTransformArrayStringToArray("a,b,c")
		assert.Equal(t, []string{"a", "b", "c"}, val)
	})
}

func Test_mustTransformSliceToMap(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		defer util.RecoverAndGraceFullyExitTestHelper(t, "")
		val := mustTransformSliceToMap("key", []string{"a", "b", "c"})
		assert.Equal(t, map[string]bool{"a": true, "b": true, "c": true}, val)
	})
	t.Run("panic - string", func(t *testing.T) {
		defer util.RecoverAndGraceFullyExitTestHelper(t, `There exist a zero type in slice ["a",""] for key 'key'`)
		mustTransformSliceToMap("key", []string{"a", ""})
	})
}

func Test_mustTransformJsonStringToMappedCryptoTickers(t *testing.T) {
	t.Run("ok - float64", func(t *testing.T) {
		config.CryptoTickers = map[string]bool{
			"BTC": true,
			"ETH": true,
		}
		defer util.RecoverAndGraceFullyExitTestHelper(t, "")
		val := mustTransformJsonStringToMappedCryptoTickers[float64]("key", config, `{"BTC":1.23,"ETH":4.56}`)
		assert.Equal(t, map[string]float64{"BTC": 1.23, "ETH": 4.56}, val)
	})
	t.Run("ok - int", func(t *testing.T) {
		config := &Config{
			CryptoTickers: map[string]bool{
				"BTC": true,
				"ETH": true,
			},
		}
		defer util.RecoverAndGraceFullyExitTestHelper(t, "")
		val := mustTransformJsonStringToMappedCryptoTickers[float64]("key", config, `{"BTC":1,"ETH":2}`)
		assert.Equal(t, map[string]float64{"BTC": 1, "ETH": 2}, val)
	})
	t.Run("panic - unable to unmarshal", func(t *testing.T) {
		defer util.RecoverAndGraceFullyExitTestHelper(t, "Unable to unmarshal '{")
		mustTransformJsonStringToMappedCryptoTickers[float64]("key", nil, `{`)
	})
	t.Run("panic - crypto ticker not exist", func(t *testing.T) {
		defer util.RecoverAndGraceFullyExitTestHelper(t, "Crypto Ticker 'SOL' does not exist for key 'key'")
		config := &Config{
			CryptoTickers: map[string]bool{
				"BTC": true,
				"ETH": true,
			},
		}
		mustTransformJsonStringToMappedCryptoTickers[float64]("key", config, `{"SOL":1}`)
	})
}

func Test_mustParseStrToType(t *testing.T) {
	t.Run("ok - float64", func(t *testing.T) {
		defer util.RecoverAndGraceFullyExitTestHelper(t, "")
		val := mustParseStrToType("key", "1.23", reflect.Float64)
		assert.Equal(t, float64(1.23), val)
	})
	t.Run("panic - unable to parse float", func(t *testing.T) {
		defer util.RecoverAndGraceFullyExitTestHelper(t, "Unable to parse key 'key', value: 'a' of type 'float64'")
		mustParseStrToType("key", "a", reflect.Float64)
	})
	t.Run("panic - invalid type", func(t *testing.T) {
		defer util.RecoverAndGraceFullyExitTestHelper(t, "Type is not allowed 'string' for key 'key'")
		mustParseStrToType("key", "a", reflect.String)
	})
}

func Test_mustGetDifferenceInDaysFromStartDate(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		defer util.RecoverAndGraceFullyExitTestHelper(t, "")
		twoDaysBefore := time.Now().Add(time.Hour * -48)
		Set(&Config{
			GoogleSheet: GoogleSheet{
				startDate: fmt.Sprintf("%02d/%02d/%04d", twoDaysBefore.Day(), twoDaysBefore.Month(), twoDaysBefore.Year()),
			},
		})
		timeInit(Get(), nil)
		val := mustGetDifferenceInDaysFromStartDate(Get())
		assert.Equal(t, 2, val)
	})
	t.Run("panic - start date is after now", func(t *testing.T) {
		defer util.RecoverAndGraceFullyExitTestHelper(t, "Start date of recording is later than today")
		twoDaysLater := time.Now().Add(time.Hour * 48)
		Set(&Config{
			GoogleSheet: GoogleSheet{
				startDate: fmt.Sprintf("%v/%v/%v", twoDaysLater.Year(), twoDaysLater.Format("01"), twoDaysLater.Day()),
			},
		})
		timeInit(Get(), nil)
		mustGetDifferenceInDaysFromStartDate(Get())
	})
}

func Test_formCellRanges(t *testing.T) {
	type args struct {
		config *GoogleSheet
	}
	tests := []struct {
		name string
		args args
		want map[string]*sheets.GridRange
	}{
		{
			name: "ok",
			args: args{
				config: &GoogleSheet{
					rowRanges: map[string]int{
						"BTC": 3,
						"ETH": 4,
					},
					columnRanges: map[string]string{
						"BTC": "E:H",
						"ETH": "I:L",
					},
					startRows: map[string]int{
						"BTC": 2,
						"ETH": 3,
					},
					differenceInDays: 3,
				},
			},
			want: map[string]*sheets.GridRange{
				"BTC": {
					StartRowIndex:    2,
					EndRowIndex:      3,
					StartColumnIndex: 4,
					EndColumnIndex:   8,
				},
				"ETH": {
					StartRowIndex:    3,
					EndRowIndex:      4,
					StartColumnIndex: 8,
					EndColumnIndex:   12,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formCellRanges(tt.args.config)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_mustParseColStringToInt(t *testing.T) {
	type args struct {
		col string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "single_letter",
			args: args{
				col: "M",
			},
			want: 12,
		},
		{
			name: "double_letter",
			args: args{
				col: "BC",
			},
			want: 54,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mustParseColStringToInt(tt.args.col); got != tt.want {
				t.Errorf("mustParseColStringToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
