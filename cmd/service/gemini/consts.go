package gemini

const (
	baseURL    = "https://api.gemini.com"
	sandboxURL = "https://api.sandbox.gemini.com"

	// public
	TickerDetailsURI = "/v1/symbols/details/%s"
	TickerV2URI      = "/v2/ticker/%s"

	// authenticated
	NewOrderURI     = "/v1/order/new"
	ActiveOrdersURI = "/v1/orders"
	OrderStatusURI  = "/v1/order/status"
	CancelOrderURI  = "/v1/order/cancel"
)

const (
	MakerTradingFee float64 = 0.002
)

const (
	MaxRetryCount = 5
)

const (
	BTC = "BTC"
	ETH = "ETH"

	SGD = "SGD"
	USD = "USD"
)
