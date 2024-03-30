package main

import (
	"fmt"
	"os"

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
	r.Use(middleware.Logger)

}
