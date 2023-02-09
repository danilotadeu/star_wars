package planet

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/danilotadeu/star_wars/app"
	errorsP "github.com/danilotadeu/star_wars/model/errors_handler"
	genericModel "github.com/danilotadeu/star_wars/model/generic"
	planetModel "github.com/danilotadeu/star_wars/model/planet"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
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
		logrus.WithFields(logrus.Fields{"trace": "api.planet.planet.ParseInt"}).Error(err)
		return c.Status(http.StatusBadRequest).JSON(errorsP.ErrorsResponse{
			Message: "Por favor envie o id",
		})
	}

	ctx := c.Context()
	planet, err := p.apps.Planet.GetOneByID(ctx, iplanetId)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "api.planet.planet.GetOneByID"}).Error(err)
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
		logrus.WithFields(logrus.Fields{"trace": "api.planet.planetDelete.ParseInt"}).Error(err)
		return c.Status(http.StatusBadRequest).JSON(errorsP.ErrorsResponse{
			Message: "Por favor envie o id",
		})
	}

	ctx := c.Context()
	err = p.apps.Planet.Delete(ctx, iplanetId)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "api.planet.planetDelete.Delete"}).Error(err)
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
// @Param name query string false "name"
// @Success      200  {object}  planetModel.ResponsePlanets
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
			logrus.WithFields(logrus.Fields{"trace": "api.planet.planets.ParseInt.limit"}).Error(err)
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
			logrus.WithFields(logrus.Fields{"trace": "api.planet.planets.ParseInt.page"}).Error(err)
			return c.Status(http.StatusBadRequest).JSON(errorsP.ErrorsResponse{
				Message: "Por favor envie o page corretamente.",
			})
		}
		ipage = pageConv
	}

	name := c.Query("name")

	planets, err := p.apps.Planet.GetAllPlanets(ctx, ipage, ilimit, name)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "api.planet.planets.GetAllPlanets"}).Error(err)
		if errors.Is(err, planetModel.ErrorPlanetNotFound) {
			return c.Status(http.StatusNotFound).JSON(errorsP.ErrorsResponse{
				Message: "Dados nao encontrados",
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(errorsP.ErrorsResponse{
			Message: "Aconteceu um erro interno..",
		})
	}

	nextPage, previousPage := genericModel.MakePagination(ipage)

	_, err = p.apps.Planet.GetAllPlanets(ctx, *nextPage, ilimit, name)
	if err != nil {
		if !errors.Is(err, planetModel.ErrorPlanetNotFound) {
			logrus.WithFields(logrus.Fields{"trace": "api.planet.planets.GetAllPlanets_1"}).Error(err)
			return c.Status(http.StatusInternalServerError).JSON(errorsP.ErrorsResponse{
				Message: "Aconteceu um erro interno..",
			})
		}
		nextPage = nil
	}

	total, err := p.apps.Planet.GetTotalPlanets(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "api.planet.planets.GetTotalPlanets"}).Error(err)
		return c.Status(http.StatusInternalServerError).JSON(errorsP.ErrorsResponse{
			Message: "Aconteceu um erro interno..",
		})
	}

	return c.Status(http.StatusOK).JSON(planetModel.ResponsePlanets{
		Data: planets,
		ResponsePagination: genericModel.Pagination{
			Count:        *total,
			NextPage:     nextPage,
			PreviousPage: previousPage,
		},
	})
}
