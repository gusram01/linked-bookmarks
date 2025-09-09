package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gusram01/linked-bookmarks/internal/ai"
	linksHttp "github.com/gusram01/linked-bookmarks/internal/link/infra/http"
	onboardingHttp "github.com/gusram01/linked-bookmarks/internal/onboarding/infra/http"
	"github.com/gusram01/linked-bookmarks/internal/platform/config"
	"github.com/gusram01/linked-bookmarks/internal/platform/database"
	"github.com/gusram01/linked-bookmarks/internal/platform/logger"
	"github.com/gusram01/linked-bookmarks/internal/platform/observability"
	storagekv "github.com/gusram01/linked-bookmarks/internal/platform/storage-kv"
	"github.com/gusram01/linked-bookmarks/internal/shared/models"
	vectordb "github.com/gusram01/linked-bookmarks/internal/vector-db"
	"github.com/gusram01/linked-bookmarks/internal/worker"
)

func main() {
	config.LoadConfigFile()

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

	ai.Start()

	clerk.SetKey(config.ENVS.AuthProviderApiKey)

	database.Initialize(&models.Link{}, &models.User{}, &models.UserLink{}, &models.Tag{}, &models.TagLink{})
	vectordb.Initialize()

	worker.CentralWorkerPool.Run()
	onboardingHttp.Bootstrap(app)
	linksHttp.Bootstrap(app)

	p := config.ENVS.ApiPort

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTTIN, syscall.SIGTERM)

	go func() {
		logger.GetLogger().Info("start listen on port ", "port", p)
		address := fmt.Sprintf(":%s", p)
		if err := app.Listen(address); err != nil {

			if err != http.ErrServerClosed {
				log.Fatalf("Could not listen on %s: %v\n", address, err)
			}
		}
	}()

	<-quit
	logger.GetLogger().Info("🚨 Shutdown signal received.")

	_ = app.Shutdown()

	logger.GetLogger().Info("Starting cleanup tasks...")

	storagekv.GetStorage().Close()
	worker.CentralWorkerPool.Shutdown()
	vectordb.VDB.Shutdown()
	// Your cleanup tasks go here

	logger.GetLogger().Info("✅ Application shut down gracefully.")
}
