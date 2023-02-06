package planet

import (
	"errors"
	"time"
)

var ErrorFilmNotFound = errors.New("Film not found")
var ErrorFilmPlanetNotFound = errors.New("Film Planet not found")

type ResultFilm struct {
	Title        string    `json:"title"`
	EpisodeID    int       `json:"episode_id"`
	OpeningCrawl string    `json:"opening_crawl"`
	Director     string    `json:"director"`
	Producer     string    `json:"producer"`
	ReleaseDate  string    `json:"release_date"`
	Characters   []string  `json:"characters"`
	Planets      []string  `json:"planets"`
	Starships    []string  `json:"starships"`
	Vehicles     []string  `json:"vehicles"`
	Species      []string  `json:"species"`
	Created      time.Time `json:"created"`
	Edited       time.Time `json:"edited"`
	URL          string    `json:"url"`
}

type Film struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Director    string    `json:"director"`
	ReleaseDate time.Time `json:"release_date"`
	CreatedAt   time.Time `json:"created_at"`
}

type FilmPlanet struct {
	FilmID    int64      `json:"film_id"`
	PlanetID  int64      `json:"planet_id"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Film      Film
}
