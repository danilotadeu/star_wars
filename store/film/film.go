package film

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	filmModel "github.com/danilotadeu/star_wars/model/film"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Store is a contract to Film..
//
//go:generate mockgen -destination ../../mock/store/film/film_store_mock.go -package mockStoreFilm . Store
type Store interface {
	GetFilms(ctx context.Context, films []string) ([]filmModel.ResultFilm, error)
	SaveFilm(ctx context.Context, film filmModel.ResultFilm) (*int64, error)
	GetOne(ctx context.Context, name string) (*filmModel.Film, error)
	SaveFilmWithPlanet(ctx context.Context, planetID, filmID int64) (*int64, error)
	GetFilmWithPlanet(ctx context.Context, planetID, filmID int64) (*filmModel.FilmPlanet, error)
	GetFilmsByPlanetIDs(ctx context.Context, planetIDs []int64) ([]filmModel.FilmPlanet, error)
}

type storeImpl struct {
	db          *sql.DB
	urlStarWars string
}

// NewApp init a film
func NewStore(db *sql.DB, urlStarWars string) Store {
	return &storeImpl{
		db:          db,
		urlStarWars: urlStarWars,
	}
}

// GetFilm get films in api star wars..
func (a *storeImpl) GetFilms(ctx context.Context, films []string) ([]filmModel.ResultFilm, error) {
	var resultsFilm []filmModel.ResultFilm
	for _, film := range films {
		id := strings.Split(film, a.urlStarWars+"/films/")[1]

		client := &http.Client{}
		url := a.urlStarWars + "/films/" + strings.TrimRight(id, "/")

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			logrus.WithFields(logrus.Fields{"trace": "store.film.GetFilms.newRequest"}).Error(err)
			return nil, err
		}
		req.Header.Add("Accept", "application/json")
		resp, err := client.Do(req)

		if err != nil {
			logrus.WithFields(logrus.Fields{"trace": "store.film.GetFilms.Do"}).Error(err)
			return nil, err
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.WithFields(logrus.Fields{"trace": "store.film.GetFilms.readAll"}).Error(err)
			return nil, err
		}

		if resp.StatusCode == http.StatusOK {
			var responseFilm filmModel.ResultFilm
			err := json.Unmarshal(respBody, &responseFilm)
			if err != nil {
				logrus.WithFields(logrus.Fields{"trace": "store.film.GetFilms.jsonUnmarshal"}).Error(err)
				return nil, err
			}

			resultsFilm = append(resultsFilm, responseFilm)
		}
	}

	return resultsFilm, nil
}

func (a *storeImpl) SaveFilm(ctx context.Context, film filmModel.ResultFilm) (*int64, error) {
	query := fmt.Sprintf("INSERT INTO film(name, director, release_date) VALUES ('%s','%s','%s')",
		film.Title, film.Director, film.ReleaseDate)
	res, err := a.db.Exec(query)

	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "store.film.SaveFilm.Exec"}).Error(err)
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "store.film.SaveFilm.LastInsertId"}).Error(err)
		return nil, err
	}

	return &lastID, nil
}

func (a *storeImpl) SaveFilmWithPlanet(ctx context.Context, planetID, filmID int64) (*int64, error) {
	query := fmt.Sprintf("INSERT INTO film_planet(planet_id, film_id) VALUES ('%d','%d')",
		planetID, filmID)

	res, err := a.db.Exec(query)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "store.film.SaveFilmWithPlanet.Exec"}).Error(err)
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "store.film.SaveFilmWithPlanet.LastInsertId"}).Error(err)
		return nil, err
	}

	return &lastID, nil
}

func (a *storeImpl) GetOne(ctx context.Context, name string) (*filmModel.Film, error) {
	res, err := a.db.Query("SELECT * FROM film WHERE name = ?", name)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "store.film.GetOne.Query"}).Error(err)
		return nil, err
	}
	defer res.Close()

	if res.Next() {
		var film filmModel.Film
		err := res.Scan(
			&film.ID,
			&film.Name,
			&film.Director,
			&film.ReleaseDate,
			&film.CreatedAt,
		)
		if err != nil {
			logrus.WithFields(logrus.Fields{"trace": "store.film.GetOne.Scan"}).Error(err)
			return nil, err
		}

		return &film, nil
	} else {
		return nil, nil
	}
}

func (a *storeImpl) GetFilmWithPlanet(ctx context.Context, planetID, filmID int64) (*filmModel.FilmPlanet, error) {
	res, err := a.db.Query("SELECT * FROM film_planet WHERE planet_id = ? and film_id = ?", planetID, filmID)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "store.film.GetFilmWithPlanet.Query"}).Error(err)
		return nil, err
	}
	defer res.Close()

	if res.Next() {
		var film filmModel.FilmPlanet
		err := res.Scan(
			&film.PlanetID,
			&film.FilmID,
			&film.CreatedAt,
			&film.DeletedAt,
		)
		if err != nil {
			logrus.WithFields(logrus.Fields{"trace": "store.film.GetFilmWithPlanet.Scan"}).Error(err)
			return nil, err
		}

		return &film, nil
	} else {
		return nil, nil
	}
}

func (a *storeImpl) GetFilmsByPlanetIDs(ctx context.Context, planetIDs []int64) ([]filmModel.FilmPlanet, error) {
	query, args, err := sqlx.In(`SELECT *
					FROM
						star_wars.film_planet
							LEFT JOIN
						star_wars.film ON star_wars.film_planet.film_id = star_wars.film.id
					WHERE star_wars.film_planet.planet_id IN (?);`, planetIDs)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "store.film.GetFilmsByPlanetIDs.In"}).Error(err)
		return nil, err
	}

	query = sqlx.Rebind(sqlx.QUESTION, query)
	res, err := a.db.Query(query, args...)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "store.film.GetFilmsByPlanetIDs.Query"}).Error(err)
		return nil, err
	}
	defer res.Close()

	var films []filmModel.FilmPlanet
	for res.Next() {
		var film filmModel.FilmPlanet
		err := res.Scan(
			&film.PlanetID,
			&film.FilmID,
			&film.CreatedAt,
			&film.DeletedAt,
			&film.Film.ID,
			&film.Film.Name,
			&film.Film.Director,
			&film.Film.ReleaseDate,
			&film.Film.CreatedAt,
		)
		if err != nil {
			logrus.WithFields(logrus.Fields{"trace": "store.film.GetFilmsByPlanetIDs.Scan"}).Error(err)
			return nil, err
		}

		films = append(films, film)
	}

	return films, nil
}
