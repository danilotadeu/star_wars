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

// ShowPlanet godoc
// @Summary      Show a planet
// @Description  get planet by ID
// @Tags         planets
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Planet ID"
// @Success      200  {object}  planetModel.PlanetDB
// @Failure      400  {object}  errorsP.ErrorsResponse
// @Failure      404  {object}  errorsP.ErrorsResponse
// @Failure      500  {object}  errorsP.ErrorsResponse
// @Router       /planets/{id} [get]
func (p *apiImpl) planet(c *fiber.Ctx) error {
	planetId := c.Params("planetId")
	iplanetId, err := strconv.ParseInt(planetId, 10, 64)
	if err != nil {
		log.Println("api.planet.planet.ParseInt", err.Error())
		return c.Status(http.StatusBadRequest).JSON(errorsP.ErrorsResponse{
			Message: "Por favor envie o id",
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
			Message: "Aconteceu um erro interno..",
		})
	}

	return c.Status(http.StatusOK).JSON(planet)
}

// DeletePlanet godoc
// @Summary      Delete a planet
// @Description  delete planet by ID
// @Tags         planets
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Planet ID"
// @Success      204
// @Failure      400  {object}  errorsP.ErrorsResponse
// @Failure      404  {object}  errorsP.ErrorsResponse
// @Failure      500  {object}  errorsP.ErrorsResponse
// @Router       /planets/{id} [delete]
func (p *apiImpl) planetDelete(c *fiber.Ctx) error {
	planetId := c.Params("planetId")
	iplanetId, err := strconv.ParseInt(planetId, 10, 64)
	if err != nil {
		log.Println("api.planet.planetDelete.ParseInt", err.Error())
		return c.Status(http.StatusBadRequest).JSON(errorsP.ErrorsResponse{
			Message: "Por favor envie o id",
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
			Message: "Aconteceu um erro interno..",
		})
	}

	return c.Status(http.StatusNoContent).JSON(true)
}

// ListPlanets godoc
// @Summary      List planets
// @Description  get planets
// @Tags         planets
// @Accept       json
// @Produce      json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success      200  {object}  []planetModel.PlanetDB
// @Failure      400  {object}  errorsP.ErrorsResponse
// @Failure      404  {object}  errorsP.ErrorsResponse
// @Failure      500  {object}  errorsP.ErrorsResponse
// @Router       /planets [get]
func (p *apiImpl) planets(c *fiber.Ctx) error {
	ctx := c.Context()

	limit := c.Query("limit")
	var ilimit int64 = 10
	if len(limit) > 0 {
		limitConv, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			log.Println("api.planet.planet.ParseInt.limit", err.Error())
			return c.Status(http.StatusBadRequest).JSON(errorsP.ErrorsResponse{
				Message: "Por favor envie o limit corretamente.",
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
			return c.Status(http.StatusBadRequest).JSON(errorsP.ErrorsResponse{
				Message: "Por favor envie o page corretamente.",
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
