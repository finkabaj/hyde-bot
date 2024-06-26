package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/finkabaj/hyde-bot/internals/backend/controllers"
	"github.com/finkabaj/hyde-bot/internals/backend/services"
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

	logger := logger.NewLogger(fs)

	r := chi.NewRouter()

	if os.Getenv("ENV") == "development" {
		r.Use(middleware.Logger)
	} else if os.Getenv("ENV") == "production" {
		r.Use(middleware.Recoverer)
	}

	var database db.Database = postgresql.NewPostgresql(logger)
	defer database.Close()

	credentials := db.DatabaseCredentials{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
	}

	if err = database.Connect(credentials); err != nil {
		logger.Fatal(err)
	}

	if err = database.Status(); err != nil {
		logger.Fatal(err)
	}

	commandsController := controllers.NewCommandsController(&database)
	commandsController.RegisterRoutes(r)

	guildService := services.NewGuildService(database)
	guildController := controllers.NewGuildController(guildService, logger)
	guildController.RegisterRoutes(r)

	reactionService := services.NewReactionService(logger, database, guildService)
	rulesController := controllers.NewRulesController(reactionService, logger)
	rulesController.RegisterRoutes(r)

	host := os.Getenv("API_HOST")
	port := os.Getenv("API_PORT")

	server := http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  600 * time.Second,
	}

	err = server.ListenAndServe()

	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("API is up and running!")
}
