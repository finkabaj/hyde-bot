package controllers

import (
	"net/http"

	"github.com/finkabaj/hyde-bot/internals/backend/services"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
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
	r.Route("/rules", func(r chi.Router) {
		r.Route("/reaction", func(r chi.Router) {
			r.Get("/{id}", rc.getReactions)
			r.Post("/", rc.postReactions)
			r.Delete("/{id}", rc.deleteReactions)
		})
	})
}

func (rc *RulesController) getReactions(w http.ResponseWriter, r *http.Request) {
	gId := chi.URLParam(r, "id")

	rules, err := rc.reactionService.GetReactionRules(gId)

	switch {
	case err == common.ErrInternal:
		common.NewErrorResponseBuilder(err).
			SetMessage("Internal server error").
			SetStatus(http.StatusInternalServerError).
			Send(w)
		return
	case err == common.ErrNotFound:
		common.NewErrorResponseBuilder(err).
			SetMessage("no rules found for the guild").
			SetStatus(http.StatusNotFound).
			Send(w)
		return
	case err != nil:
		common.NewErrorResponseBuilder(err).
			SetMessage("something realy bad happened").
			SetStatus(http.StatusTeapot).
			Send(w)
		return
	}

	if err := common.MarshalBody(w, http.StatusOK, rules); err != nil {
		rc.logger.Error(err, logger.LogFields{"message": "Error while marhsaling getReactionsRules response"})
		common.NewErrorResponseBuilder(common.ErrInternal).
			SetMessage("Internal server error").
			SetStatus(http.StatusInternalServerError).
			Send(w)
	}
}

func (rc *RulesController) postReactions(w http.ResponseWriter, r *http.Request) {

}

func (rc *RulesController) deleteReactions(w http.ResponseWriter, r *http.Request) {

}
