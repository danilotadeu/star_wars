package main

import (
	"context"
	"log"
	"os"

	"github.com/danilotadeu/star_wars/app"
	"github.com/danilotadeu/star_wars/server"
	"github.com/danilotadeu/star_wars/store"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	server := server.New()
	db := server.ConnectDatabase()
	store := store.Register(db, os.Getenv("URL_STARWARS_API"))
	app := app.Register(store)

	err = app.Planet.CreatePlanetsAndFilms(context.Background())
	if err != nil {
		log.Println("Happened a problem to import planet and movies: ", err.Error())
		panic(err)
	}
	db.Close()

	log.Println("Planets and movies created with successfully !!!")
}
