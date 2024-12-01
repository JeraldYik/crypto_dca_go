package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func Init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
		PadLevelText:  true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)
}

func Info(location, msg string, tags ...any) {
	logrus.Info(formatLogMsg(location, msg, tags...))
}

func Warn(location, msg string, tags ...any) {
	logrus.Warn(formatLogMsg(location, msg, tags...))
}

func Error(location, msg string, err error, tags ...any) {
	logrus.WithField("error", err).Error(formatLogMsg(location, msg, tags...))
}

func Fatal(location, msg string, err error, tags ...any) {
	// Calls os.Exit(1) after logging
	logrus.WithField("error", err).Fatal(formatLogMsg(location, msg, tags...))
}

func Panic(location, msg string, err error, tags ...any) {
	// Calls panic() after logging
	logrus.WithField("error", err).Panic(formatLogMsg(location, msg, tags...))
}

func formatLogMsg(location, msg string, tags ...any) string {
	return fmt.Sprintf("[%s] %v", location, fmt.Sprintf(msg, tags...))
}
