package limiter

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gusram01/linked-bookmarks/internal"
	storagekv "github.com/gusram01/linked-bookmarks/internal/platform/storage-kv"
)

var limiterHandler fiber.Handler

func new() fiber.Handler {

	return func(c *fiber.Ctx) error {

		if limiterHandler != nil {
			return limiterHandler(c)
		}

		limiterHandler = limiter.New(
			limiter.Config{
				LimitReached: func(c *fiber.Ctx) error {
					return c.Status(fiber.StatusTooManyRequests).JSON(
						internal.NewGcResponse(
							nil,
							errors.New("too many requests, please try again later"),
						),
					)
				},
				Storage:            storagekv.GetStorage(),
				SkipFailedRequests: true,
				Max:                20,
				Expiration:         30 * time.Second,
				LimiterMiddleware:  limiter.SlidingWindow{},
			},
		)

		return limiterHandler(c)
	}
}

func LimiterMiddleware() fiber.Handler {
	if limiterHandler == nil {
		return new()
	}

	return limiterHandler
}
