package config

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_initConfig(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// defer os.Clearenv()
		twoDaysBefore := TestNow.Add(-1 * time.Hour * 48)
		s := fmt.Sprintf("%02d/%02d/%04d", twoDaysBefore.Day(), twoDaysBefore.Month(), twoDaysBefore.Year())

		os.Setenv(string(env_EnvKey), "sandbox")
		os.Setenv(string(cryptoTickers_EnvKey), "BTC,ETH")
		os.Setenv(string(geminiApiKey_EnvKey), "gemini_api_key")
		os.Setenv(string(geminiApiSecret_EnvKey), "gemini_api_secret")
		os.Setenv(string(dailyFiatAmounts_EnvKey), `{"BTC":1,"ETH":2}`)
		os.Setenv(string(orderPriceToBidPriceRatio_EnvKey), "0.999")
		os.Setenv(string(googleServiceAccountEmail_EnvKey), "google_service_account_email")
		os.Setenv(string(googleServiceAccountPrivateKey_EnvKey), "google_service_account_private_key")
		os.Setenv(string(googleSheetID_EnvKey), "google_sheets_id")
		os.Setenv(string(googleSheetName_EnvKey), "google_sheets_name")
		os.Setenv(string(startRows_EnvKey), `{"BTC":1,"ETH":2}`)
		os.Setenv(string(columnRanges_EnvKey), `{"BTC": "E:H", "ETH": "I:L"}`)
		os.Setenv(string(startDate_EnvKey), s)
		os.Setenv(string(dbUsername_EnvKey), "db_username")
		os.Setenv(string(dbPassword_EnvKey), "db_password")
		os.Setenv(string(dbName_EnvKey), "db_name")
		os.Setenv(string(dbHost_EnvKey), "db_host")
		os.Setenv(string(sentryDsn_EnvKey), "sentry_dsn")

		// expected
		c := initConfig()
		timeInit(c, &TestNow)
		addTimeRelatedConfigs(c)

		// actual
		TestInit(nil, &TestNow)
		addTimeRelatedConfigs(Get())

		assert.Equal(t, c, Get())
	})
}
