package gemini

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Ticker Details
func (api *Api) tickerDetails(ticker string) (TickerDetails, error) {
	quoteCurrency := AppendTickerWithQuoteCurrency(ticker)
	path := fmt.Sprintf(TickerDetailsURI, quoteCurrency)

	log.Printf("[gemini.tickerDetails] path:%s\n", path)

	var tickerDetails TickerDetails

	body, err := api.request(http.MethodGet, path, nil)
	if err != nil {
		return tickerDetails, err
	}

	if err := json.Unmarshal(body, &tickerDetails); err != nil {
		return tickerDetails, err
	}

	log.Printf("[gemini.tickerDetails] tickerDetails: %+v\n", tickerDetails)

	return tickerDetails, nil
}

// TickerV2
func (api *Api) tickerV2(ticker string) (TickerV2, error) {
	quoteCurrency := AppendTickerWithQuoteCurrency(ticker)
	path := fmt.Sprintf(TickerV2URI, quoteCurrency)

	log.Printf("[gemini.tickerV2] path:%s\n", path)

	var tickerV2 TickerV2

	body, err := api.request(http.MethodGet, path, nil)
	if err != nil {
		return tickerV2, err
	}

	if err := json.Unmarshal(body, &tickerV2); err != nil {
		return tickerV2, err
	}

	log.Printf("[gemini.tickerV2] tickerV2: %+v\n", tickerV2)

	return tickerV2, nil
}
