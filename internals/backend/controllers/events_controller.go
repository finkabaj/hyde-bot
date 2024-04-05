package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/finkabaj/hyde-bot/internals/db"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/finkabaj/hyde-bot/internals/utils/guild"
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
	body := r.Body
	defer body.Close()

	var guild *guild.GuildCreate

	if err := common.UnmarshalResponse(body, &guild); err != nil {
		logger.Error(err, logger.LogFields{"mesage": "error while unmarshaling guild info"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := common.MarshalRequest(w, &guild); err != nil {
		logger.Error(err, logger.LogFields{"message": "error while marshalling guild info"})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
