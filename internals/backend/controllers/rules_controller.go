package controllers

import (
	"github.com/finkabaj/hyde-bot/internals/backend/services"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/go-chi/chi/v5"
)

type RulesController struct {
	reactionService services.IReactionService
	logger          logger.ILogger
}

var rulesController *RulesController

func NewRulesController(reactionService services.IReactionService, logger logger.ILogger) *RulesController {
	if rulesController == nil {
		rulesController = &RulesController{
			reactionService: reactionService,
			logger:          logger,
		}
	}
	return rulesController
}

func (rc *RulesController) RegisterRoutes(r *chi.Mux) {

}
