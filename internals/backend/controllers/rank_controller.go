package controllers

import (
	"net/http"

	"github.com/finkabaj/hyde-bot/internals/backend/middleware"
	"github.com/finkabaj/hyde-bot/internals/backend/services"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/ranks"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/go-chi/chi/v5"
)

type RankController struct {
	rankService services.IRankService
	logger      logger.ILogger
}

var rankController *RankController

func NewRankController(rs services.IRankService, l logger.ILogger) *RankController {
	if rankController == nil {
		rankController = &RankController{
			rankService: rs,
			logger:      l,
		}
	}
	return rankController
}

func (rc *RankController) RegisterRoutes(r *chi.Mux) {
	r.Route("/rank", func(r chi.Router) {
		r.Get("/{gID}", rc.getRanks)
		r.With(middleware.ValidateJson[ranks.Ranks]()).Post("/", rc.postRanks)
		r.Delete("/{gID}", rc.deleteRanks)
		r.Delete("/{gID}/{rID}", rc.deleteRank)
	})
}

func (rc *RankController) getRanks(w http.ResponseWriter, r *http.Request) {

}

func (rc *RankController) postRanks(w http.ResponseWriter, r *http.Request) {
	ranksFromCtx, ok := middleware.JsonFromContext(r.Context()).(ranks.Ranks)

	if !ok {
		rc.logger.Error(common.ErrInternal, map[string]any{"details": "error while validating postRanks"})
		common.SendInternalError(w, "error while validating")
		return
	}

	newRanks, err := rc.rankService.CreateRanks(ranksFromCtx)

	switch {
	case err == common.ErrBadRequest:
		common.SendBadRequestError(w, "invalid request")
		return
	case err == common.ErrInternal:
		common.SendInternalError(w, "internal error")
		return
	case err == common.ErrNotFound:
		common.SendNotFoundError(w, "guild id or owner id not found")
		return
	case err != nil:
		common.NewErrorResponseBuilder(err).SetStatus(http.StatusInternalServerError).SetMessage("internal error").Send(w)
		return
	}

	if err := common.MarshalBody(w, http.StatusCreated, &newRanks); err != nil {
		rc.logger.Error(err, map[string]any{"details": "error while marshaling response in postRanks"})
		common.SendInternalError(w, err.Error())
	}
}

func (rc *RankController) deleteRanks(w http.ResponseWriter, r *http.Request) {

}

func (rc *RankController) deleteRank(w http.ResponseWriter, r *http.Request) {

}
