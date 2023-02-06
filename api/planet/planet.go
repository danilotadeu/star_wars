package planet

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/danilotadeu/star_wars/app"
	errorsP "github.com/danilotadeu/star_wars/model/errors_handler"
	planetModel "github.com/danilotadeu/star_wars/model/planet"
	"github.com/gofiber/fiber/v2"
)

type apiImpl struct {
	apps *app.Container
}

// NewAPI planet function..
func NewAPI(g fiber.Router, apps *app.Container) {
	api := apiImpl{
		apps: apps,
	}

	g.Get("/", api.planets)
	g.Get("/:planetId", api.planet)
	g.Delete("/:planetId", api.planetDelete)
}

func (p *apiImpl) planet(c *fiber.Ctx) error {
	planetId := c.Params("planetId")
	iplanetId, err := strconv.ParseInt(planetId, 10, 64)
	if err != nil {
		log.Println("api.planet.planet.ParseInt", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(errorsP.ErrorsResponse{
			Message: "Por favor tente mais tarde...",
		})
	}

	ctx := c.Context()
	planet, err := p.apps.Planet.GetOneByID(ctx, iplanetId)
	if err != nil {
		log.Println("api.planet.planet.GetOneByID", err.Error())
		if errors.Is(err, planetModel.ErrorPlanetNotFound) {
			return c.Status(http.StatusNotFound).JSON(errorsP.ErrorsResponse{
				Message: fmt.Sprintf("Planeta (%d) não encontrado", iplanetId),
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(errorsP.ErrorsResponse{
			Message: "Por favor tente mais tarde...",
		})
	}

	return c.Status(http.StatusOK).JSON(planet)
}

func (p *apiImpl) planetDelete(c *fiber.Ctx) error {
	planetId := c.Params("planetId")
	iplanetId, err := strconv.ParseInt(planetId, 10, 64)
	if err != nil {
		log.Println("api.planet.planetDelete.ParseInt", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(errorsP.ErrorsResponse{
			Message: "Por favor tente mais tarde...",
		})
	}

	ctx := c.Context()
	err = p.apps.Planet.Delete(ctx, iplanetId)
	if err != nil {
		log.Println("api.planet.planetDelete.Delete", err.Error())
		if errors.Is(err, planetModel.ErrorPlanetNotFound) {
			return c.Status(http.StatusNotFound).JSON(errorsP.ErrorsResponse{
				Message: fmt.Sprintf("Planeta (%d) não encontrado", iplanetId),
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(errorsP.ErrorsResponse{
			Message: "Por favor tente mais tarde...",
		})
	}

	return c.Status(http.StatusOK).JSON(true)
}

func (p *apiImpl) planets(c *fiber.Ctx) error {
	ctx := c.Context()

	limit := c.Query("limit")
	var ilimit int64 = 10
	if len(limit) > 0 {
		limitConv, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			log.Println("api.planet.planet.ParseInt.limit", err.Error())
			return c.Status(http.StatusInternalServerError).JSON(errorsP.ErrorsResponse{
				Message: "Por favor tente mais tarde...",
			})
		}
		ilimit = limitConv
	}

	page := c.Query("page")
	var ipage int64

	if len(page) > 0 {
		pageConv, err := strconv.ParseInt(page, 10, 64)
		if err != nil {
			log.Println("api.planet.planet.ParseInt.page", err.Error())
			return c.Status(http.StatusInternalServerError).JSON(errorsP.ErrorsResponse{
				Message: "Por favor tente mais tarde...",
			})
		}
		ipage = pageConv
	}

	name := c.Query("name")

	planets, err := p.apps.Planet.GetAllPlanets(ctx, ipage, ilimit, name)
	if err != nil {
		log.Println("api.planet.planet.GetAllPlanets", err.Error())
		if errors.Is(err, planetModel.ErrorPlanetNotFound) {
			return c.Status(http.StatusNotFound).JSON(errorsP.ErrorsResponse{
				Message: "Dados nao encontrados",
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(errorsP.ErrorsResponse{
			Message: "Por favor tente mais tarde...",
		})
	}

	return c.Status(http.StatusOK).JSON(planets)
}
