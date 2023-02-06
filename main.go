package main

import (
	"log"

	serverInit "github.com/danilotadeu/star_wars/server"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	server serverInit.Server
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	server = serverInit.New()
}

func main() {
	server.Start()
}
