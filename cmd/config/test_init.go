package config

import (
	"time"

	"google.golang.org/api/sheets/v4"
)

type ConfigUpdateable struct {
	IsSandboxEnv    *bool
	DailyFiatAmount map[string]float64
}

var TestNow = time.Date(2024, time.November, 3, 14, 30, 0, 0, time.UTC)
var TestNowDate = time.Date(2024, time.November, 3, 0, 0, 0, 0, time.UTC)
var TestNowDateStr = "03/11/2024"

func TestInit(u *ConfigUpdateable, now *time.Time) {
	config = &Config{
		IsSandboxEnv: true,
		CryptoTickers: map[string]bool{
			"BTC": true,
			"ETH": true,
		},
		GeminiApi: GeminiApi{
			ApiKey:    "gemini_api_key",
			ApiSecret: "gemini_api_secret",
		},
		OrderMetadata: OrderMetadata{
			DailyFiatAmount: map[string]float64{
				"BTC": 1,
				"ETH": 2,
			},
			OrderPriceToBidPriceRatio: 0.999,
		},
		GoogleSheet: GoogleSheet{
			ServiceAccountEmail:      "google_service_account_email",
			ServiceAccountPrivateKey: "google_service_account_private_key",
			SheetID:                  "google_sheets_id",
			SheetName:                "google_sheets_name",
			rowRanges: map[string]int{
				"BTC": 3,
				"ETH": 4,
			},
			CellRanges: map[string]*sheets.GridRange{
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
			columnRanges: map[string]string{
				"BTC": "E:H",
				"ETH": "I:L",
			},
			startDate: "01/11/2024",
			startRows: map[string]int{
				"BTC": 1,
				"ETH": 2,
			},
			differenceInDays: 2,
		},
		Db: Db{
			Name:     "db_name",
			Host:     "db_host",
			Username: "db_username",
			Password: "db_password",
		},
		Sentry: Sentry{
			Dsn: "sentry_dsn",
		},
	}

	if u != nil {
		if u.IsSandboxEnv != nil {
			config.IsSandboxEnv = *u.IsSandboxEnv
		}
		if u.DailyFiatAmount != nil {
			config.OrderMetadata.DailyFiatAmount = u.DailyFiatAmount
		}
	}

	timeInit(config, now)
}
