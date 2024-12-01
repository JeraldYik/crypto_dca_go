package cmd

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/jeraldyik/crypto_dca_go/cmd/config"
	"github.com/jeraldyik/crypto_dca_go/cmd/service/db"
	"github.com/jeraldyik/crypto_dca_go/cmd/service/gemini"
	"github.com/jeraldyik/crypto_dca_go/cmd/service/google_sheets"
	"github.com/jeraldyik/crypto_dca_go/cmd/util"
	"github.com/jeraldyik/crypto_dca_go/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/sheets/v4"
)

func TestRun(t *testing.T) {
	ctx := util.TestContext()
	config.TestInit(nil, &config.TestNow)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	gemini.MustInitClient()

	tests := []struct {
		name    string
		setup   func(*mocks.MockGoogleSheetsRepository, *mocks.MockOrderRepository) func()
		wantErr bool
	}{
		{
			name: "ok",
			setup: func(gs *mocks.MockGoogleSheetsRepository, orderDB *mocks.MockOrderRepository) func() {
				defer httpmock.Reset()
				responder := httpmock.NewStringResponder(http.StatusOK, `{
					"tick_size": 1E-8,
					"quote_increment": 0.01
				}`)
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerDetailsURI, "btcsgd"), responder)
				responder = httpmock.NewStringResponder(http.StatusOK, `{
					"tick_size": 1E-6,
					"quote_increment": 0.01
				}`)
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerDetailsURI, "ethsgd"), responder)

				responder = httpmock.NewStringResponder(http.StatusOK, `{
					"bid": "9345.70"
				}`)
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerV2URI, "btcsgd"), responder)
				responder = httpmock.NewStringResponder(http.StatusOK, `{
					"bid": "9345.70"
				}`)
				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(gemini.TickerV2URI, "ethsgd"), responder)

				responder = httpmock.NewStringResponder(http.StatusOK, `{
						"order_id": "106817811",
						"avg_execution_price": "3632.8508430064554",
						"is_live": true,
						"is_cancelled": false,
						"executed_amount": "3.7567928949",
						"client_order_id": "20190110-4738721"
				}`)
				httpmock.RegisterResponder(http.MethodPost, gemini.NewOrderURI, responder)

				gs.EXPECT().GetSheetID().Return(int64(1234), nil)
				gs.EXPECT().BatchUpdate(&sheets.BatchUpdateSpreadsheetRequest{
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
				}).Return(nil)

				orderDB.EXPECT().BulkInsert([]*db.Order{
					{
						Ticker:            "btcsgd",
						CreatedForDay:     config.TestNowDate,
						FiatDepositInSGD:  1.002,
						PricePerCoinInSGD: 1000,
						CoinAmount:        1,
					},
					{
						Ticker:            "ethsgd",
						CreatedForDay:     config.TestNowDate,
						FiatDepositInSGD:  2.004,
						PricePerCoinInSGD: 1000,
						CoinAmount:        1,
					},
				}).Return(nil)

				return func() {
					httpmock.Reset()
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockGS := mocks.NewMockGoogleSheetsRepository(ctrl)
			mockOrderDB := mocks.NewMockOrderRepository(ctrl)
			google_sheets.Set(mockGS)
			db.Set(mockOrderDB)

			teardown := tt.setup(mockGS, mockOrderDB)
			err := Run(ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			teardown()
		})
	}
}
