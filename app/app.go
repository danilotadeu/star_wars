package app

import (
	"log"

	"github.com/danilotadeu/star_wars/app/planet"
	"github.com/danilotadeu/star_wars/store"
)

// Container ...
type Container struct {
	Planet planet.App
}

// Register app container
func Register(store *store.Container) *Container {
	container := &Container{
		Planet: planet.NewApp(store),
	}

	log.Println("Registered -> App")
	return container
}
