package observability

import (
	"time"

	"github.com/getsentry/sentry-go"
	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gusram01/linked-bookmarks/internal/platform/config"
	"github.com/gusram01/linked-bookmarks/internal/platform/logger"
)

var sentryHandler fiber.Handler

func new() fiber.Handler {

	return func(c *fiber.Ctx) error {

		if sentryHandler != nil {
			return sentryHandler(c)
		}

		if err := sentry.Init(sentry.ClientOptions{
			Dsn: config.ENVS.SentryDsn,
		}); err != nil {
			logger.GetLogger().ErrorContext(c.UserContext(), "observability initialization failed: ", "error", err.Error())

			return err
		}

		sentryHandler = sentryfiber.New(sentryfiber.Options{
			Repanic:         false,
			WaitForDelivery: true,
			Timeout:         3 * time.Second,
		})

		return sentryHandler(c)
	}

}

func SentryMiddleware() fiber.Handler {
	if sentryHandler == nil {
		return new()
	}

	return sentryHandler
}
