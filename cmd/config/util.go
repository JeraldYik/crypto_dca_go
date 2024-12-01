package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jeraldyik/crypto_dca_go/internal/logger"
	"google.golang.org/api/sheets/v4"
)

func mustBeDefined(key envKey, val string, ok bool) {
	location := "config.mustBeDefined"
	if !ok || val == "" {
		errStr := fmt.Sprintf("Missing environment variable '%s'", key)
		logger.Panic(location, errStr, errors.New(errStr))
	}
}

func mustRetrieveConfigFromEnv(key envKey) string {
	val, ok := os.LookupEnv(string(key))
	mustBeDefined(key, val, ok)
	return val
}

// Can contain square brackets or without
func mustTransformArrayStringToArray(arrayString string) []string {
	// to check for square brackets
	if len(arrayString) >= 2 && (string(arrayString[0]) == "[" || string(arrayString[len(arrayString)-1]) == "]") {
		arrayString = arrayString[1 : len(arrayString)-1]
	}
	return strings.Split(arrayString, ",")
}

func mustTransformSliceToMap[T string](key envKey, s []T) map[T]bool {
	location := "config.mustTransformSliceToMap"
	m := make(map[T]bool)
	for _, v := range s {
		reflectValue := reflect.ValueOf(v)
		if reflect.DeepEqual(reflectValue, reflect.Zero(reflectValue.Type()).Interface()) {
			errStr := fmt.Sprintf("There exist a zero type in slice %+v for key '%s'", s, key)
			logger.Panic(location, errStr, errors.New(errStr))
		}
		m[v] = true
	}
	return m
}

func mustTransformJsonStringToMappedCryptoTickers[T float64 | int | string](key envKey, config *Config, s string) map[string]T {
	location := "config.mustTransformJsonStringToMappedCryptoTickers"
	m := make(map[string]T)
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		errStr := fmt.Sprintf("Unable to unmarshal '%s'", key)
		logger.Panic(location, errStr, errors.New(errStr))
	}
	mustCheckIfCryptoTickerExist(key, config, m)

	return m
}

func mustCheckIfCryptoTickerExist[T float64 | int | string](key envKey, config *Config, m map[string]T) {
	location := "config.mustCheckIfCryptoTickerExist"
	for cryptoTicker := range m {
		if _, ok := config.CryptoTickers[cryptoTicker]; !ok {
			errStr := fmt.Sprintf("Crypto Ticker '%s' does not exist for key '%s'", cryptoTicker, key)
			logger.Panic(location, errStr, errors.New(errStr))
		}
	}
}

// TODO: refactor this
func mustParseStrToType[T float64](key envKey, s string, t reflect.Kind) T {
	location := "config.mustParseStrToType"
	switch t {
	case reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			errStr := fmt.Sprintf("Unable to parse key '%s', value: '%s' of type '%s'", key, s, t)
			logger.Panic(location, errStr, errors.New(errStr))
		}
		return T(f)
	default:
		errStr := fmt.Sprintf("Type '%s' is not allowed for key '%s'", t, key)
		logger.Panic(location, errStr, errors.New(errStr))
	}
	return 0
}

// Date in format of YYYY/MM/DD
func mustGetDifferenceInDaysFromStartDate(config *Config) int {
	location := "config.mustGetDifferenceInDaysFromStartDate"
	now := GetTime().now
	startDate := mustParseStrToTime(startDate_EnvKey, config.GoogleSheet.startDate)
	if startDate.After(now) {
		errStr := "Start date of recording is later than today"
		logger.Panic(location, errStr, errors.New(errStr))
	}
	duration := now.Sub(startDate)
	return int(math.Floor(duration.Hours() / 24))
}

func mustParseStrToTime(key envKey, s string) time.Time {
	location := "config.mustParseStrToTime"
	t, err := time.Parse("02/01/2006", s)
	if err != nil {
		errStr := fmt.Sprintf("Unable to parse date string '%s' for key '%s'", s, key)
		logger.Panic(location, errStr, errors.New(errStr))
	}
	return t
}

func formRowRanges(c *GoogleSheet) map[string]int {
	rowRanges := make(map[string]int)
	for ticker, startRow := range c.startRows {
		rowRanges[ticker] = startRow + c.differenceInDays
	}
	return rowRanges
}

func mustParseColStringToInt(col string) int64 {
	location := "config.mustParseColStringToInt"
	idx := int64(0)
	for _, c := range col {
		if c < 'A' || c > 'Z' {
			errStr := fmt.Sprintf("Col '%s' is of invalid format", col)
			logger.Panic(location, errStr, errors.New(errStr))
		}
		idx = idx*26 + int64(c-'A') + 1
	}
	return idx - 1
}

func formCellRanges(config *GoogleSheet) map[string]*sheets.GridRange {
	cellRanges := make(map[string]*sheets.GridRange)
	for ticker, row := range config.rowRanges {
		cols := strings.Split(config.columnRanges[ticker], ":")
		cellRanges[ticker] = &sheets.GridRange{
			StartRowIndex:    int64(row - 1),
			EndRowIndex:      int64(row),
			StartColumnIndex: mustParseColStringToInt(cols[0]),
			EndColumnIndex:   mustParseColStringToInt(cols[1]) + 1,
		}
	}

	return cellRanges
}
