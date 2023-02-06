package api

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/danilotadeu/star_wars/api/planet"
	"github.com/danilotadeu/star_wars/app"
	"github.com/gofiber/fiber/v2"
)

// Register ...
func Register(apps *app.Container) *fiber.App {
	fiberRoute := fiber.New()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Gracefully shutting down...")
		_ = fiberRoute.Shutdown()
	}()

	// Planets
	planet.NewAPI(fiberRoute.Group("/planets"), apps)

	log.Println("Registered -> Api")
	return fiberRoute
}
