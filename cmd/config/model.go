package config

import "google.golang.org/api/sheets/v4"

// Public vars are to be used in the application
//
// Private vars are only declared in env config, and not used elsewhere
type Config struct {
	IsSandboxEnv  bool
	CryptoTickers map[string]bool
	OrderMetadata OrderMetadata
	GeminiApi     GeminiApi
	GoogleSheet   GoogleSheet
	Db            Db
	Sentry        Sentry
}

type OrderMetadata struct {
	DailyFiatAmount           map[string]float64
	OrderPriceToBidPriceRatio float64
}

type GeminiApi struct {
	ApiKey    string
	ApiSecret string
}

// Public vars are to be used in the application
//
// Private vars are only declared in env config, and not used elsewhere
type GoogleSheet struct {
	ServiceAccountEmail        string
	ServiceAccountPrivateKey   string
	ServiceAccountPrivateKeyID string
	SheetID                    string
	SheetName                  string
	CellRanges                 map[string]*sheets.GridRange
	columnRanges               map[string]string
	startDate                  string
	startRows                  map[string]int
	differenceInDays           int
	rowRanges                  map[string]int
}

// Using Supabase
type Db struct {
	Name     string
	Host     string
	Username string
	Password string
	ApiUrl   string
	ApiKey   string
}

type Sentry struct {
	Dsn string
}
