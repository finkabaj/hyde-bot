package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mogs "github.com/finkabaj/hyde-bot/internals/backend/mocks"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/finkabaj/hyde-bot/internals/utils/rule"
	"github.com/stretchr/testify/assert"
)

var mockReactionService *mogs.MockReactionService = mogs.NewMockReactionService()
var rc *RulesController = NewRulesController(mockReactionService, mogs.NewMockLogger())

const rac = rule.ReactActionCount

func init() {
	rc.RegisterRoutes(r)
}

func TestCreateReactionRules(t *testing.T) {
	t.Run("Positive", testCreateReactionRulePositive)
	t.Run("NegativeConflict", testCreateReactionRuleNegativeConflict)
	t.Run("NegativeBadRequest", testCreateReactionRuleNegativeBadRequest)
	t.Run("NegativeInternalError", testCreateReactionRuleNegativeInternalError)
}

func TestGetReactionsRules(t *testing.T) {
	t.Run("Positive", testGetReactionRulesPositive)
	t.Run("NegativeNotFound", testGetReactionRulesNotFound)
	t.Run("NegativeInternalError", testGetReactionRulesInternalError)
	t.Run("TeapodStatus", testGetReactionRulesTeapot)
}

func TestDeleteReactionRules(t *testing.T) {
	t.Run("Positive", testDeleteReactionRulesPositive)
	t.Run("NegativeNotFound", testDeleteReactionRulesNotFound)
	t.Run("NegativeInternalError", testDeleteReactionRulesInternalError)
	t.Run("BadRequest", testDeleteReactionRulesBadRequest)
}

func testCreateReactionRulePositive(t *testing.T) {
	expectedResponse := []rule.ReactionRule{
		{
			EmojiName:  "🤰",
			IsCustom:   false,
			RuleAuthor: "J3nxJ5WHIoHJinXjSX",
			GuildId:    "QaK6KDIezh0ckrQhySh",
			Actions:    [rac]rule.ReactAction{rule.Delete, rule.Ban},
		},
		{
			EmojiName:  "💦",
			IsCustom:   false,
			RuleAuthor: "J3nxJ5WHIoHJinXjSD",
			GuildId:    "QaK6KDIezh0ckrQhyS",
			Actions:    [rac]rule.ReactAction{rule.Ban},
		},
		{
			EmojiId:    "12321",
			EmojiName:  "bust",
			IsCustom:   true,
			RuleAuthor: "QaK6KDIezh0ckrQhyShD",
			GuildId:    "QaK6KDIezh0ckrQhyS",
			Actions:    [rac]rule.ReactAction{rule.Kick},
		},
	}

	mockReactionService.On("CreateReactionRules", expectedResponse).Return(expectedResponse, nil)

	var byf bytes.Buffer
	json.NewEncoder(&byf).Encode(expectedResponse)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/rules/reaction/", &byf)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualResponse []rule.ReactionRule
	common.UnmarshalBody(rr.Result().Body, &actualResponse)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}

func testCreateReactionRuleNegativeInternalError(t *testing.T) {
	expectedResponse := common.NewErrorResponseBuilder(common.ErrInternal).
		SetStatus(http.StatusInternalServerError).
		Get()
	sendedBody := []rule.ReactionRule{
		{
			EmojiName:  "🤰",
			IsCustom:   false,
			RuleAuthor: "J3nxJ5WHIoHJinXjIE",
			GuildId:    "QaK6KDIezh0ckrQhy",
			Actions:    [rac]rule.ReactAction{rule.Ban},
		},
	}

	mockReactionService.On("CreateReactionRules", sendedBody).Return([]rule.ReactionRule{}, common.ErrInternal)

	var byf bytes.Buffer
	json.NewEncoder(&byf).Encode(sendedBody)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/rules/reaction/", &byf)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualResponse common.ErrorResponse
	common.UnmarshalBody(rr.Result().Body, &actualResponse)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, expectedResponse, &actualResponse)

	mockReactionService.AssertExpectations(t)
}

func testCreateReactionRuleNegativeConflict(t *testing.T) {
	expectedResponse := common.NewErrorResponseBuilder(rule.ErrRuleReactionConflict).
		SetMessage("rule on this reaction already exists").
		SetStatus(http.StatusConflict).
		Get()
	sendedBody := []rule.ReactionRule{
		{
			EmojiName:  "🤰",
			IsCustom:   false,
			RuleAuthor: "J3nxJ5WHIoHJinXjSX",
			GuildId:    "QaK6KDIezh0ckrQhysh",
			Actions:    [rac]rule.ReactAction{rule.Ban},
		},
	}

	mockReactionService.On("CreateReactionRules", sendedBody).Return([]rule.ReactionRule{}, rule.ErrRuleReactionConflict)

	var byf bytes.Buffer
	json.NewEncoder(&byf).Encode(sendedBody)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/rules/reaction/", &byf)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualResponse common.ErrorResponse
	common.UnmarshalBody(rr.Result().Body, &actualResponse)

	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Equal(t, expectedResponse, &actualResponse)

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
			IsCustom:   false,
			RuleAuthor: "J3nxJ5WHIoHJinXjxx",
			GuildId:    "QaK6KDIezh0ckrQhyxx",
			Actions:    [rac]rule.ReactAction{rule.Ban},
		},
	}

	mockReactionService.On("CreateReactionRules", sendedBody).Return([]rule.ReactionRule{}, common.ErrBadRequest)

	var byf bytes.Buffer
	json.NewEncoder(&byf).Encode(sendedBody)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/rules/reaction", &byf)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualResponse common.ErrorResponse
	common.UnmarshalBody(rr.Result().Body, &actualResponse)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, expectedResponse, &actualResponse)

	mockReactionService.AssertExpectations(t)
}

func testGetReactionRulesPositive(t *testing.T) {
	expectedResponse := []rule.ReactionRule{
		{
			EmojiName:  "🤰",
			RuleAuthor: "J3nxJ5WHIoHJinXjSD",
			GuildId:    "QaK6KDIezh0ckrQhy",
			Actions:    [rac]rule.ReactAction{rule.Ban},
		},
		{
			EmojiName:  "💦",
			RuleAuthor: "J3nxJ5WHIoHJinXjSD",
			GuildId:    "QaK6KDIezh0ckrQhy",
			Actions:    [rac]rule.ReactAction{rule.Ban, rule.Kick},
		},
		{
			EmojiId:    "12321",
			RuleAuthor: "QaK6KDIezh0ckrQhyShD",
			GuildId:    "QaK7KDIezh0ckrQhy",
			Actions:    [rac]rule.ReactAction{rule.Ban},
		},
	}

	mockReactionService.On("GetReactionRules", "QaK6KDIezh0ckrQP8").Return(expectedResponse, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/rules/reaction/QaK6KDIezh0ckrQP8", nil)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualResponse []rule.ReactionRule
	common.UnmarshalBodyBytes(rr.Body.Bytes(), &actualResponse)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expectedResponse, actualResponse)

	mockReactionService.AssertExpectations(t)
}

func testGetReactionRulesTeapot(t *testing.T) {
	wtfErr := errors.New("wtf is this")
	expectedResponse := common.NewErrorResponseBuilder(wtfErr).
		SetMessage("something realy bad happened").
		SetStatus(http.StatusTeapot).
		Get()

	mockReactionService.On("GetReactionRules", "QaK6KDIezh0ckrQTe").Return([]rule.ReactionRule{}, wtfErr)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/rules/reaction/QaK6KDIezh0ckrQTe", nil)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualResponse common.ErrorResponse
	common.UnmarshalBody(rr.Result().Body, &actualResponse)

	assert.Equal(t, http.StatusTeapot, rr.Code)
	assert.Equal(t, expectedResponse, &actualResponse)

	mockReactionService.AssertExpectations(t)
}

func testGetReactionRulesNotFound(t *testing.T) {
	expectedResponse := common.NewErrorResponseBuilder(common.ErrNotFound).
		SetMessage("no rules found for the guild").
		SetStatus(http.StatusNotFound).
		Get()

	mockReactionService.On("GetReactionRules", "QaK6KDIezh0ckrQB9").Return([]rule.ReactionRule{}, common.ErrNotFound)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/rules/reaction/QaK6KDIezh0ckrQB9", nil)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualResponse common.ErrorResponse
	common.UnmarshalBody(rr.Result().Body, &actualResponse)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, expectedResponse, &actualResponse)

	mockReactionService.AssertExpectations(t)
}

func testGetReactionRulesInternalError(t *testing.T) {
	expectedResponse := common.NewErrorResponseBuilder(common.ErrInternal).
		SetMessage("Internal server error").
		SetStatus(http.StatusInternalServerError).
		Get()

	mockReactionService.On("GetReactionRules", "QaK6KDIezh0ckrQIE").Return([]rule.ReactionRule{}, common.ErrInternal)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/rules/reaction/QaK6KDIezh0ckrQIE", nil)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualResponse common.ErrorResponse
	common.UnmarshalBody(rr.Result().Body, &actualResponse)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, expectedResponse, &actualResponse)

	mockReactionService.AssertExpectations(t)
}

func testDeleteReactionRulesPositive(t *testing.T) {
	expectedResponse := common.OkResponse{Message: "successfully deleted 3 rules"}
	query := []rule.DeleteReactionRuleQuery{
		{
			EmojiName: "🤰",
		},
		{
			EmojiName: "💦",
		},
		{
			EmojiId: "12321",
		},
	}
	gId := "QaK6KDIezh0ckrQPo"

	encodedQuery := rule.EncodeDeleteReactQuery(query)

	mockReactionService.On("DeleteReactionRules", query, gId).Return(nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/rules/reaction/%s?%s", gId, encodedQuery), nil)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualResponse common.OkResponse
	common.UnmarshalBody(rr.Result().Body, &actualResponse)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expectedResponse, actualResponse)

	mockReactionService.AssertExpectations(t)
}

func testDeleteReactionRulesNotFound(t *testing.T) {
	expectedResponse := common.NewErrorResponseBuilder(common.ErrNotFound).
		SetMessage("no rules found for the guild").
		SetStatus(http.StatusNotFound).
		Get()
	query := []rule.DeleteReactionRuleQuery{
		{
			EmojiName: "f",
		},
	}
	gId := "QaK6KDIezh0ckrQnF"

	encodedQuery := rule.EncodeDeleteReactQuery(query)

	mockReactionService.On("DeleteReactionRules", query, gId).Return(common.ErrNotFound)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/rules/reaction/%s?%s", gId, encodedQuery), nil)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualResponse common.ErrorResponse
	common.UnmarshalBody(rr.Result().Body, &actualResponse)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, expectedResponse, &actualResponse)

	mockReactionService.AssertExpectations(t)
}

func testDeleteReactionRulesInternalError(t *testing.T) {
	expectedResponse := common.NewErrorResponseBuilder(common.ErrInternal).
		SetStatus(http.StatusInternalServerError).
		Get()
	query := []rule.DeleteReactionRuleQuery{
		{
			EmojiId: "123",
		},
	}
	gId := "QaK6KDIezh0ckrQIe"

	encodedQuery := rule.EncodeDeleteReactQuery(query)

	mockReactionService.On("DeleteReactionRules", query, gId).Return(common.ErrInternal)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/rules/reaction/%s?%s", gId, encodedQuery), nil)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualResponse common.ErrorResponse
	common.UnmarshalBody(rr.Result().Body, &actualResponse)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, expectedResponse, &actualResponse)

	mockReactionService.AssertExpectations(t)
}

func testDeleteReactionRulesBadRequest(t *testing.T) {
	expectedResponse := common.NewErrorResponseBuilder(common.ErrBadRequest).
		SetMessage("invalid request query").
		SetStatus(http.StatusBadRequest).
		Get()
	query := []rule.DeleteReactionRuleQuery{
		{
			EmojiId: "123132",
		},
	}
	gId := "QaK6KDIezh0ckrQBR"

	encodedQuery := rule.EncodeDeleteReactQuery(query)

	mockReactionService.On("DeleteReactionRules", query, gId).Return(common.ErrBadRequest)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/rules/reaction/%s?%s", gId, encodedQuery), nil)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualResponse common.ErrorResponse
	common.UnmarshalBody(rr.Result().Body, &actualResponse)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, expectedResponse, &actualResponse)

	mockReactionService.AssertExpectations(t)
}
