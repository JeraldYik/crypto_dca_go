package main

import (
	"context"

	"github.com/jeraldyik/crypto_dca_go/cmd"
	"github.com/jeraldyik/crypto_dca_go/cmd/config"
	"github.com/jeraldyik/crypto_dca_go/cmd/service/db"
	"github.com/jeraldyik/crypto_dca_go/cmd/service/gemini"
	"github.com/jeraldyik/crypto_dca_go/cmd/service/google_sheets"
	"github.com/jeraldyik/crypto_dca_go/cmd/service/sentry"
	"github.com/jeraldyik/crypto_dca_go/cmd/util"
)

func main() {
	defer util.RecoverAndGraceFullyExit()
	ctx := context.Background()

	// setup
	config.MustInit()
	gemini.MustInitClient()
	google_sheets.MustInit(ctx)
	db.MustInit()
	defer db.Close()
	sentry.MustInit()
	defer sentry.Flush()

	// run
	if err := cmd.Run(ctx); err != nil {
		// sentry.CaptureErr(err)
	}
}
