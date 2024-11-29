package google_sheets

import (
	"context"
	"log"

	"github.com/jeraldyik/crypto_dca_go/cmd/config"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

//go:generate mockgen -source=/cmd/service/google_sheets/main.go -destination=/mocks/mock_GoogleSheetsRepository.go -package=mocks
type GoogleSheetsRepository interface {
	GetSheetID() (int64, error)
	BatchUpdate(req *sheets.BatchUpdateSpreadsheetRequest) error
}

type GoogleSheets struct {
	sheets *sheets.Service
}

var googleSheets GoogleSheetsRepository

func MustInit(ctx context.Context) {
	// Create a JWT configurations object for the Google service account
	conf := &jwt.Config{
		Email:      config.Get().GoogleSheet.ServiceAccountEmail,
		PrivateKey: []byte(config.Get().GoogleSheet.ServiceAccountPrivateKey),
		TokenURL:   googleServiceAccountTokenUri,
		Scopes:     []string{sheets.SpreadsheetsScope},
	}
	client := conf.Client(ctx)
	// Create a service object for Google sheets
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Panicf("[google_sheets.MustInit] Failed to initialise Google Sheets, err: %+v", err)
	}
	Set(&GoogleSheets{sheets: srv})
}

func Get() GoogleSheetsRepository {
	return googleSheets
}

func Set(gs GoogleSheetsRepository) {
	googleSheets = gs
}
