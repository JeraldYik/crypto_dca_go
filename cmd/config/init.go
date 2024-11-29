package config

import (
	"reflect"
	"strings"
)

func initConfig() *Config {
	config := &Config{}

	env := mustRetrieveConfigFromEnv(env_EnvKey)
	config.IsSandboxEnv = env != string(production)

	cryptoTickers := mustRetrieveConfigFromEnv(cryptoTickers_EnvKey)
	cryptoTickerSlice := mustTransformArrayStringToArray(cryptoTickers)
	config.CryptoTickers = mustTransformSliceToMap(cryptoTickers_EnvKey, cryptoTickerSlice)

	geminiApiKey := mustRetrieveConfigFromEnv(geminiApiKey_EnvKey)
	config.GeminiApi.ApiKey = geminiApiKey

	geminiApiSecret := mustRetrieveConfigFromEnv(geminiApiSecret_EnvKey)
	config.GeminiApi.ApiSecret = geminiApiSecret

	dailyFiatAmounts := mustRetrieveConfigFromEnv(dailyFiatAmounts_EnvKey)
	config.OrderMetadata.DailyFiatAmount = mustTransformJsonStringToMappedCryptoTickers[float64](dailyFiatAmounts_EnvKey, config, dailyFiatAmounts)

	orderPriceToBidPriceRatio := mustRetrieveConfigFromEnv(orderPriceToBidPriceRatio_EnvKey)
	config.OrderMetadata.OrderPriceToBidPriceRatio = mustParseStrToType(orderPriceToBidPriceRatio_EnvKey, orderPriceToBidPriceRatio, reflect.Float64)

	googleServiceAccountEmail := mustRetrieveConfigFromEnv(googleServiceAccountEmail_EnvKey)
	config.GoogleSheet.ServiceAccountEmail = googleServiceAccountEmail

	googleServiceAccountPrivateKey := mustRetrieveConfigFromEnv(googleServiceAccountPrivateKey_EnvKey)
	config.GoogleSheet.ServiceAccountPrivateKey = strings.ReplaceAll(googleServiceAccountPrivateKey, "\\n", "\n")

	googleSheetID := mustRetrieveConfigFromEnv(googleSheetID_EnvKey)
	config.GoogleSheet.SheetID = googleSheetID

	googleSheetName := mustRetrieveConfigFromEnv(googleSheetName_EnvKey)
	config.GoogleSheet.SheetName = googleSheetName

	startRows := mustRetrieveConfigFromEnv(startRows_EnvKey)
	config.GoogleSheet.startRows = mustTransformJsonStringToMappedCryptoTickers[int](startRows_EnvKey, config, startRows)

	startDate := mustRetrieveConfigFromEnv(startDate_EnvKey)
	config.GoogleSheet.startDate = startDate

	columnRanges := mustRetrieveConfigFromEnv(columnRanges_EnvKey)
	config.GoogleSheet.columnRanges = mustTransformJsonStringToMappedCryptoTickers[string](columnRanges_EnvKey, config, columnRanges)

	dbUserName := mustRetrieveConfigFromEnv(dbUsername_EnvKey)
	config.Db.Username = dbUserName

	dbPassword := mustRetrieveConfigFromEnv(dbPassword_EnvKey)
	config.Db.Password = dbPassword

	dbName := mustRetrieveConfigFromEnv(dbName_EnvKey)
	config.Db.Name = dbName

	dbHost := mustRetrieveConfigFromEnv(dbHost_EnvKey)
	config.Db.Host = dbHost

	sentryDsn := mustRetrieveConfigFromEnv(sentryDsn_EnvKey)
	config.Sentry.Dsn = sentryDsn

	return config
}

func addTimeRelatedConfigs(config *Config) {
	config.GoogleSheet.differenceInDays = mustGetDifferenceInDaysFromStartDate(config)

	config.GoogleSheet.rowRanges = formRowRanges(&config.GoogleSheet)

	config.GoogleSheet.CellRanges = formCellRanges(&config.GoogleSheet)
}
