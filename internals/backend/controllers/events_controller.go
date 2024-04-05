package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/finkabaj/hyde-bot/internals/db"
)

type EventsController struct {
	DB *db.Database
}

var eventsController *EventsController

func NewEventsController(db *db.Database) *EventsController {
	if eventsController == nil {
		eventsController = &EventsController{}
	}
	return eventsController
}

func (c *EventsController) RegisterRoutes(router *chi.Mux) {
	router.Post("/guild", handleNewGuild)
}

func handleNewGuild(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("OK"))
}
