package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/finkabaj/hyde-bot/internals/backend/controllers"
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

	controller := controllers.NewCommandsController()
	controller.RegisterRoutes(r)

	host := os.Getenv("API_HOST")
	port := os.Getenv("API_PORT")

	http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r)
}
