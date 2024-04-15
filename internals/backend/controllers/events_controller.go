package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/finkabaj/hyde-bot/internals/backend/middleware"
	"github.com/finkabaj/hyde-bot/internals/backend/services"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/finkabaj/hyde-bot/internals/utils/guild"
)

type EventsController struct {
	service services.IEventsService
	logger  logger.ILogger
}

var eventsController *EventsController

func NewEventsController(es services.IEventsService, l logger.ILogger) *EventsController {
	if eventsController == nil {
		eventsController = &EventsController{
			service: es,
			logger:  l,
		}
	}
	return eventsController
}

func (c *EventsController) RegisterRoutes(router *chi.Mux) {
	router.Route("/guild", func(r chi.Router) {
		r.With(middleware.ValidateJson[guild.GuildCreate]()).Post("/", c.postGuild)
		r.Get("/{id}", c.getGuild)
	})
}

func (ec *EventsController) postGuild(w http.ResponseWriter, r *http.Request) {
	g := r.Context().Value(middleware.ValidateJsonCtxKey).(guild.GuildCreate)

	newGuild, err := ec.service.CreateGuild(&g)

	if err == guild.ErrGuildConflict {
		common.NewErrorResponseBuilder(err).
			SetStatus(http.StatusConflict).
			SetMessage(fmt.Sprintf("Guild with id: %s already exists", g.GuildId)).
			Send(w)
		return
	} else if err != nil {
		ec.logger.Error(err, logger.LogFields{"message": "error while creating new guild"})
		common.NewErrorResponseBuilder(err).
			SetStatus(http.StatusInternalServerError).
			Send(w)
		return
	}

	if err := common.MarshalBody(w, http.StatusCreated, &newGuild); err != nil {
		ec.logger.Error(err, logger.LogFields{"message": "error while marshalling guild info"})
		common.NewErrorResponseBuilder(err).
			SetStatus(http.StatusInternalServerError).
			Send(w)
	}
}

func (ec *EventsController) getGuild(w http.ResponseWriter, r *http.Request) {
	gId := chi.URLParam(r, "id")

	g, err := ec.service.GetGuild(gId)

	if err != nil {
		ec.logger.Error(err)
		common.NewErrorResponseBuilder(err).
			SetStatus(http.StatusInternalServerError).
			Send(w)
		return
	}

	if g == nil {
		common.NewErrorResponseBuilder(common.ErrNotFound).
			SetStatus(http.StatusNotFound).
			SetMessage(fmt.Sprintf("No guild with id: %s found", gId)).
			Send(w)
		return
	}

	if err := common.MarshalBody(w, http.StatusOK, &g); err != nil {
		ec.logger.Error(err, logger.LogFields{"message": "Error while marshaling get guild"})
		common.NewErrorResponseBuilder(err).
			SetStatus(http.StatusInternalServerError).
			Send(w)
	}
}
