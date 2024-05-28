package controllers

import (
	"fmt"
	"net/http"

	"github.com/finkabaj/hyde-bot/internals/backend/middleware"
	"github.com/finkabaj/hyde-bot/internals/backend/services"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/finkabaj/hyde-bot/internals/utils/rule"
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
			r.With(middleware.ValidateJson[[]rule.ReactionRule]()).Post("/", rc.postReactions)
			r.With(middleware.ValidateQuery(rule.DecodeDeleteReactQuery)).Delete("/{id}", rc.deleteReactions)
		})
	})
}

func (rc *RulesController) getReactions(w http.ResponseWriter, r *http.Request) {
	gId := chi.URLParam(r, "id")

	rules, err := rc.reactionService.GetReactionRules(gId)

	switch {
	case err == common.ErrInternal:
		common.SendInternalError(w, "Internal server error")
		return
	case err == common.ErrNotFound:
		common.SendNotFoundError(w, "no rules found for the guild")
		return
	case err != nil:
		common.NewErrorResponseBuilder(err).
			SetMessage("something realy bad happened").
			SetStatus(http.StatusTeapot).
			Send(w)
		return
	}

	if err := common.MarshalBody(w, http.StatusOK, &rules); err != nil {
		rc.logger.Error(err, map[string]any{"details": "Error while marhsaling getReactionsRules response"})
		common.SendInternalError(w)
	}
}

func (rc *RulesController) postReactions(w http.ResponseWriter, r *http.Request) {
	rRules, ok := middleware.JsonFromContext(r.Context()).([]rule.ReactionRule)

	if !ok {
		rc.logger.Error(common.ErrInternal, map[string]any{"details": "error while validating postReactions"})
		common.SendInternalError(w, "error while validating")
		return
	}

	newRules, err := rc.reactionService.CreateReactionRules(rRules)

	switch err {
	case common.ErrInternal:
		common.SendInternalError(w)
		return
	case rule.ErrRuleReactionConflict:
		common.NewErrorResponseBuilder(err).
			SetStatus(http.StatusConflict).
			SetMessage("rule on this reaction already exists").
			Send(w)
		return
	case common.ErrBadRequest:
		common.SendBadRequestError(w, "invalid request body")
		return
	case common.ErrNotFound:
		common.SendNotFoundError(w, "provided guild not found")
		return
	}

	if err := common.MarshalBody(w, http.StatusCreated, &newRules); err != nil {
		rc.logger.Error(common.ErrInternal, map[string]any{"details": "error while marshaling postReactions"})
		common.NewErrorResponseBuilder(err).
			SetMessage("Error while marshing response").
			SetStatus(http.StatusInternalServerError).
			Send(w)
	}
}

func (rc *RulesController) deleteReactions(w http.ResponseWriter, r *http.Request) {
	gId := chi.URLParam(r, "id")
	query, ok := middleware.QueryFromContext(r.Context()).([]rule.DeleteReactionRuleQuery)

	if !ok {
		rc.logger.Error(common.ErrInternal, map[string]any{"details": "no value found in context"})
		common.SendInternalError(w)
		return
	}

	err := rc.reactionService.DeleteReactionRules(query, gId)

	switch err {
	case common.ErrNotFound:
		common.SendNotFoundError(w, "no rules found for the guild")
		return
	case common.ErrInternal:
		common.SendInternalError(w)
		return
	case common.ErrBadRequest:
		common.SendBadRequestError(w, "invalid request query")
		return
	}

	ruleWord := "rule"
	if len(query) > 1 {
		ruleWord = "rules"
	}
	res := common.OkResponse{Message: fmt.Sprintf("successfully deleted %d %s", len(query), ruleWord)}

	if err := common.MarshalBody(w, http.StatusOK, &res); err != nil {
		rc.logger.Error(err, map[string]any{"details": "error while marhsalong deleteReactions response"})
		common.SendInternalError(w, "error while marshaling deleteReactions response")
	}
}
