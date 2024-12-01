package sentry

import (
	"time"

	sentry "github.com/getsentry/sentry-go"
	"github.com/jeraldyik/crypto_dca_go/cmd/config"
	"github.com/jeraldyik/crypto_dca_go/internal/logger"
)

func MustInit() {
	location := "sentry.MustInit"
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: config.Get().Sentry.Dsn,
	}); err != nil {
		logger.Panic(location, "Failed to initialize sentry, err: %+v", err)
	}
}

func Flush() {
	sentry.Flush(2 * time.Second) // Recommended in doc
}

func CaptureErr(err error) {
	sentry.CaptureException(err)
}
