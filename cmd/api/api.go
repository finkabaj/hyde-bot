package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/finkabaj/hyde-bot/internals/backend/controllers"
	"github.com/finkabaj/hyde-bot/internals/db"
	"github.com/finkabaj/hyde-bot/internals/db/postgresql"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Printf("Error loading .env: %s", err.Error())
		os.Exit(1)
	}

	fs, err := os.OpenFile("log/api.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		os.Exit(1)
	}

	logger.Init(fs)

	r := chi.NewRouter()

	if os.Getenv("ENV") == "development" {
		r.Use(middleware.Logger)
	} else if os.Getenv("ENV") == "production" {
		r.Use(middleware.Recoverer)
	}

	var database db.Database = &postgresql.Postgresql{}

	credentials := db.DatabaseCredentials{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
	}

	if err = database.Connect(&credentials); err != nil {
		logger.Fatal(err)
	}

	if err = database.Status(); err != nil {
		logger.Fatal(err)
	}

	defer database.Close()

	controller := controllers.NewCommandsController(&database)
	controller.RegisterRoutes(r)

	host := os.Getenv("API_HOST")
	port := os.Getenv("API_PORT")

	http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r)
}
