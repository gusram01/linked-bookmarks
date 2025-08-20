package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	linksHttp "github.com/gusram01/linked-bookmarks/internal/link/infra/http"
	onboardingHttp "github.com/gusram01/linked-bookmarks/internal/onboarding/infra/http"
	"github.com/gusram01/linked-bookmarks/internal/platform/config"
	"github.com/gusram01/linked-bookmarks/internal/platform/database"
	"github.com/gusram01/linked-bookmarks/internal/platform/logger"
	"github.com/gusram01/linked-bookmarks/internal/platform/observability"
	storagekv "github.com/gusram01/linked-bookmarks/internal/platform/storage-kv"
	"github.com/gusram01/linked-bookmarks/internal/shared/models"
)

func main() {

	app := fiber.New(fiber.Config{
		Prefork: true,
		AppName: "linked-bookmarks",
	})
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(healthcheck.New())
	app.Use(logger.SetupFiberLogger())
	app.Use(observability.SentryMiddleware())

	app.Use(limiter.New(limiter.Config{
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "Too many requests, please try again later.",
				"data":    nil,
			})
		},
		Storage: storagekv.GetStorage(),
	}))

	clerk.SetKey(config.Config("GC_MARK_AUTH_KEY"))

	database.Initialize(&models.Link{}, &models.User{}, &models.UserLink{})

	onboardingHttp.Bootstrap(app)
	linksHttp.Bootstrap(app)

	p := config.Config("GC_MARK_PORT")

	go func() {
		fmt.Printf("start listen on port: %s \n", p)
		if err := app.Listen(fmt.Sprintf(":%s", p)); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(
		c,
		os.Interrupt,
		syscall.SIGTTIN,
		syscall.SIGTERM,
	)

	<-c
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()

	fmt.Println("Running cleanup tasks...")

	storagekv.GetStorage().Close()
	// Your cleanup tasks go here
	// db.Close()
	// redisConn.Close()
	fmt.Println("Fiber was successful shutdown.")
}
