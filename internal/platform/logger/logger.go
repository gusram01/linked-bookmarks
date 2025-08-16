package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	slogfiber "github.com/samber/slog-fiber"
	slogformatter "github.com/samber/slog-formatter"
)

var logger = slog.New(slogformatter.NewFormatterHandler(
	slogformatter.TimezoneConverter(time.UTC),
	slogformatter.TimeFormatter(time.RFC3339, nil),
)(
	slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}),
))

var config = slogfiber.Config{
	WithSpanID:    true,
	WithTraceID:   true,
	WithRequestID: true,
}

func SetupFiberLogger() fiber.Handler {
	return slogfiber.NewWithConfig(logger, config)
}

func GetLogger() *slog.Logger {
	return logger
}
