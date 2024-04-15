package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mogs "github.com/finkabaj/hyde-bot/internals/backend/mocks"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/finkabaj/hyde-bot/internals/utils/rule"
	"github.com/stretchr/testify/assert"
)

// TODO: After tests are written, initialize the controller and service
var mockReactionService *mogs.MockReactionService = mogs.NewMockReactionService()
var rc *RulesController = NewRulesController(mockReactionService, mogs.NewMockLogger())

func init() {
	rc.RegisterRoutes(r)
}

func TestCreateReactionRules(t *testing.T) {
	t.Run("Positive", testCreateReactionRulePositive)
	t.Run("NegativeConflict", testCreateReactionRuleNegativeConflict)
	t.Run("NegativeIncompatible", testCreateReactionRuleNegativeIncompatible)
	t.Run("NegativeBadRequest", testCreateReactionRuleNegativeBadRequest)
}

func TestGetReactionsRules(t *testing.T) {
	t.Run("Positive", testGetReactionRulesPositive)
	t.Run("NegativeNotFound", testGetReactionRulesNotFound)
	t.Run("NegativeInternalError", testGetReactionRulesInternalError)
}

func TestDeleteReactionRules(t *testing.T) {

}

func testCreateReactionRulePositive(t *testing.T) {
	expectedResponse := []rule.ReactionRule{
		{
			EmojiName:  "🤰",
			RuleAuthor: "J3nxJ5WHIoHJinXjSX",
			GuildId:    "QaK6KDIezh0ckrQhySh",
			Action:     rule.Delete,
		},
		{
			EmojiName:  "💦",
			RuleAuthor: "J3nxJ5WHIoHJinXjSD",
			GuildId:    "QaK6KDIezh0ckrQhyS",
			Action:     rule.Ban,
		},
		{
			EmojiId:    "12321",
			RuleAuthor: "QaK6KDIezh0ckrQhyShD",
			GuildId:    "QaK6KDIezh0ckrQhyS",
			Action:     rule.Kick,
		},
	}

	mockReactionService.On("CreateReactionRules", &expectedResponse).Return(&expectedResponse, nil)

	var byf bytes.Buffer
	json.NewEncoder(&byf).Encode(expectedResponse)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "rules/reaction", &byf)
	r.ServeHTTP(rr, req)

	var actualResponse []rule.ReactionRule
	common.UnmarshalBody(rr.Result().Body, actualResponse)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func testCreateReactionRuleNegativeConflict(t *testing.T) {
	expectedResponse := common.NewErrorResponseBuilder(rule.ErrRuleReactionConflict).
		SetMessage("rule on this reaction already exists").
		SetStatus(http.StatusConflict).
		Get()
	sendedBody := []rule.ReactionRule{
		{
			EmojiName:  "🤰",
			RuleAuthor: "J3nxJ5WHIoHJinXjSX",
			GuildId:    "QaK6KDIezh0ckrQhysh",
			Action:     rule.Delete,
		},
	}

	mockReactionService.On("CreateReactionRules", &sendedBody).Return(nil, rule.ErrRuleReactionConflict)

	var byf bytes.Buffer
	json.NewEncoder(&byf).Encode(sendedBody)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "rules/reaction", &byf)
	r.ServeHTTP(rr, req)

	var actualResponse common.ErrorResponse
	common.UnmarshalBody(rr.Result().Body, actualResponse)

	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Equal(t, expectedResponse, actualResponse)

	mockReactionService.AssertExpectations(t)
}

func testCreateReactionRuleNegativeIncompatible(t *testing.T) {
	expectedResponse := common.NewErrorResponseBuilder(rule.ErrRuleReactionIncompatible).
		SetMessage("either emoji name or emoji id must be provided").
		SetStatus(http.StatusConflict).
		Get()
	sendedBody := []rule.ReactionRule{
		{
			EmojiName:  "🤰",
			EmojiId:    "asdsad",
			RuleAuthor: "J3nxJ5WHIoHJinXjSX",
			GuildId:    "QaK6KDIezh0ckrQhysh",
			Action:     rule.Delete,
		},
	}

	mockReactionService.On("CreateReactionRules", &sendedBody).Return(nil, rule.ErrRuleReactionIncompatible)

	var byf bytes.Buffer
	json.NewEncoder(&byf).Encode(expectedResponse)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "rules/reaction", &byf)
	r.ServeHTTP(rr, req)

	var actualResponse common.ErrorResponse
	common.UnmarshalBody(rr.Result().Body, actualResponse)

	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Equal(t, expectedResponse, actualResponse)

	mockReactionService.AssertExpectations(t)
}

func testCreateReactionRuleNegativeBadRequest(t *testing.T) {
	expectedResponse := common.NewErrorResponseBuilder(common.ErrBadRequest).
		SetMessage("invalid request body").
		SetStatus(http.StatusBadRequest).
		Get()
	sendedBody := []rule.ReactionRule{
		{
			EmojiName:  "🤰",
			RuleAuthor: "J3nxJ5WHIoHJinXjxx",
			GuildId:    "QaK6KDIezh0ckrQhyxx",
			Action:     rule.Delete,
		},
	}

	mockReactionService.On("CreateReactionRules", &sendedBody).Return(nil, common.ErrBadRequest)

	var byf bytes.Buffer
	json.NewEncoder(&byf).Encode(sendedBody)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "rules/reaction", &byf)
	r.ServeHTTP(rr, req)

	var actualResponse common.ErrorResponse
	common.UnmarshalBody(rr.Result().Body, actualResponse)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, expectedResponse, actualResponse)

	mockReactionService.AssertExpectations(t)
}

func testGetReactionRulesPositive(t *testing.T) {
	expectedResponse := []rule.ReactionRule{
		{
			EmojiName:  "🤰",
			RuleAuthor: "J3nxJ5WHIoHJinXjSD",
			GuildId:    "QaK6KDIezh0ckrQhy",
			Action:     rule.Delete,
		},
		{
			EmojiName:  "💦",
			RuleAuthor: "J3nxJ5WHIoHJinXjSD",
			GuildId:    "QaK6KDIezh0ckrQhy",
			Action:     rule.Ban,
		},
		{
			EmojiId:    "12321",
			RuleAuthor: "QaK6KDIezh0ckrQhyShD",
			GuildId:    "QaK7KDIezh0ckrQhy",
			Action:     rule.Kick,
		},
	}

	mockReactionService.On("GetReactionRules", "QaK6KDIezh0ckrQhy").Return(&expectedResponse, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "rules/reaction/QaK6KDIezh0ckrQhy", nil)
	r.ServeHTTP(rr, req)

	var actualResponse []rule.ReactionRule
	common.UnmarshalBody(rr.Result().Body, actualResponse)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expectedResponse, actualResponse)

	mockReactionService.AssertExpectations(t)
}

func testGetReactionRulesNotFound(t *testing.T) {
	expectedResponse := common.NewErrorResponseBuilder(common.ErrNotFound).
		SetMessage("no rules found for the guild").
		SetStatus(http.StatusNotFound).
		Get()

	mockReactionService.On("GetReactionRules", "QaK6KDIezh0ckrQhd").Return(nil, common.ErrNotFound)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "rules/reaction/QaK6KDIezh0ckrQhy", nil)
	r.ServeHTTP(rr, req)

	var actualResponse common.ErrorResponse
	common.UnmarshalBody(rr.Result().Body, actualResponse)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, expectedResponse, actualResponse)

	mockReactionService.AssertExpectations(t)
}

func testGetReactionRulesInternalError(t *testing.T) {
	expectedResponse := common.NewErrorResponseBuilder(common.ErrInternal).
		SetMessage("internal server error").
		SetStatus(http.StatusInternalServerError).
		Get()

	mockReactionService.On("GetReactionRules", "QaK6KDIezh0ckrQIE").Return(nil, common.ErrInternal)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "rules/reaction/QaK6KDIezh0ckrQIE", nil)
	r.ServeHTTP(rr, req)

	var actualResponse common.ErrorResponse
	common.UnmarshalBody(rr.Result().Body, actualResponse)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, expectedResponse, actualResponse)

	mockReactionService.AssertExpectations(t)
}
