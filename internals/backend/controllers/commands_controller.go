package controllers

import (
	"github.com/finkabaj/hyde-bot/internals/db"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type CommandsController struct {
	DB *db.Database
}

var commandsController *CommandsController

func NewCommandsController(db *db.Database) *CommandsController {
	if commandsController == nil {
		commandsController = &CommandsController{}
	}
	return commandsController
}

func (c *CommandsController) RegisterRoutes(router *chi.Mux) {
	router.Get("/", handleGetRoot)
}

func handleGetRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
