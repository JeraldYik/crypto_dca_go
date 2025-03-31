package config

type envKey string

const (
	env_EnvKey envKey = "ENV"

	cryptoTickers_EnvKey             envKey = "CRYPTO_TICKERS"
	geminiApiKey_EnvKey              envKey = "GEMINI_API_KEY"
	geminiApiSecret_EnvKey           envKey = "GEMINI_API_SECRET"
	dailyFiatAmounts_EnvKey          envKey = "DAILY_FIAT_AMOUNTS"
	orderPriceToBidPriceRatio_EnvKey envKey = "ORDER_PRICE_TO_BID_PRICE_RATIO"

	googleServiceAccountEmail_EnvKey      envKey = "GOOGLE_SERVICE_ACCOUNT_EMAIL"
	googleServiceAccountPrivateKey_EnvKey envKey = "GOOGLE_SERVICE_ACCOUNT_PRIVATE_KEY"
	googleSheetID_EnvKey                  envKey = "GOOGLE_SHEET_ID"
	googleSheetName_EnvKey                envKey = "GOOGLE_SHEET_NAME"
	columnRanges_EnvKey                   envKey = "COLUMN_RANGES"
	startRows_EnvKey                      envKey = "START_ROWS"
	startDate_EnvKey                      envKey = "START_DATE"

	dbUsername_EnvKey envKey = "DB_USERNAME"
	dbPassword_EnvKey envKey = "DB_PASSWORD"
	dbName_EnvKey     envKey = "DB_NAME"
	dbHost_EnvKey     envKey = "DB_HOST"
	dbApiUrl_EnvKey   envKey = "DB_API_URL"
	dbApiKey_EnvKey   envKey = "DB_API_KEY"
	sentryDsn_EnvKey  envKey = "SENTRY_DSN"
)

const (
	production = "production" // ! NOT TO BE USED. To determine if env is production or not
)

// These 2 variables determine the looping logic for leaving orders open, querying, cancelling and re-create order with a different bid price
const (
	OrderOpenThenCancelWindowCount  = 23 // outer loop
	OrderOpenQueryStatusWindowCount = 60 // inner loop
)
