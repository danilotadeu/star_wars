package planet

import (
	"errors"
	"time"

	filmModel "github.com/danilotadeu/star_wars/model/film"
)

var ErrorPlanetNotFound = errors.New("Planet not found")

type ResultPlanet struct {
	Count    int         `json:"count"`
	Next     *string     `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []Planet    `json:"results"`
}

type Planet struct {
	Name           string    `json:"name"`
	RotationPeriod string    `json:"rotation_period"`
	OrbitalPeriod  string    `json:"orbital_period"`
	Diameter       string    `json:"diameter"`
	Climate        string    `json:"climate"`
	Gravity        string    `json:"gravity"`
	Terrain        string    `json:"terrain"`
	SurfaceWater   string    `json:"surface_water"`
	Population     string    `json:"population"`
	Residents      []string  `json:"residents"`
	Films          []string  `json:"films"`
	Created        time.Time `json:"created"`
	Edited         time.Time `json:"edited"`
	URL            string    `json:"url"`
}

type PlanetDB struct {
	ID        int64            `json:"id"`
	Name      string           `json:"name"`
	Climate   string           `json:"climate"`
	Terrain   string           `json:"terrain"`
	CreatedAt time.Time        `json:"created_at"`
	DeletedAt *time.Time       `json:"deleted_at,omitempty"`
	Films     []filmModel.Film `json:"films,omitempty"`
}
