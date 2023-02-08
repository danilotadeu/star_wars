package api

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/danilotadeu/star_wars/api/planet"
	"github.com/danilotadeu/star_wars/app"
	_ "github.com/danilotadeu/star_wars/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/sirupsen/logrus"
)

// @title		Star Wars API
// @version		1.0
// @BasePath	/api
func Register(apps *app.Container, port string) {
	fiberRoute := fiber.New()

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt)
	go func() {
		<-gracefulShutdown
		fmt.Println("Gracefully shutting down...")
		_ = fiberRoute.Shutdown()
	}()

	baseAPI := fiberRoute.Group("/api")

	// Planets
	planet.NewAPI(baseAPI.Group("/planets"), apps)

	fiberRoute.Get("/swagger/*", swagger.HandlerDefault)

	logrus.WithFields(logrus.Fields{"trace": "api"}).Infof("Registered - Api")
	fiberRoute.Listen(":" + port)
}
