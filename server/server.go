package server

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/danilotadeu/star_wars/api"
	"github.com/danilotadeu/star_wars/app"
	"github.com/danilotadeu/star_wars/store"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const LOGS_PATH = "log/logrus.log"

// Server is a interface to define contract to server up
type Server interface {
	Start()
	ConnectDatabase() *sql.DB
}

type server struct {
	App   *app.Container
	Store *store.Container
	Db    *sql.DB
}

// New is instance the server
func New() Server {
	return &server{}
}

func (e *server) Start() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename: LOGS_PATH,
		MaxSize:  50, // megabytes
	}))

	e.Db = e.ConnectDatabase()
	e.Store = store.Register(e.Db, os.Getenv("URL_STARWARS_API"))
	e.App = app.Register(e.Store)
	api.Register(e.App, os.Getenv("PORT"))

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt)
	go func() {
		<-gracefulShutdown
		_ = e.Db.Close()
	}()
}

func (e *server) ConnectDatabase() *sql.DB {
	connectionMysql := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE"))
	db, err := sql.Open("mysql", connectionMysql)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		log.Println("error db.Ping(): ", err.Error())
		panic(err)
	}

	return db
}
