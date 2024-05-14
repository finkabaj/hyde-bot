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

type GuildController struct {
	service services.IGuildService
	logger  logger.ILogger
}

var guildController *GuildController

func NewGuildController(es services.IGuildService, l logger.ILogger) *GuildController {
	if guildController == nil {
		guildController = &GuildController{
			service: es,
			logger:  l,
		}
	}
	return guildController
}

func (c *GuildController) RegisterRoutes(router *chi.Mux) {
	router.Route("/guild", func(r chi.Router) {
		r.With(middleware.ValidateJson[guild.GuildCreate]()).Post("/", c.postGuild)
		r.Get("/{id}", c.getGuild)
	})
}

func (ec *GuildController) postGuild(w http.ResponseWriter, r *http.Request) {
	g := r.Context().Value(middleware.ValidateJsonCtxKey).(guild.GuildCreate)

	newGuild, err := ec.service.CreateGuild(g)

	switch {
	case err == guild.ErrGuildConflict:
		common.NewErrorResponseBuilder(err).
			SetStatus(http.StatusConflict).
			SetMessage("Guild already exists").
			Send(w)
		return
	case err == common.ErrInternal:
		common.SendInternalError(w)
		return
	case err != nil:
		common.SendInternalError(w, "Unexpected error in postGuild")
		return
	}

	if err := common.MarshalBody(w, http.StatusCreated, &newGuild); err != nil {
		ec.logger.Error(err, map[string]any{"details": "error while marshalling guild info"})
		common.SendInternalError(w)
	}
}

func (ec *GuildController) getGuild(w http.ResponseWriter, r *http.Request) {
	gId := chi.URLParam(r, "id")

	g, err := ec.service.GetGuild(gId)

	switch {
	case err == common.ErrNotFound:
		common.SendNotFoundError(w, fmt.Sprintf("No guild with id: %s found", gId))
		return
	case err == common.ErrInternal:
		common.SendInternalError(w)
		return
	case err != nil:
		common.SendInternalError(w, "Unexpected error at getGuild")
	}

	if err := common.MarshalBody(w, http.StatusOK, &g); err != nil {
		ec.logger.Error(err, map[string]any{"details": "Error while marshaling get guild"})
		common.SendInternalError(w)
	}
}
