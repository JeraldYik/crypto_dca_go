package google_sheets

import (
	"errors"
	"log"

	"github.com/jeraldyik/crypto_dca_go/cmd/config"
	"google.golang.org/api/sheets/v4"
)

func (gs *GoogleSheets) GetSheetID() (int64, error) {
	config := config.Get().GoogleSheet
	spreadsheet, err := gs.sheets.Spreadsheets.Get(config.SheetID).Do()
	if err != nil {
		log.Printf("[google_sheets.GetSheetID] Unable to fetch spreadsheet details, err: %+v\n", err)
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
	config := config.Get().GoogleSheet
	_, err := gs.sheets.Spreadsheets.BatchUpdate(config.SheetID, req).Do()
	return err
}
