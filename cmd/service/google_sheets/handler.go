package google_sheets

import (
	"errors"

	"github.com/jeraldyik/crypto_dca_go/cmd/config"
	"github.com/jeraldyik/crypto_dca_go/cmd/util"
	"github.com/jeraldyik/crypto_dca_go/internal/logger"
	"google.golang.org/api/sheets/v4"
)

func (gs *GoogleSheets) GetSheetID() (int64, error) {
	location := "google_sheets.GetSheetID"
	config := config.Get().GoogleSheet
	spreadsheet, err := gs.sheets.Spreadsheets.Get(config.SheetID).Do()
	if err != nil {
		logger.Error(location, "Unable to fetch spreadsheet details", err)
		return 0, err
	}
	for _, sheet := range spreadsheet.Sheets {
		if sheet.Properties.Title == config.SheetName {
			return sheet.Properties.SheetId, nil
		}
	}
	return 0, errors.New("no matching sheet")
}

func (gs *GoogleSheets) BatchUpdate(req *sheets.BatchUpdateSpreadsheetRequest) error {
	location := "google_sheets.BatchUpdate"
	config := config.Get().GoogleSheet
	resp, err := gs.sheets.Spreadsheets.BatchUpdate(config.SheetID, req).Do()
	logger.Info(location, "resp: %v", util.SafeJsonDump(resp))
	return err
}
