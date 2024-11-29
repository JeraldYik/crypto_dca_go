package gemini

import (
	"github.com/jeraldyik/crypto_dca_go/cmd/config"
)

var api *Api

func MustInitClient() {
	c := config.Get().GeminiApi
	api = New(c.ApiKey, c.ApiSecret, true)
}

func GetClient() *Api {
	return api
}
