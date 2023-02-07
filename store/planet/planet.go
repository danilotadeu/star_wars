package planet

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	planetModel "github.com/danilotadeu/star_wars/model/planet"
)

// Store is a contract to Planet..
//
//go:generate mockgen -destination ../../mock/store/planet/planet_store_mock.go -package mockStorePlanet . Store
type Store interface {
	GetPlanets(ctx context.Context) ([]planetModel.ResultPlanet, error)
	SavePlanet(ctx context.Context, planet planetModel.Planet) (*int64, error)
	GetOne(ctx context.Context, name string) (*planetModel.PlanetDB, error)
	GetOneByID(ctx context.Context, id int64) (*planetModel.PlanetDB, error)
	GetAll(ctx context.Context, page, limit int64, name string) ([]*planetModel.PlanetDB, error)
	Delete(ctx context.Context, id int64) error
}

type storeImpl struct {
	db          *sql.DB
	urlStarWars string
}

// NewApp init a planet
func NewStore(db *sql.DB, urlStarWars string) Store {
	return &storeImpl{
		db:          db,
		urlStarWars: urlStarWars,
	}
}

// GetPlanets get all planets..
func (a *storeImpl) GetPlanets(ctx context.Context) ([]planetModel.ResultPlanet, error) {
	page := 1
	var resultsPlanet []planetModel.ResultPlanet

	for {
		client := &http.Client{}
		url := a.urlStarWars + "/planets"

		if page > 1 {
			url += "/?page=" + strconv.Itoa(page)
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("store.planet.GetPlanets.newRequest", err.Error())
			return nil, err
		}
		req.Header.Add("Accept", "application/json")
		resp, err := client.Do(req)

		if err != nil {
			log.Println("store.planet.GetPlanets.Do", err.Error())
			return nil, err
		}

		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("store.planet.GetPlanets.readAll", err.Error())
			return nil, err
		}

		if resp.StatusCode == http.StatusOK {
			var responsePlanet planetModel.ResultPlanet
			err := json.Unmarshal(respBody, &responsePlanet)
			if err != nil {
				log.Println("store.planet.GetPlanets.jsonUnmarshal", err.Error())
				return nil, err
			}

			resultsPlanet = append(resultsPlanet, responsePlanet)
			if responsePlanet.Next == nil {
				break
			}

			page++
		}
	}

	return resultsPlanet, nil
}

func (a *storeImpl) SavePlanet(ctx context.Context, planet planetModel.Planet) (*int64, error) {
	query := fmt.Sprintf("INSERT INTO planet(name, climate, terrain) VALUES ('%s','%s','%s')",
		planet.Name, planet.Climate, planet.Terrain)
	res, err := a.db.Exec(query)

	if err != nil {
		log.Println("store.planet.SavePlanet.Exec", err.Error())
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Println("store.planet.SavePlanet.LastInsertId", err.Error())
		return nil, err
	}

	return &lastId, nil
}

func (a *storeImpl) GetOne(ctx context.Context, name string) (*planetModel.PlanetDB, error) {
	res, err := a.db.Query("SELECT * FROM planet WHERE deleted_at IS NULL and name = ?", name)
	if err != nil {
		log.Println("store.film.GetOne.Query", err.Error())
		return nil, err
	}
	defer res.Close()

	if res.Next() {
		var planet planetModel.PlanetDB
		err := res.Scan(
			&planet.ID,
			&planet.Name,
			&planet.Climate,
			&planet.Terrain,
			&planet.CreatedAt,
			&planet.DeletedAt,
		)
		if err != nil {
			log.Println("store.film.GetOne.Scan", err.Error())
			return nil, err
		}

		return &planet, nil
	} else {
		return nil, nil
	}
}

func (a *storeImpl) GetOneByID(ctx context.Context, id int64) (*planetModel.PlanetDB, error) {
	res, err := a.db.Query("SELECT * FROM planet WHERE deleted_at IS NULL and id = ?", id)
	if err != nil {
		log.Println("store.planet.GetOneByID.Query", err.Error())
		return nil, err
	}
	defer res.Close()

	if res.Next() {
		var planet planetModel.PlanetDB
		err := res.Scan(
			&planet.ID,
			&planet.Name,
			&planet.Climate,
			&planet.Terrain,
			&planet.CreatedAt,
			&planet.DeletedAt,
		)
		if err != nil {
			log.Println("store.planet.GetOneByID.Scan", err.Error())
			return nil, err
		}

		return &planet, nil
	} else {
		return nil, planetModel.ErrorPlanetNotFound
	}
}

func (a *storeImpl) GetAll(ctx context.Context, page, limit int64, name string) ([]*planetModel.PlanetDB, error) {
	query := `SELECT * FROM planet WHERE deleted_at IS NULL`
	params := []interface{}{}
	if len(name) > 0 {
		params = append(params, "%"+name+"%")
		query += ` AND name LIKE ? `
	}

	query += ` LIMIT ? OFFSET ?`
	params = append(params, limit, page)

	res, err := a.db.Query(query, params...)
	if err != nil {
		log.Println("store.planet.GetAll.Query", err.Error())
		return nil, err
	}
	defer res.Close()

	var results []*planetModel.PlanetDB
	for res.Next() {
		var planet planetModel.PlanetDB
		err := res.Scan(
			&planet.ID,
			&planet.Name,
			&planet.Climate,
			&planet.Terrain,
			&planet.CreatedAt,
			&planet.DeletedAt,
		)
		if err != nil {
			log.Println("store.planet.GetAll.Scan", err.Error())
			return nil, err
		}
		results = append(results, &planet)
	}

	return results, nil
}

func (a *storeImpl) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("UPDATE planet SET deleted_at = '%s' WHERE id = '%d'", time.Now().Format("2006-01-02 15:04:05"), id)
	res, err := a.db.Exec(query)
	if err != nil {
		log.Println("store.planet.Delete.Exec_1", err.Error())
		return err

	}
	_, err = res.RowsAffected()
	if err != nil {
		log.Println("store.planet.Delete.RowsAffected_1", err.Error())
		return err
	}

	resFilmPlanet, err := a.db.Exec("UPDATE film_planet SET deleted_at = ? WHERE planet_id = ?", time.Now().Format("2006-01-02 15:04:05"), id)
	if err != nil {
		log.Println("store.planet.Delete.Exec_2", err.Error())
		return err

	}
	_, err = resFilmPlanet.RowsAffected()
	if err != nil {
		log.Println("store.planet.Delete.RowsAffected_2", err.Error())
		return err
	}

	return nil
}
