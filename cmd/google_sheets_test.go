package cmd

import (
	"testing"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/jeraldyik/crypto_dca_go/cmd/config"
	"github.com/jeraldyik/crypto_dca_go/cmd/util"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/sheets/v4"
)

func Test_formBatchUpdateRequest(t *testing.T) {
	config.TestInit(nil, &config.TestNow)

	t.Run("ok", func(t *testing.T) {
		sheetID := 1234
		postOrders := treemap.NewWithStringComparator()
		postOrders.Put("BTC", PostOrder{
			ActualFiatDeposit: 1.002,
			AvgExecutionPrice: 1000,
			ExecutedAmount:    1,
		})
		postOrders.Put("ETH", PostOrder{
			ActualFiatDeposit: 2.004,
			AvgExecutionPrice: 1000,
			ExecutedAmount:    1,
		})
		got := formBatchUpdateRequest(int64(sheetID), postOrders)
		assert.Equal(t, &sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{
				{
					UpdateCells: &sheets.UpdateCellsRequest{
						Range: &sheets.GridRange{
							SheetId:          1234,
							StartRowIndex:    2,
							EndRowIndex:      3,
							StartColumnIndex: 4,
							EndColumnIndex:   8,
						},
						Rows: []*sheets.RowData{
							{
								Values: []*sheets.CellData{
									{UserEnteredValue: &sheets.ExtendedValue{StringValue: &config.TestNowDateStr}},
									{UserEnteredValue: &sheets.ExtendedValue{NumberValue: util.PtrOf(1.002)}},
									{UserEnteredValue: &sheets.ExtendedValue{NumberValue: util.PtrOf(float64(1000))}},
									{UserEnteredValue: &sheets.ExtendedValue{NumberValue: util.PtrOf(float64(1))}},
								},
							},
						},
						Fields: "userEnteredValue",
					},
				},
				{
					UpdateCells: &sheets.UpdateCellsRequest{
						Range: &sheets.GridRange{
							SheetId:          1234,
							StartRowIndex:    3,
							EndRowIndex:      4,
							StartColumnIndex: 8,
							EndColumnIndex:   12,
						},
						Rows: []*sheets.RowData{
							{
								Values: []*sheets.CellData{
									{UserEnteredValue: &sheets.ExtendedValue{StringValue: &config.TestNowDateStr}},
									{UserEnteredValue: &sheets.ExtendedValue{NumberValue: util.PtrOf(2.004)}},
									{UserEnteredValue: &sheets.ExtendedValue{NumberValue: util.PtrOf(float64(1000))}},
									{UserEnteredValue: &sheets.ExtendedValue{NumberValue: util.PtrOf(float64(1))}},
								},
							},
						},
						Fields: "userEnteredValue",
					},
				},
			},
		}, got)
	})
}
