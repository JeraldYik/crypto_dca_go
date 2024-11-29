package sentry

import (
	"log"
	"time"

	sentry "github.com/getsentry/sentry-go"
	"github.com/jeraldyik/crypto_dca_go/cmd/config"
)

func MustInit() {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: config.Get().Sentry.Dsn,
	}); err != nil {
		log.Panicf("[sentry.MustInit] Failed to initialize sentry, err: %+v\n", err)
	}
}

func Flush() {
	sentry.Flush(2 * time.Second) // Recommended in doc
}

func CaptureErr(err error) {
	sentry.CaptureException(err)
}
