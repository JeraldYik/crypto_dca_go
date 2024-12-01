package gemini

import "time"

type TickerDetails struct {
	Symbol                string  `json:"symbol"`
	BaseCurrency          string  `json:"base_currency"`
	QuoteCurrency         string  `json:"quote_currency"`
	TickSize              float64 `json:"tick_size"`
	QuoteIncrement        float64 `json:"quote_increment"`
	MinOrderSize          string  `json:"min_order_size"`
	Status                string  `json:"status"`
	WrapEnabled           bool    `json:"wrap_enabled"`
	ProductType           string  `json:"product_type"`
	ContractType          string  `json:"contract_type"`
	ContractPriceCurrency string  `json:"contract_price_currency"`
}

type TickerV2 struct {
	Symbol  string   `json:"symbol"`
	Open    float64  `json:"open,string"`
	High    float64  `json:"high,string"`
	Low     float64  `json:"low,string"`
	Close   float64  `json:"close,string"`
	Changes []string `json:"changes"`
	Bid     float64  `json:"bid,string"`
	Ask     float64  `json:"ask,string"`
}

type Order struct {
	OrderID           string   `json:"order_id"`
	ClientOrderID     string   `json:"client_order_id"`
	Symbol            string   `json:"symbol"`
	Exchange          string   `json:"exchange"`
	Price             float64  `json:"price,string"`
	AvgExecutionPrice float64  `json:"avg_execution_price,string"`
	Side              string   `json:"side"`
	Type              string   `json:"type"`
	Options           []string `json:"options"`
	Timestamp         string   `json:"timestamp"`
	Timestampms       int64    `json:"timestampms"`
	IsLive            bool     `json:"is_live"`
	IsCancelled       bool     `json:"is_cancelled"`
	Reason            string   `json:"reason"`
	WasForced         bool     `json:"was_forced"`
	ExecutedAmount    float64  `json:"executed_amount,string"`
	RemainingAmount   float64  `json:"remaining_amount,string"`
	OriginalAmount    float64  `json:"original_amount,string"`
	IsHidden          bool     `json:"is_hidden"`
}

type Trade struct {
	Timestamp    int64     `json:"timestamp"`
	Timestampms  int64     `json:"timestampms"`
	TimestampmsT time.Time `json:"timestampmst,omitempty"`
	TradeID      int64     `json:"tid"`
	Price        float64   `json:"price,string"`
	Amount       float64   `json:"amount,string"`
	Exchange     string    `json:"exchange"`
	Type         string    `json:"type"`
	Broken       bool      `json:"broken,omitempty"`
}

type CancelResult struct {
	Result  string              `json:"result"`
	Details CancelResultDetails `json:"details"`
}

type CancelResultDetails struct {
	CancelledOrders []float64 `json:"cancelledOrders"`
	CancelRejects   []float64 `json:"cancelRejects"`
}

type FundBalance struct {
	Currency               string  `json:"currency"`
	Amount                 float64 `json:"amount,string"`
	Available              float64 `json:"available,string"`
	AvailableForWithdrawal float64 `json:"availableForWithdrawal,string"`
	Type                   string  `json:"type"`
}
