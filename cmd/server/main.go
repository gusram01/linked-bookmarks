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
	"github.com/gusram01/linked-bookmarks/internal/healthcheck"
	links "github.com/gusram01/linked-bookmarks/internal/link/infra"
	linksHttp "github.com/gusram01/linked-bookmarks/internal/link/infra/http"
	"github.com/gusram01/linked-bookmarks/internal/platform/config"
	"github.com/gusram01/linked-bookmarks/internal/platform/database"
)

func main(){

	app := fiber.New(fiber.Config{
		Prefork: true,
		AppName: "linked-bookmarks",
	})

	app.Use(cors.New())
    clerk.SetKey(config.Config("GC_MARK_AUTH_KEY"))

    database.Initialize(&links.LinkModel{})

    healthcheck.Bootstrap(app)
    linksHttp.Bootstrap(app)

	p := config.Config("GC_MARK_PORT")

	go func() {
		/*
		* TODO: the prefork setting is creating multiple
		* process in the same port. How handle this
		* to prevent unexpected behaviors in prod ??
		*/
		fmt.Printf("start listen on port: %s \n", p)
		if err := app.Listen(fmt.Sprintf(":%s",p));  err != nil {
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

	// Your cleanup tasks go here
	// db.Close()
	// redisConn.Close()
	fmt.Println("Fiber was successful shutdown.")
}
