package cmd

import (
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/jeraldyik/crypto_dca_go/cmd/config"
	"github.com/jeraldyik/crypto_dca_go/cmd/service/google_sheets"
	"github.com/jeraldyik/crypto_dca_go/cmd/util"
	"github.com/jeraldyik/crypto_dca_go/internal/logger"
	"google.golang.org/api/sheets/v4"
)

func formBatchUpdateRequest(sheetID int64, postOrders *treemap.Map) *sheets.BatchUpdateSpreadsheetRequest {
	cellRanges := config.Get().GoogleSheet.CellRanges

	req := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: make([]*sheets.Request, postOrders.Size()),
	}
	i := 0
	it := postOrders.Iterator()
	for it.Next() {
		ticker, postOrder := it.Key().(string), it.Value().(PostOrder)
		cellRange := cellRanges[ticker]
		cellRange.SheetId = sheetID

		req.Requests[i] = &sheets.Request{
			UpdateCells: &sheets.UpdateCellsRequest{
				Range: cellRange,
				Rows: []*sheets.RowData{
					{
						Values: []*sheets.CellData{ // per row, i.e. per ticker
							{UserEnteredValue: &sheets.ExtendedValue{StringValue: util.PtrOf(config.GetTime().GetNowDateString())}},
							{UserEnteredValue: &sheets.ExtendedValue{NumberValue: util.PtrOf(postOrder.ActualFiatDeposit)}},
							{UserEnteredValue: &sheets.ExtendedValue{NumberValue: util.PtrOf(postOrder.AvgExecutionPrice)}},
							{UserEnteredValue: &sheets.ExtendedValue{NumberValue: util.PtrOf(postOrder.ExecutedAmount)}},
						},
					},
				},
				Fields: "userEnteredValue",
			},
		}
		i++
	}

	return req
}

func batchUpdate(postOrders *treemap.Map) error {
	location := "cmd.batchUpdate"
	sheetID, err := google_sheets.Get().GetSheetID()
	if err != nil {
		logger.Error(location, "Getting google sheets sheet ID", err)
		return err
	}
	googleSheetsReq := formBatchUpdateRequest(sheetID, postOrders)
	err = google_sheets.Get().BatchUpdate(googleSheetsReq)
	if err != nil {
		logger.Error(location, "Batch updating google sheets", err)
		return err
	}

	return nil
}
