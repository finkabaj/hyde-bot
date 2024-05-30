package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mogs "github.com/finkabaj/hyde-bot/internals/backend/mocks"
	"github.com/finkabaj/hyde-bot/internals/ranks"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var mockRankService *mogs.MockRankService = &mogs.MockRankService{}
var mockRankController *RankController = NewRankController(mockRankService, mogs.NewMockLogger())

func init() {
	mockRankController.RegisterRoutes(r)
}

func CreateMockRanks() ranks.Ranks {
	gID := uuid.NewString()
	oID := uuid.NewString()
	return ranks.Ranks{
		Ranks: []ranks.Rank{
			{
				XP:    100,
				ID:    uuid.NewString(),
				Role:  nil,
				Level: 1,
			},
			{
				XP:    200,
				ID:    uuid.NewString(),
				Role:  nil,
				Level: 2,
			},
		},
		GuildID: gID,
		OwnerID: oID,
	}
}

func TestCreateRanks(t *testing.T) {
	d1 := CreateMockRanks()
	d2 := CreateMockRanks()
	d3 := CreateMockRanks()
	d4 := CreateMockRanks()
	d5 := CreateMockRanks()
	data := []struct {
		name           string
		input          ranks.Ranks
		expected       ranks.Ranks
		expectedStatus int
		expectedError  common.ErrorResponse
		rawError       error
	}{
		{
			name:           "Positive",
			input:          d1,
			expected:       d1,
			expectedError:  common.ErrorResponse{},
			expectedStatus: http.StatusCreated,
			rawError:       nil,
		},
		{
			name:           "BadRequest",
			input:          d2,
			expected:       ranks.Ranks{},
			expectedError:  *common.NewErrorResponseBuilder(common.ErrBadRequest).SetStatus(http.StatusBadRequest).SetMessage("invalid request").Get(),
			expectedStatus: http.StatusBadRequest,
			rawError:       common.ErrBadRequest,
		},
		{
			name:           "Internal",
			input:          d3,
			expected:       ranks.Ranks{},
			expectedError:  *common.NewErrorResponseBuilder(common.ErrInternal).SetStatus(http.StatusInternalServerError).SetMessage("internal error").Get(),
			expectedStatus: http.StatusInternalServerError,
			rawError:       common.ErrInternal,
		},
		{
			name:           "NotFound",
			input:          d4,
			expected:       ranks.Ranks{},
			expectedError:  *common.NewErrorResponseBuilder(common.ErrNotFound).SetStatus(http.StatusNotFound).SetMessage("guild id or owner id not found").Get(),
			expectedStatus: http.StatusNotFound,
			rawError:       common.ErrNotFound,
		},
		{
			name:           "Unexpected",
			input:          d5,
			expected:       ranks.Ranks{},
			expectedError:  *common.NewErrorResponseBuilder(errors.New("tnd")).SetStatus(http.StatusInternalServerError).SetMessage("internal error").Get(),
			expectedStatus: http.StatusInternalServerError,
			rawError:       errors.New("tnd"),
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			mockRankService.On("CreateRanks", d.input).Return(d.expected, d.rawError)

			var byf bytes.Buffer
			json.NewEncoder(&byf).Encode(d.input)

			rr := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/rank/", &byf)
			r.ServeHTTP(rr, req)
			defer rr.Result().Body.Close()

			assert.Equal(t, d.expectedStatus, rr.Code)

			if reflect.DeepEqual(d.expectedError, common.ErrorResponse{}) {
				var aResp ranks.Ranks
				common.UnmarshalBody(rr.Result().Body, &aResp)
				assert.Equal(t, d.expected, aResp)
			} else {
				var aResp common.ErrorResponse
				common.UnmarshalBody(rr.Result().Body, &aResp)
				assert.Equal(t, d.expectedError, aResp)
			}

			mockRankService.AssertExpectations(t)
		})
	}
}

func TestGetRanks(t *testing.T) {
	d1 := CreateMockRanks()
	data := []struct {
		name           string
		inputGuildID   string
		expectedRanks  ranks.Ranks
		expectedStatus int
		expectedError  common.ErrorResponse
		rawError       error
	}{
		{
			name:           "Positive",
			inputGuildID:   d1.GuildID,
			expectedRanks:  d1,
			expectedError:  common.ErrorResponse{},
			expectedStatus: http.StatusOK,
			rawError:       nil,
		},
		{
			name:           "NotFound",
			inputGuildID:   "2",
			expectedRanks:  ranks.Ranks{},
			expectedError:  *common.NewErrorResponseBuilder(common.ErrNotFound).SetStatus(http.StatusNotFound).SetMessage("guild id not found").Get(),
			expectedStatus: http.StatusNotFound,
			rawError:       common.ErrNotFound,
		},
		{
			name:           "Internal",
			inputGuildID:   "3",
			expectedRanks:  ranks.Ranks{},
			expectedError:  *common.NewErrorResponseBuilder(common.ErrInternal).SetStatus(http.StatusInternalServerError).SetMessage("internal error").Get(),
			expectedStatus: http.StatusInternalServerError,
			rawError:       common.ErrInternal,
		},
		{
			name:           "Unexpected",
			inputGuildID:   "4",
			expectedRanks:  ranks.Ranks{},
			expectedError:  *common.NewErrorResponseBuilder(errors.New("tnd")).SetStatus(http.StatusInternalServerError).SetMessage("internal error").Get(),
			expectedStatus: http.StatusInternalServerError,
			rawError:       errors.New("tnd"),
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			mockRankService.On("GetRanks", d.inputGuildID).Return(d.expectedRanks, d.rawError)

			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/rank/%s", d.inputGuildID), nil)
			r.ServeHTTP(rr, req)
			defer rr.Result().Body.Close()

			fmt.Println(rr.Body.String())

			assert.Equal(t, d.expectedStatus, rr.Code)

			if reflect.DeepEqual(d.expectedError, common.ErrorResponse{}) {
				var aResp ranks.Ranks
				common.UnmarshalBody(rr.Result().Body, &aResp)
				assert.Equal(t, d.expectedRanks, aResp)
			} else {
				var aResp common.ErrorResponse
				common.UnmarshalBody(rr.Result().Body, &aResp)
				assert.Equal(t, d.expectedError, aResp)
			}
		})
	}
}
