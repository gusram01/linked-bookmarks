package healthcheck

import "github.com/gofiber/fiber/v2"

func registerRoutes(c fiber.Router) {
    health := c.Group("/health")

    health.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "ok": true,
            "health": 100,
        })
    })
}

func Bootstrap(r fiber.Router) {
    registerRoutes(r)
}
