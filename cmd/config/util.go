package config

import (
	"encoding/json"
	"log"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/sheets/v4"
)

func mustBeDefined(key envKey, val string, ok bool) {
	if !ok || val == "" {
		log.Panicf("[config.mustBeDefined] Missing environment variable '%s'\n", key)
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
	m := make(map[T]bool)
	for _, v := range s {
		reflectValue := reflect.ValueOf(v)
		if reflect.DeepEqual(reflectValue, reflect.Zero(reflectValue.Type()).Interface()) {
			log.Panicf("[config.mustTransformSliceToMap] There exist a zero type in slice %+v for key '%s'\n", s, key)
		}
		m[v] = true
	}
	return m
}

func mustTransformJsonStringToMappedCryptoTickers[T float64 | int | string](key envKey, config *Config, s string) map[string]T {
	m := make(map[string]T)
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		log.Panicf("[config.mustTransformJsonStringToMappedCryptoTickers] Unable to unmarshal '%s'\n", s)
	}
	mustCheckIfCryptoTickerExist(key, config, m)

	return m
}

func mustCheckIfCryptoTickerExist[T float64 | int | string](key envKey, config *Config, m map[string]T) {
	for cryptoTicker := range m {
		if _, ok := config.CryptoTickers[cryptoTicker]; !ok {
			log.Panicf("[config.mustCheckIfCryptoTickerExist] Crypto Ticker '%s' does not exist for key '%s'\n", cryptoTicker, key)
		}
	}
}

// TODO: refactor this
func mustParseStrToType[T float64](key envKey, s string, t reflect.Kind) T {
	switch t {
	case reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			log.Panicf("[config.mustParseStrToType] Unable to parse key '%s', value: '%s' of type '%s'\n", key, s, t)
		}
		return T(f)
	default:
		log.Panicf("[config.mustParseStrToType] Type is not allowed '%s' for key '%s'\n", t, key)
	}
	return 0
}

// Date in format of YYYY/MM/DD
func mustGetDifferenceInDaysFromStartDate(config *Config) int {
	now := GetTime().now
	startDate := mustParseStrToTime(startDate_EnvKey, config.GoogleSheet.startDate)
	if startDate.After(now) {
		log.Panicf("[config.mustGetDifferenceInDaysFromStartDate] Start date of recording is later than today\n")
	}
	duration := now.Sub(startDate)
	return int(math.Floor(duration.Hours() / 24))
}

func mustParseStrToTime(key envKey, s string) time.Time {
	t, err := time.Parse("02/01/2006", s)
	if err != nil {
		log.Panicf("[config.mustParseStrToTime] Unable to parse time string '%s' for key '%s'\n", s, key)
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
	idx := int64(0)
	for _, c := range col {
		if c < 'A' || c > 'Z' {
			log.Panicf("[config.mustParseColStringToInt] Col '%s' is of invalid format\n", col)
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
