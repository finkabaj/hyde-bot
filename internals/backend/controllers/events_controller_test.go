package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	mogs "github.com/finkabaj/hyde-bot/internals/backend/mocks"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/finkabaj/hyde-bot/internals/utils/guild"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

var mockService *mogs.MockEventsService
var r *chi.Mux
var rr *httptest.ResponseRecorder
var ec *EventsController

func TestMain(m *testing.M) {
	mockService = mogs.NewMockEventsService()
	r = chi.NewRouter()
	rr = httptest.NewRecorder()
	ec = NewEventsController(mockService)
	ec.RegisterRoutes(r)

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestGetGuildPositive(t *testing.T) {
	gId := "1c"
	expectedResponse := guild.Guild{
		GuildId: gId,
		OwnerId: "one ass",
	}
	var expectedError error = nil

	mockService.On("GetGuild", gId).Return(expectedResponse, expectedError)

	req := httptest.NewRequest("GET", fmt.Sprintf("/guild/%s", gId), nil)
	r.ServeHTTP(rr, req)

	var actualRespose guild.Guild
	common.UnmarshalBody(rr.Result().Body, &actualRespose)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, expectedResponse, actualRespose)

	mockService.AssertExpectations(t)
}
