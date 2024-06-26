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
	"github.com/finkabaj/hyde-bot/internals/utils/guild"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

var mockGuildService *mogs.MockGuildService = mogs.NewMockGuildService()
var r *chi.Mux = chi.NewRouter()
var ec *GuildController = NewGuildController(mockGuildService, mogs.NewMockLogger())

func init() {
	ec.RegisterRoutes(r)
}

func TestGetGuild(t *testing.T) {
	t.Run("Positive", testGetGuildPositive)
	t.Run("NegativeNotFound", testGetGuildNegativeNotFound)
	t.Run("NegativeInternalError", testGetGuildNegativeInternalError)
	t.Run("NegativeWtf", testGetGuildNegativeWtf)
}

func TestCreateGuild(t *testing.T) {
	t.Run("Positive", testCreateGuildPositive)
	t.Run("NegativeValidationNil", testCreateGuildNegativeValidationNil)
	t.Run("NegativeValidation", testCreateGuildNegativeValidation)
	t.Run("NegativeConflict", testCreateGuildNegativeConflict)
}

func testGetGuildPositive(t *testing.T) {
	gId := "positive"
	expectedResponse := guild.Guild{
		GuildId: "positive",
		OwnerId: "positive",
	}

	mockGuildService.On("GetGuild", gId).Return(expectedResponse, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/guild/%s", gId), nil)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualRespose guild.Guild
	common.UnmarshalBody(rr.Result().Body, &actualRespose)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, expectedResponse, actualRespose)

	mockGuildService.AssertExpectations(t)
}

func testGetGuildNegativeNotFound(t *testing.T) {
	gId := "negativeNotFound"
	expectedResponse := common.NewErrorResponseBuilder(common.ErrNotFound).
		SetStatus(http.StatusNotFound).
		SetMessage(fmt.Sprintf("No guild with id: %s found", gId)).
		Get()

	mockGuildService.On("GetGuild", gId).Return(guild.Guild{}, common.ErrNotFound)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/guild/%s", gId), nil)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualRespose *common.ErrorResponse
	if err := common.UnmarshalBody(rr.Result().Body, &actualRespose); err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, expectedResponse, actualRespose)

	mockGuildService.AssertExpectations(t)
}

func testGetGuildNegativeInternalError(t *testing.T) {
	gId := "negativeInternalError"
	expectedResponse := common.NewErrorResponseBuilder(common.ErrInternal).
		SetStatus(http.StatusInternalServerError).
		Get()

	mockGuildService.On("GetGuild", gId).Return(guild.Guild{}, common.ErrInternal)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/guild/%s", gId), nil)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualRespose *common.ErrorResponse
	if err := common.UnmarshalBody(rr.Result().Body, &actualRespose); err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, expectedResponse, actualRespose)

	mockGuildService.AssertExpectations(t)
}

func testGetGuildNegativeWtf(t *testing.T) {
	gId := "negativeWtf"
	expectedError := errors.New("wtf")
	expectedResponse := common.NewErrorResponseBuilder(common.ErrInternal).
		SetStatus(http.StatusInternalServerError).
		SetMessage("Unexpected error at getGuild").
		Get()

	mockGuildService.On("GetGuild", gId).Return(guild.Guild{
		GuildId: "posi vibes",
		OwnerId: "spinnu~",
	}, expectedError)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/guild/%s", gId), nil)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualRespose *common.ErrorResponse
	if err := common.UnmarshalBody(rr.Result().Body, &actualRespose); err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, expectedResponse, actualRespose)

	mockGuildService.AssertExpectations(t)
}

func testCreateGuildPositive(t *testing.T) {
	expectedResponse := guild.Guild{
		GuildId: "QaK6KDIezh0ckrQhySh",
		OwnerId: "J3nxJ5WHIoHJinXjSX",
	}

	mockGuildService.On("CreateGuild", guild.GuildCreate{GuildId: "QaK6KDIezh0ckrQhySh", OwnerId: "J3nxJ5WHIoHJinXjSX"}).Return(expectedResponse, nil)

	var byf bytes.Buffer
	json.NewEncoder(&byf).Encode(expectedResponse)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/guild", &byf)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualRespose guild.Guild
	common.UnmarshalBody(rr.Result().Body, &actualRespose)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, expectedResponse, actualRespose)

	mockGuildService.AssertExpectations(t)
}

func testCreateGuildNegativeValidationNil(t *testing.T) {
	expectedResponse := common.NewErrorResponseBuilder(common.ErrEmptyBody).
		SetStatus(http.StatusBadRequest).
		SetMessage("empty request body").
		Get()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/guild", nil)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualRespose *common.ErrorResponse
	if err := common.UnmarshalBody(rr.Result().Body, &actualRespose); err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, expectedResponse, actualRespose)
}

func testCreateGuildNegativeValidation(t *testing.T) {
	expectedResponse := common.NewErrorResponseBuilder(common.ErrValidation).
		SetStatus(http.StatusBadRequest).
		SetValidationFields(map[string]string{"guildId": "required"}).
		Get()
	sendedBody := guild.GuildCreate{OwnerId: "ass"}

	var byf bytes.Buffer

	json.NewEncoder(&byf).Encode(sendedBody)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/guild", &byf)
	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualRespose *common.ErrorResponse
	if err := common.UnmarshalBody(rr.Result().Body, &actualRespose); err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, expectedResponse, actualRespose)

	mockGuildService.AssertExpectations(t)
}

func testCreateGuildNegativeConflict(t *testing.T) {
	expectedResponse := common.NewErrorResponseBuilder(guild.ErrGuildConflict).
		SetStatus(http.StatusConflict).
		SetMessage("Guild already exists").
		Get()

	sendedBody := guild.GuildCreate{GuildId: "SAS6KDIezh0ckrQhySh", OwnerId: "COCxJ5WHIoHJinXjSX"}
	var byf bytes.Buffer
	json.NewEncoder(&byf).Encode(sendedBody)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/guild", &byf)

	mockGuildService.On("CreateGuild", sendedBody).Return(guild.Guild{}, guild.ErrGuildConflict)

	r.ServeHTTP(rr, req)
	defer rr.Result().Body.Close()

	var actualRespose *common.ErrorResponse
	if err := common.UnmarshalBody(rr.Result().Body, &actualRespose); err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Equal(t, expectedResponse, actualRespose)

	mockGuildService.AssertExpectations(t)
}
