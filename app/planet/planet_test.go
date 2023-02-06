package planet

import (
	"context"
	"fmt"
	"testing"
	"time"

	mockStoreFilm "github.com/danilotadeu/star_wars/mock/store/film"
	mockStorePlanet "github.com/danilotadeu/star_wars/mock/store/planet"
	filmModel "github.com/danilotadeu/star_wars/model/film"
	planetModel "github.com/danilotadeu/star_wars/model/planet"
	"github.com/danilotadeu/star_wars/store"
	"github.com/golang/mock/gomock"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreatePlanetsAndFilms(t *testing.T) {
	cases := map[string]struct {
		prepareMock func(planetStore *mockStorePlanet.MockStore, filmStore *mockStoreFilm.MockStore)
		expectedErr error
	}{
		"should save planet and films": {
			prepareMock: func(planetStore *mockStorePlanet.MockStore, filmStore *mockStoreFilm.MockStore) {
				planetStore.EXPECT().GetPlanets(gomock.Any()).Return([]planetModel.ResultPlanet{
					{
						Count:    60,
						Next:     new(string),
						Previous: nil,
						Results: []planetModel.Planet{
							{
								Name: "Planet 1",
							},
							{
								Name: "Planet 2",
							},
						},
					},
				}, nil)
				planetStore.EXPECT().GetOne(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
				var planetID int64 = 1
				planetStore.EXPECT().SavePlanet(gomock.Any(), gomock.Any()).AnyTimes().Return(&planetID, nil)
				filmStore.EXPECT().GetFilms(gomock.Any(), gomock.Any()).AnyTimes().Return([]filmModel.ResultFilm{
					{
						Title: "Film 1",
					},
					{
						Title: "Film 2",
					},
				}, nil)
				filmStore.EXPECT().GetOne(gomock.Any(), gomock.Any()).AnyTimes().Return(&filmModel.Film{
					ID: 1,
				}, nil)
				filmStore.EXPECT().GetFilmWithPlanet(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
				var filmPlanetID int64 = 1
				filmStore.EXPECT().SaveFilmWithPlanet(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&filmPlanetID, nil)
			},
			expectedErr: nil,
		},
		"should save planet and films when planet exist": {
			prepareMock: func(planetStore *mockStorePlanet.MockStore, filmStore *mockStoreFilm.MockStore) {
				planetStore.EXPECT().GetPlanets(gomock.Any()).Return([]planetModel.ResultPlanet{
					{
						Count:    60,
						Next:     new(string),
						Previous: nil,
						Results: []planetModel.Planet{
							{
								Name: "Planet 1",
							},
							{
								Name: "Planet 2",
							},
						},
					},
				}, nil)
				planetStore.EXPECT().GetOne(gomock.Any(), gomock.Any()).AnyTimes().Return(&planetModel.PlanetDB{
					ID: 1,
				}, nil)
				filmStore.EXPECT().GetFilms(gomock.Any(), gomock.Any()).AnyTimes().Return([]filmModel.ResultFilm{
					{
						Title: "Film 1",
					},
					{
						Title: "Film 2",
					},
				}, nil)
				filmStore.EXPECT().GetOne(gomock.Any(), gomock.Any()).AnyTimes().Return(&filmModel.Film{
					ID: 1,
				}, nil)
				filmStore.EXPECT().GetFilmWithPlanet(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
				var filmPlanetID int64 = 1
				filmStore.EXPECT().SaveFilmWithPlanet(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&filmPlanetID, nil)
			},
			expectedErr: nil,
		},
		"should return error when get planets": {
			prepareMock: func(planetStore *mockStorePlanet.MockStore, filmStore *mockStoreFilm.MockStore) {
				planetStore.EXPECT().GetPlanets(gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
		"should return error when get one": {
			prepareMock: func(planetStore *mockStorePlanet.MockStore, filmStore *mockStoreFilm.MockStore) {
				planetStore.EXPECT().GetPlanets(gomock.Any()).Return([]planetModel.ResultPlanet{
					{
						Count:    60,
						Next:     new(string),
						Previous: nil,
						Results: []planetModel.Planet{
							{
								Name: "Planet 1",
							},
							{
								Name: "Planet 2",
							},
						},
					},
				}, nil)
				planetStore.EXPECT().GetOne(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, fmt.Errorf("error"))

			},
			expectedErr: fmt.Errorf("error"),
		},
		"should return error when save planet": {
			prepareMock: func(planetStore *mockStorePlanet.MockStore, filmStore *mockStoreFilm.MockStore) {
				planetStore.EXPECT().GetPlanets(gomock.Any()).Return([]planetModel.ResultPlanet{
					{
						Count:    60,
						Next:     new(string),
						Previous: nil,
						Results: []planetModel.Planet{
							{
								Name: "Planet 1",
							},
							{
								Name: "Planet 2",
							},
						},
					},
				}, nil)
				planetStore.EXPECT().GetOne(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
				planetStore.EXPECT().SavePlanet(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
		"should return error when save films": {
			prepareMock: func(planetStore *mockStorePlanet.MockStore, filmStore *mockStoreFilm.MockStore) {
				planetStore.EXPECT().GetPlanets(gomock.Any()).Return([]planetModel.ResultPlanet{
					{
						Count:    60,
						Next:     new(string),
						Previous: nil,
						Results: []planetModel.Planet{
							{
								Name: "Planet 1",
							},
							{
								Name: "Planet 2",
							},
						},
					},
				}, nil)
				planetStore.EXPECT().GetOne(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
				var planetID int64 = 1
				planetStore.EXPECT().SavePlanet(gomock.Any(), gomock.Any()).AnyTimes().Return(&planetID, nil)
				filmStore.EXPECT().GetFilms(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			filmStoreMock := mockStoreFilm.NewMockStore(ctrl)
			planetStoreMock := mockStorePlanet.NewMockStore(ctrl)

			cs.prepareMock(planetStoreMock, filmStoreMock)
			app := NewApp(&store.Container{
				Film:   filmStoreMock,
				Planet: planetStoreMock,
			})

			// when
			err := app.CreatePlanetsAndFilms(ctx)

			// then
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func TestSaveFilms(t *testing.T) {
	cases := map[string]struct {
		inputFilms      []string
		inputPlanetID   int64
		prepareMock     func(filmStore *mockStoreFilm.MockStore)
		expectedPlanets []*planetModel.PlanetDB
		expectedErr     error
	}{
		"should save films already exists in table films": {
			inputFilms:    []string{"1", "2"},
			inputPlanetID: 5,
			prepareMock: func(filmStore *mockStoreFilm.MockStore) {
				filmStore.EXPECT().GetFilms(gomock.Any(), gomock.Any()).Return([]filmModel.ResultFilm{
					{
						Title: "Film 1",
					},
					{
						Title: "Film 2",
					},
				}, nil)
				filmStore.EXPECT().GetOne(gomock.Any(), gomock.Any()).AnyTimes().Return(&filmModel.Film{
					ID: 1,
				}, nil)
				filmStore.EXPECT().GetFilmWithPlanet(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
				var filmPlanetID int64 = 1
				filmStore.EXPECT().SaveFilmWithPlanet(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&filmPlanetID, nil)
			},
			expectedErr: nil,
		},
		"should save films that not exists": {
			inputFilms:    []string{"1", "2"},
			inputPlanetID: 5,
			prepareMock: func(filmStore *mockStoreFilm.MockStore) {
				filmStore.EXPECT().GetFilms(gomock.Any(), gomock.Any()).Return([]filmModel.ResultFilm{
					{
						Title: "Film 1",
					},
					{
						Title: "Film 2",
					},
				}, nil)
				filmStore.EXPECT().GetOne(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
				var filmID int64 = 1
				filmStore.EXPECT().SaveFilm(gomock.Any(), gomock.Any()).AnyTimes().Return(&filmID, nil)
				filmStore.EXPECT().GetFilmWithPlanet(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
				var filmPlanetID int64 = 1
				filmStore.EXPECT().SaveFilmWithPlanet(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&filmPlanetID, nil)
			},
			expectedErr: nil,
		},
		"should save films but not save films and planets": {
			inputFilms:    []string{"1", "2"},
			inputPlanetID: 5,
			prepareMock: func(filmStore *mockStoreFilm.MockStore) {
				filmStore.EXPECT().GetFilms(gomock.Any(), gomock.Any()).Return([]filmModel.ResultFilm{
					{
						Title: "Film 1",
					},
					{
						Title: "Film 2",
					},
				}, nil)
				filmStore.EXPECT().GetOne(gomock.Any(), gomock.Any()).AnyTimes().Return(&filmModel.Film{
					ID: 1,
				}, nil)
				filmStore.EXPECT().GetFilmWithPlanet(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&filmModel.FilmPlanet{
					FilmID:   1,
					PlanetID: 1,
				}, nil)
			},
			expectedErr: nil,
		},
		"should return error when save film": {
			inputFilms:    []string{"1", "2"},
			inputPlanetID: 5,
			prepareMock: func(filmStore *mockStoreFilm.MockStore) {
				filmStore.EXPECT().GetFilms(gomock.Any(), gomock.Any()).Return([]filmModel.ResultFilm{
					{
						Title: "Film 1",
					},
					{
						Title: "Film 2",
					},
				}, nil)
				filmStore.EXPECT().GetOne(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
				filmStore.EXPECT().SaveFilm(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
		"should return error when get films": {
			inputFilms:    []string{"1", "2"},
			inputPlanetID: 5,
			prepareMock: func(filmStore *mockStoreFilm.MockStore) {
				filmStore.EXPECT().GetFilms(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
		"should return error when get one": {
			inputFilms:    []string{"1", "2"},
			inputPlanetID: 5,
			prepareMock: func(filmStore *mockStoreFilm.MockStore) {
				filmStore.EXPECT().GetFilms(gomock.Any(), gomock.Any()).Return([]filmModel.ResultFilm{
					{
						Title: "Film 1",
					},
					{
						Title: "Film 2",
					},
				}, nil)
				filmStore.EXPECT().GetOne(gomock.Any(), gomock.Any()).AnyTimes().Return(nil, fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
		"should return error when get film with planet": {
			inputFilms:    []string{"1", "2"},
			inputPlanetID: 5,
			prepareMock: func(filmStore *mockStoreFilm.MockStore) {
				filmStore.EXPECT().GetFilms(gomock.Any(), gomock.Any()).Return([]filmModel.ResultFilm{
					{
						Title: "Film 1",
					},
					{
						Title: "Film 2",
					},
				}, nil)
				filmStore.EXPECT().GetOne(gomock.Any(), gomock.Any()).AnyTimes().Return(&filmModel.Film{
					ID: 1,
				}, nil)
				filmStore.EXPECT().GetFilmWithPlanet(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
		"should return error when save film with planet": {
			inputFilms:    []string{"1", "2"},
			inputPlanetID: 5,
			prepareMock: func(filmStore *mockStoreFilm.MockStore) {
				filmStore.EXPECT().GetFilms(gomock.Any(), gomock.Any()).Return([]filmModel.ResultFilm{
					{
						Title: "Film 1",
					},
					{
						Title: "Film 2",
					},
				}, nil)
				filmStore.EXPECT().GetOne(gomock.Any(), gomock.Any()).AnyTimes().Return(&filmModel.Film{
					ID: 1,
				}, nil)
				filmStore.EXPECT().GetFilmWithPlanet(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, nil)
				filmStore.EXPECT().SaveFilmWithPlanet(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil, fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			filmStoreMock := mockStoreFilm.NewMockStore(ctrl)

			cs.prepareMock(filmStoreMock)
			app := NewApp(&store.Container{
				Film: filmStoreMock,
			})

			// when
			err := app.SaveFilms(ctx, cs.inputFilms, cs.inputPlanetID)

			// then
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func TestGetAllPlanets(t *testing.T) {
	dateString := "2021-11-22"
	date, _ := time.Parse("2006-01-02", dateString)
	planetExpected := []*planetModel.PlanetDB{
		{
			ID:        1,
			Name:      "Planet 1",
			Climate:   "Climate 1",
			Terrain:   "Terrain 1",
			CreatedAt: date,
			Films: []filmModel.Film{
				{
					ID:          1,
					Name:        "Film 1",
					Director:    "Director 1",
					ReleaseDate: date,
					CreatedAt:   date,
				},
				{
					ID:          2,
					Name:        "Film 2",
					Director:    "Director 2",
					ReleaseDate: date,
					CreatedAt:   date,
				},
			},
		},
	}
	cases := map[string]struct {
		inputPage       int64
		inputOffset     int64
		inputName       string
		prepareMock     func(planetStore *mockStorePlanet.MockStore, filmStore *mockStoreFilm.MockStore)
		expectedPlanets []*planetModel.PlanetDB
		expectedErr     error
	}{
		"should get all planets": {
			inputPage:   0,
			inputOffset: 5,
			inputName:   "Planet 1",
			prepareMock: func(planetStore *mockStorePlanet.MockStore, filmStore *mockStoreFilm.MockStore) {
				planetStore.EXPECT().GetAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*planetModel.PlanetDB{
					{
						ID:        1,
						Name:      "Planet 1",
						Climate:   "Climate 1",
						Terrain:   "Terrain 1",
						CreatedAt: date,
					},
				}, nil)
				filmStore.EXPECT().GetFilmsByPlanetIDs(gomock.Any(), gomock.Any()).Return([]filmModel.FilmPlanet{
					{
						FilmID:    1,
						PlanetID:  1,
						CreatedAt: date,
						Film: filmModel.Film{
							ID:          1,
							Name:        "Film 1",
							Director:    "Director 1",
							ReleaseDate: date,
							CreatedAt:   date,
						},
					},
					{
						FilmID:    2,
						PlanetID:  1,
						CreatedAt: date,
						Film: filmModel.Film{
							ID:          2,
							Name:        "Film 2",
							Director:    "Director 2",
							ReleaseDate: date,
							CreatedAt:   date,
						},
					},
				}, nil)
			},
			expectedPlanets: planetExpected,
			expectedErr:     nil,
		},
		"should return empty planets": {
			inputPage:   0,
			inputOffset: 5,
			inputName:   "Planet 1",
			prepareMock: func(planetStore *mockStorePlanet.MockStore, filmStore *mockStoreFilm.MockStore) {
				planetStore.EXPECT().GetAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			expectedPlanets: nil,
			expectedErr:     planetModel.ErrorPlanetNotFound,
		},
		"should return error when get all planets": {
			inputPage:   0,
			inputOffset: 5,
			inputName:   "Planet 1",
			prepareMock: func(planetStore *mockStorePlanet.MockStore, filmStore *mockStoreFilm.MockStore) {
				planetStore.EXPECT().GetAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			expectedPlanets: nil,
			expectedErr:     fmt.Errorf("error"),
		},
		"should return error when get films by planet ids": {
			inputPage:   0,
			inputOffset: 5,
			inputName:   "Planet 1",
			prepareMock: func(planetStore *mockStorePlanet.MockStore, filmStore *mockStoreFilm.MockStore) {
				planetStore.EXPECT().GetAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*planetModel.PlanetDB{
					{
						ID:        1,
						Name:      "Planet 1",
						Climate:   "Climate 1",
						Terrain:   "Terrain 1",
						CreatedAt: date,
					},
				}, nil)
				filmStore.EXPECT().GetFilmsByPlanetIDs(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			expectedPlanets: nil,
			expectedErr:     fmt.Errorf("error"),
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			planetStoreMock := mockStorePlanet.NewMockStore(ctrl)
			filmStoreMock := mockStoreFilm.NewMockStore(ctrl)

			cs.prepareMock(planetStoreMock, filmStoreMock)
			app := NewApp(&store.Container{
				Planet: planetStoreMock,
				Film:   filmStoreMock,
			})

			// when
			planets, err := app.GetAllPlanets(ctx, cs.inputPage, cs.inputOffset, cs.inputName)

			// then
			assert.Equal(t, cs.expectedErr, err)
			assert.Equal(t, cs.expectedPlanets, planets)
		})
	}
}

func TestDelete(t *testing.T) {
	cases := map[string]struct {
		inputPlanet int64
		prepareMock func(planetStore *mockStorePlanet.MockStore)
		expectedErr error
	}{
		"should delete a planet by id": {
			inputPlanet: 1,
			prepareMock: func(planetStore *mockStorePlanet.MockStore) {
				planetStore.EXPECT().GetOneByID(gomock.Any(), gomock.Any()).Return(&planetModel.PlanetDB{
					ID: 1,
				}, nil)
				planetStore.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
		"should return error when get a planet by id": {
			inputPlanet: 1,
			prepareMock: func(planetStore *mockStorePlanet.MockStore) {
				planetStore.EXPECT().GetOneByID(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
		"should return error when delete a planet by id": {
			inputPlanet: 1,
			prepareMock: func(planetStore *mockStorePlanet.MockStore) {
				planetStore.EXPECT().GetOneByID(gomock.Any(), gomock.Any()).Return(&planetModel.PlanetDB{
					ID: 1,
				}, nil)
				planetStore.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
			},
			expectedErr: fmt.Errorf("error"),
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			planetStoreMock := mockStorePlanet.NewMockStore(ctrl)

			cs.prepareMock(planetStoreMock)
			app := NewApp(&store.Container{
				Planet: planetStoreMock,
			})

			// when
			err := app.Delete(ctx, cs.inputPlanet)

			// then
			assert.Equal(t, cs.expectedErr, err)
		})
	}
}

func TestGetOneByID(t *testing.T) {
	dateString := "2021-11-22"
	date, _ := time.Parse("2006-01-02", dateString)
	planetExpected := planetModel.PlanetDB{
		ID:        1,
		Name:      "Planet 1",
		Climate:   "Climate 1",
		Terrain:   "Terrain 1",
		CreatedAt: date,
		Films: []filmModel.Film{
			{
				ID:          1,
				Name:        "Film 1",
				Director:    "Director 1",
				ReleaseDate: date,
				CreatedAt:   date,
			},
			{
				ID:          2,
				Name:        "Film 2",
				Director:    "Director 2",
				ReleaseDate: date,
				CreatedAt:   date,
			},
		},
	}

	cases := map[string]struct {
		inputPlanet    int64
		prepareMock    func(planetStore *mockStorePlanet.MockStore, filmStore *mockStoreFilm.MockStore)
		expectedPlanet *planetModel.PlanetDB
		expectedErr    error
	}{
		"should return a planet with success": {
			inputPlanet: 1,
			prepareMock: func(planetStore *mockStorePlanet.MockStore, filmStore *mockStoreFilm.MockStore) {
				planetStore.EXPECT().GetOneByID(gomock.Any(), gomock.Any()).Return(&planetModel.PlanetDB{
					ID:        1,
					Name:      "Planet 1",
					Climate:   "Climate 1",
					Terrain:   "Terrain 1",
					CreatedAt: date,
				}, nil)
				filmStore.EXPECT().GetFilmsByPlanetIDs(gomock.Any(), gomock.Any()).Return([]filmModel.FilmPlanet{
					{
						FilmID:    1,
						PlanetID:  1,
						CreatedAt: date,
						Film: filmModel.Film{
							ID:          1,
							Name:        "Film 1",
							Director:    "Director 1",
							ReleaseDate: date,
							CreatedAt:   date,
						},
					},
					{
						FilmID:    2,
						PlanetID:  1,
						CreatedAt: date,
						Film: filmModel.Film{
							ID:          2,
							Name:        "Film 2",
							Director:    "Director 2",
							ReleaseDate: date,
							CreatedAt:   date,
						},
					},
				}, nil)
			},
			expectedPlanet: &planetExpected,
			expectedErr:    nil,
		},
		"should return error when get a planet by id": {
			inputPlanet: 1,
			prepareMock: func(planetStore *mockStorePlanet.MockStore, filmStore *mockStoreFilm.MockStore) {
				planetStore.EXPECT().GetOneByID(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			expectedPlanet: nil,
			expectedErr:    fmt.Errorf("error"),
		},
		"should return error when get the films by planet id": {
			inputPlanet: 1,
			prepareMock: func(planetStore *mockStorePlanet.MockStore, filmStore *mockStoreFilm.MockStore) {
				planetStore.EXPECT().GetOneByID(gomock.Any(), gomock.Any()).Return(&planetModel.PlanetDB{
					ID:        1,
					Name:      "Planet 1",
					Climate:   "Climate 1",
					Terrain:   "Terrain 1",
					CreatedAt: date,
				}, nil)
				filmStore.EXPECT().GetFilmsByPlanetIDs(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			expectedPlanet: nil,
			expectedErr:    fmt.Errorf("error"),
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			// given
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()

			planetStoreMock := mockStorePlanet.NewMockStore(ctrl)
			filmStoreMock := mockStoreFilm.NewMockStore(ctrl)

			cs.prepareMock(planetStoreMock, filmStoreMock)
			app := NewApp(&store.Container{
				Planet: planetStoreMock,
				Film:   filmStoreMock,
			})

			// when
			planetDB, err := app.GetOneByID(ctx, cs.inputPlanet)

			// then
			assert.Equal(t, cs.expectedErr, err)
			assert.Equal(t, planetDB, cs.expectedPlanet)

		})
	}
}
