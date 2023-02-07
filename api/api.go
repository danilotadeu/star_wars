package api

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/danilotadeu/star_wars/api/planet"
	"github.com/danilotadeu/star_wars/app"
	_ "github.com/danilotadeu/star_wars/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title		Star Wars API
// @version		1.0
// @BasePath	/api
func Register(apps *app.Container, port string) {
	fiberRoute := fiber.New()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Gracefully shutting down...")
		_ = fiberRoute.Shutdown()
	}()

	baseAPI := fiberRoute.Group("/api")

	// Planets
	planet.NewAPI(baseAPI.Group("/planets"), apps)

	fiberRoute.Get("/swagger/*", swagger.HandlerDefault)

	log.Println("Registered -> Api")
	fiberRoute.Listen(":" + port)
}
