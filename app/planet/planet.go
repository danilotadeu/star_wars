package planet

import (
	"context"
	"log"

	planetModel "github.com/danilotadeu/star_wars/model/planet"
	"github.com/danilotadeu/star_wars/store"
)

//go:generate mockgen -destination ../../mock/app/planet/planet_app_mock.go -package mockAppPlanet . App
type App interface {
	CreatePlanetsAndFilms(ctx context.Context) error
	SaveFilms(ctx context.Context, films []string, planetID int64) error
	GetOneByID(ctx context.Context, planetID int64) (*planetModel.PlanetDB, error)
	GetAllPlanets(ctx context.Context, page, offset int64, name string) ([]*planetModel.PlanetDB, error)
	Delete(ctx context.Context, planetID int64) error
}

type appImpl struct {
	store *store.Container
}

// NewApp init a planet
func NewApp(store *store.Container) App {
	return &appImpl{
		store: store,
	}
}

// CreatePlanetsAndFilms create planets and films..
func (a *appImpl) CreatePlanetsAndFilms(ctx context.Context) error {
	planetResults, err := a.store.Planet.GetPlanets(ctx)
	if err != nil {
		log.Println("app.planet.CreatePlanetsAndFilms.Store.Planet.GetPlanets", err.Error())
		return err
	}

	for _, planetResult := range planetResults {
		for _, planet := range planetResult.Results {
			planetExist, err := a.store.Planet.GetOne(ctx, planet.Name)
			if err != nil {
				log.Println("app.planet.CreatePlanetsAndFilms.Store.Planet.GetOne", err.Error())
				return err
			}

			var planetID *int64
			if planetExist != nil {
				planetID = &planetExist.ID
			}

			if planetID == nil {
				planetID, err = a.store.Planet.SavePlanet(ctx, planet)
				if err != nil {
					log.Println("app.planet.CreatePlanetsAndFilms.Store.Planet.SavePlanet", err.Error())
					return err
				}
			}

			err = a.SaveFilms(ctx, planet.Films, *planetID)
			if err != nil {
				log.Println("app.planet.CreatePlanetsAndFilms.SaveFilms", err.Error())
				return err
			}
		}
	}

	return nil
}

func (a *appImpl) SaveFilms(ctx context.Context, films []string, planetID int64) error {
	filmsResult, err := a.store.Film.GetFilms(ctx, films)
	if err != nil {
		log.Println("app.planet.SaveFilms.Store.Film.GetFilms", err.Error())
		return err
	}

	for _, film := range filmsResult {
		filmExists, err := a.store.Film.GetOne(ctx, film.Title)
		if err != nil {
			log.Println("app.planet.SaveFilms.Store.Film.GetOne", err.Error())
			return err
		}

		var filmID *int64
		if filmExists != nil {
			filmID = &filmExists.ID
		} else {
			filmID, err = a.store.Film.SaveFilm(ctx, film)
			if err != nil {
				log.Println("app.planet.SaveFilms.Store.Film.SaveFilm", err.Error())
				return err
			}
		}

		filmPlanet, err := a.store.Film.GetFilmWithPlanet(ctx, planetID, *filmID)
		if err != nil {
			log.Println("app.planet.SaveFilms.Store.Film.GetFilmWithPlanet", err.Error())
			return err
		}

		if filmPlanet == nil {
			_, err = a.store.Film.SaveFilmWithPlanet(ctx, planetID, *filmID)
			if err != nil {
				log.Println("app.planet.SaveFilms.Store.Film.SaveFilmWithPlanet", err.Error())
				return err
			}
		}
	}

	return nil
}

func (a *appImpl) GetOneByID(ctx context.Context, planetID int64) (*planetModel.PlanetDB, error) {
	planet, err := a.store.Planet.GetOneByID(ctx, planetID)
	if err != nil {
		log.Println("app.planet.GetOneByID.Store.Planet.GetOneByID", err.Error())
		return nil, err
	}

	films, err := a.store.Film.GetFilmsByPlanetIDs(ctx, []int64{planet.ID})
	if err != nil {
		log.Println("app.planet.GetOneByID.Store.Film.GetFilmsByPlanetIDs", err.Error())
		return nil, err
	}

	for _, film := range films {
		planet.Films = append(planet.Films, film.Film)
	}

	return planet, nil
}

func (a *appImpl) GetAllPlanets(ctx context.Context, page, offset int64, name string) ([]*planetModel.PlanetDB, error) {
	planets, err := a.store.Planet.GetAll(ctx, page, offset, name)
	if err != nil {
		log.Println("app.planet.GetAllPlanets.Store.Planet.GetAll", err.Error())
		return nil, err
	}

	if len(planets) == 0 {
		return nil, planetModel.ErrorPlanetNotFound
	}

	planetIDs := make([]int64, 0, len(planets))
	for _, planet := range planets {
		planetIDs = append(planetIDs, planet.ID)
	}

	films, err := a.store.Film.GetFilmsByPlanetIDs(ctx, planetIDs)
	if err != nil {
		log.Println("app.planet.GetAllPlanets.Store.Film.GetFilmsByPlanetIDs", err.Error())
		return nil, err
	}

	for _, planet := range planets {
		for _, film := range films {
			if planet.ID == film.PlanetID {
				planet.Films = append(planet.Films, film.Film)
			}
		}
	}

	return planets, nil
}

func (a *appImpl) Delete(ctx context.Context, planetID int64) error {
	planet, err := a.store.Planet.GetOneByID(ctx, planetID)
	if err != nil {
		log.Println("app.planet.Delete.Store.Planet.GetOneByID", err.Error())
		return err
	}

	err = a.store.Planet.Delete(ctx, planet.ID)
	if err != nil {
		log.Println("app.planet.Delete.Store.Planet.Delete", err.Error())
		return err
	}
	return nil
}
