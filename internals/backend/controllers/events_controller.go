package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/finkabaj/hyde-bot/internals/backend/services"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/finkabaj/hyde-bot/internals/utils/guild"
)

type EventsController struct {
	service *services.EventsService
}

var eventsController *EventsController

func NewEventsController(es *services.EventsService) *EventsController {
	if eventsController == nil {
		eventsController = &EventsController{
			service: es,
		}
	}
	return eventsController
}

func (c *EventsController) RegisterRoutes(router *chi.Mux) {
	router.Post("/guild", c.postGuild)
	router.Get("/guild/{id}", c.getGuild)
}

func (ec *EventsController) postGuild(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	var guild *guild.GuildCreate

	if err := common.UnmarshalBody(body, &guild); err != nil {
		logger.Error(err, logger.LogFields{"mesage": "error while unmarshaling guild info"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := common.MarshalBody(w, http.StatusCreated, &guild); err != nil {
		logger.Error(err, logger.LogFields{"message": "error while marshalling guild info"})
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (ec *EventsController) getGuild(w http.ResponseWriter, r *http.Request) {
	gId := chi.URLParam(r, "id")

	if gId == "" {
		logger.Debug("Guild ID is empty")
		common.WriteError(w, errors.New("Empty guild id"), http.StatusBadRequest, "Provice guild id to get guild info")
		return
	}

	g, err := ec.service.GetGuild(gId)

	if err != nil {
		logger.Error(err)
		common.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	if g == nil {
		common.WriteError(w, errors.New("No guild found"), http.StatusBadRequest, fmt.Sprintf("No guild with id: %s found", gId))
		return
	}

	if err := common.MarshalBody(w, http.StatusOK, &g); err != nil {
		logger.Error(err, logger.LogFields{"message": "Error while marshaling get guild"})
		common.WriteError(w, err, http.StatusInternalServerError, "Error while marshaling")
	}
}
