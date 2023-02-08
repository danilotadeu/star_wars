package planet

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danilotadeu/star_wars/app"
	mockAppPlanet "github.com/danilotadeu/star_wars/mock/app/planet"
	planetModel "github.com/danilotadeu/star_wars/model/planet"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"gopkg.in/go-playground/assert.v1"
)

func TestHandlerDelete(t *testing.T) {
	endpoint := "/planets/:planetId"
	cases := map[string]struct {
		InputParamID       string
		ExpectedErr        error
		ExpectedStatusCode int
		PrepareMockApp     func(mockPlanetApp *mockAppPlanet.MockApp)
	}{
		"should delete the planet": {
			InputParamID: "1",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockPlanetApp *mockAppPlanet.MockApp) {
				mockPlanetApp.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
			ExpectedStatusCode: http.StatusNoContent,
		},
		"should return error with parse int": {
			InputParamID: "xpto",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockPlanetApp *mockAppPlanet.MockApp) {
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		"should return with planet not found": {
			InputParamID: "1",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockPlanetApp *mockAppPlanet.MockApp) {
				mockPlanetApp.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(planetModel.ErrorPlanetNotFound)
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		"should return error": {
			InputParamID: "1",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockPlanetApp *mockAppPlanet.MockApp) {
				mockPlanetApp.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			mockPlanetApp := mockAppPlanet.NewMockApp(ctrl)
			cs.PrepareMockApp(mockPlanetApp)

			h := apiImpl{
				apps: &app.Container{
					Planet: mockPlanetApp,
				},
			}

			app := fiber.New()
			app.Delete(endpoint, h.planetDelete)
			req := httptest.NewRequest(http.MethodDelete, strings.ReplaceAll(endpoint, ":planetId", cs.InputParamID), nil).WithContext(ctx)
			req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Errorf("Error app.Test: %s", err.Error())
				return
			}

			assert.Equal(t, cs.ExpectedErr, err)
			assert.Equal(t, cs.ExpectedStatusCode, resp.StatusCode)
		})
	}
}

func TestHandlerGetPlanetByID(t *testing.T) {
	endpoint := "/planets/:planetId"
	cases := map[string]struct {
		InputParamID       string
		ExpectedErr        error
		ExpectedStatusCode int
		PrepareMockApp     func(mockPlanetApp *mockAppPlanet.MockApp)
	}{
		"should return success with planet": {
			InputParamID: "1",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockPlanetApp *mockAppPlanet.MockApp) {
				mockPlanetApp.EXPECT().GetOneByID(gomock.Any(), gomock.Any()).Return(&planetModel.PlanetDB{
					ID:      1,
					Name:    "Planet 1",
					Climate: "Climate 1",
					Terrain: "Terrain 1",
				}, nil)
			},
			ExpectedStatusCode: http.StatusOK,
		},
		"should return error with parse int": {
			InputParamID: "xpto",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockPlanetApp *mockAppPlanet.MockApp) {
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		"should return with planet not found": {
			InputParamID: "1",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockPlanetApp *mockAppPlanet.MockApp) {
				mockPlanetApp.EXPECT().GetOneByID(gomock.Any(), gomock.Any()).Return(nil, planetModel.ErrorPlanetNotFound)
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		"should return error": {
			InputParamID: "1",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockPlanetApp *mockAppPlanet.MockApp) {
				mockPlanetApp.EXPECT().GetOneByID(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			mockPlanetApp := mockAppPlanet.NewMockApp(ctrl)
			cs.PrepareMockApp(mockPlanetApp)

			h := apiImpl{
				apps: &app.Container{
					Planet: mockPlanetApp,
				},
			}

			app := fiber.New()
			app.Get(endpoint, h.planet)
			req := httptest.NewRequest(http.MethodGet, strings.ReplaceAll(endpoint, ":planetId", cs.InputParamID), nil).WithContext(ctx)
			req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Errorf("Error app.Test: %s", err.Error())
				return
			}

			assert.Equal(t, cs.ExpectedErr, err)
			assert.Equal(t, cs.ExpectedStatusCode, resp.StatusCode)
		})
	}
}

func TestHandlerGetPlanets(t *testing.T) {
	cases := map[string]struct {
		InputPage          string
		InputLimit         string
		ExpectedErr        error
		ExpectedStatusCode int
		PrepareMockApp     func(mockPlanetApp *mockAppPlanet.MockApp)
	}{
		"should return success with planet": {
			InputPage:   "1",
			InputLimit:  "10",
			ExpectedErr: nil,
			PrepareMockApp: func(mockPlanetApp *mockAppPlanet.MockApp) {
				mockPlanetApp.EXPECT().GetAllPlanets(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*planetModel.PlanetDB{
					{
						ID:      1,
						Name:    "Planet 1",
						Climate: "Climate 1",
						Terrain: "Terrain 1",
					},
				}, nil)
			},
			ExpectedStatusCode: http.StatusOK,
		},
		"should return error with parse int page": {
			ExpectedErr: nil,
			InputPage:   "xpto",
			PrepareMockApp: func(mockPlanetApp *mockAppPlanet.MockApp) {
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		"should return error with parse int limit": {
			ExpectedErr: nil,
			InputLimit:  "xpto",
			PrepareMockApp: func(mockPlanetApp *mockAppPlanet.MockApp) {
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		"should return with planet not found": {
			ExpectedErr: nil,
			PrepareMockApp: func(mockPlanetApp *mockAppPlanet.MockApp) {
				mockPlanetApp.EXPECT().GetAllPlanets(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, planetModel.ErrorPlanetNotFound)
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		"should return error": {
			ExpectedErr: nil,
			PrepareMockApp: func(mockPlanetApp *mockAppPlanet.MockApp) {
				mockPlanetApp.EXPECT().GetAllPlanets(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			mockPlanetApp := mockAppPlanet.NewMockApp(ctrl)
			cs.PrepareMockApp(mockPlanetApp)

			h := apiImpl{
				apps: &app.Container{
					Planet: mockPlanetApp,
				},
			}
			endpoint := "/planets"
			app := fiber.New()
			app.Get(endpoint, h.planets)

			if len(cs.InputPage) > 0 && len(cs.InputLimit) > 0 {
				endpoint += "?page=" + cs.InputPage + "&limit=" + cs.InputLimit
			} else if len(cs.InputPage) > 0 {
				endpoint += "?page=" + cs.InputPage
			} else if len(cs.InputLimit) > 0 {
				endpoint += "?limit=" + cs.InputLimit
			}

			req := httptest.NewRequest(http.MethodGet, endpoint, nil).WithContext(ctx)
			req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Errorf("Error app.Test: %s", err.Error())
				return
			}

			assert.Equal(t, cs.ExpectedErr, err)
			assert.Equal(t, cs.ExpectedStatusCode, resp.StatusCode)
		})
	}
}
