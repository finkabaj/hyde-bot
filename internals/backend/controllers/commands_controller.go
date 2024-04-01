package controllers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type CommandsController struct {
}

var commandsController *CommandsController

// NewCommandsController is a function that returns a new commands controller or existing commands controller
func NewCommandsController() *CommandsController {
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
