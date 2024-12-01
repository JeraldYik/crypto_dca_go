package gemini

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jeraldyik/crypto_dca_go/internal/logger"
)

// Ticker Details
func (api *Api) tickerDetails(ticker string) (TickerDetails, error) {
	location := "gemini.tickerDetails"
	quoteCurrency := AppendTickerWithQuoteCurrency(ticker)
	path := fmt.Sprintf(TickerDetailsURI, quoteCurrency)

	logger.Info(location, "path:%s", path)

	var tickerDetails TickerDetails

	body, err := api.request(http.MethodGet, path, nil)
	if err != nil {
		return tickerDetails, err
	}

	if err := json.Unmarshal(body, &tickerDetails); err != nil {
		return tickerDetails, err
	}

	logger.Info(location, "tickerDetails: %+v", tickerDetails)

	return tickerDetails, nil
}

// TickerV2
func (api *Api) tickerV2(ticker string) (TickerV2, error) {
	location := "gemini.tickerV2"
	quoteCurrency := AppendTickerWithQuoteCurrency(ticker)
	path := fmt.Sprintf(TickerV2URI, quoteCurrency)

	logger.Info(location, "path:%s", path)

	var tickerV2 TickerV2

	body, err := api.request(http.MethodGet, path, nil)
	if err != nil {
		return tickerV2, err
	}

	if err := json.Unmarshal(body, &tickerV2); err != nil {
		return tickerV2, err
	}

	logger.Info(location, "tickerV2: %+v", tickerV2)

	return tickerV2, nil
}
