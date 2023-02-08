package app

import (
	"github.com/danilotadeu/star_wars/app/planet"
	"github.com/danilotadeu/star_wars/store"
	"github.com/sirupsen/logrus"
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

	logrus.WithFields(logrus.Fields{"trace": "app"}).Infof("Registered - App")
	return container
}
