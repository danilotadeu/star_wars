package store

import (
	"database/sql"

	"github.com/danilotadeu/star_wars/store/film"
	"github.com/danilotadeu/star_wars/store/planet"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

// Container ...
type Container struct {
	Planet planet.Store
	Film   film.Store
}

// Register store container
func Register(db *sql.DB, urlStarWars string) *Container {
	container := &Container{
		Planet: planet.NewStore(db, urlStarWars),
		Film:   film.NewStore(db, urlStarWars),
	}

	logrus.WithFields(logrus.Fields{"trace": "store"}).Infof("Registered - Store")
	return container
}
