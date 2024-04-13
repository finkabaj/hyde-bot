package controllers

import (
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
var ec *EventsController = NewEventsController(mockService)

func init() {
	ec.RegisterRoutes(r)
}

func TestGetGuildPositive(t *testing.T) {
	gId := "1c"
	expectedResponse := guild.Guild{
		GuildId: gId,
		OwnerId: "one ass",
	}
	var expectedError error = nil

	mockService.On("GetGuild", gId).Return(&expectedResponse, expectedError)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", fmt.Sprintf("/guild/%s", gId), nil)
	r.ServeHTTP(rr, req)

	var actualRespose guild.Guild
	common.UnmarshalBody(rr.Result().Body, &actualRespose)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, expectedResponse, actualRespose)

	mockService.AssertExpectations(t)
}

func TestGetGuildNegative(t *testing.T) {
	gId := "1ass"
	expectedResponse := common.NewErrorResponseBuilder(guild.ErrGuildNotFound).
		SetStatus(http.StatusNotFound).
		SetMessage(fmt.Sprintf("No guild with id: %s found", gId)).
		Get()

	mockService.On("GetGuild", gId).Return(nil, nil)

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
