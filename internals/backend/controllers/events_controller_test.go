package controllers

import (
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

var mockService *mogs.MockEventsService = mogs.NewMockEventsService()
var r *chi.Mux = chi.NewRouter()
var ec *EventsController = NewEventsController(mockService, mogs.NewMockLogger())

func TestGetGuild(t *testing.T) {
	ec.RegisterRoutes(r)
	t.Run("Positive", testGetGuildPositive)
	t.Run("NegativeNotFound", testGetGuildNegativeNotFound)
	t.Run("NegativeInternalError", testGetGuildNegativeInternalError)
}

func testGetGuildPositive(t *testing.T) {
	gId := "positive"
	expectedResponse := guild.Guild{
		GuildId: "positive",
		OwnerId: "positive",
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/guild/%s", gId), nil)
	r.ServeHTTP(rr, req)

	var actualRespose guild.Guild
	common.UnmarshalBody(rr.Result().Body, &actualRespose)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, expectedResponse, actualRespose)

	mockService.AssertExpectations(t)
}

func testGetGuildNegativeNotFound(t *testing.T) {
	gId := "negativeNotFound"
	expectedResponse := common.NewErrorResponseBuilder(guild.ErrGuildNotFound).
		SetStatus(http.StatusNotFound).
		SetMessage(fmt.Sprintf("No guild with id: %s found", gId)).
		Get()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/guild/%s", gId), nil)
	r.ServeHTTP(rr, req)

	var actualRespose *common.ErrorResponse
	if err := common.UnmarshalBody(rr.Result().Body, &actualRespose); err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, rr.Code, http.StatusNotFound)
	assert.Equal(t, expectedResponse, actualRespose)

	mockService.AssertExpectations(t)
}

func testGetGuildNegativeInternalError(t *testing.T) {
	gId := "negativeInternalError"
	expectedError := errors.New("Internal error")
	expectedResponse := common.NewErrorResponseBuilder(expectedError).
		SetStatus(http.StatusInternalServerError).
		Get()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/guild/%s", gId), nil)
	r.ServeHTTP(rr, req)

	var actualRespose *common.ErrorResponse
	if err := common.UnmarshalBody(rr.Result().Body, &actualRespose); err != nil {
		fmt.Println(err)
	}

	assert.Equal(t, rr.Code, http.StatusInternalServerError)
	assert.Equal(t, expectedResponse, actualRespose)

	mockService.AssertExpectations(t)
}
