package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mogs "github.com/finkabaj/hyde-bot/internals/backend/mocks"
	"github.com/finkabaj/hyde-bot/internals/ranks"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreateRanks(t *testing.T) {
	posi := ranks.Ranks{
		Ranks: []ranks.Rank{
			{
				XP:    100,
				ID:    "c4658312-6000-403f-8eee-ab1628e05368",
				Role:  nil,
				Level: 1,
			},
			{
				XP:    200,
				ID:    "3eb2d525-0973-4e0e-8b83-caf1edd17cfe",
				Role:  nil,
				Level: 2,
			},
		},
		GuildID: "1",
		OwnerID: "12",
	}

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
			input:          posi,
			expected:       posi,
			expectedError:  common.ErrorResponse{},
			expectedStatus: http.StatusCreated,
			rawError:       nil,
		},
		{
			name:           "BadRequest",
			input:          posi,
			expected:       ranks.Ranks{},
			expectedError:  *common.NewErrorResponseBuilder(common.ErrBadRequest).SetStatus(http.StatusBadRequest).SetMessage("invalid request").Get(),
			expectedStatus: http.StatusBadRequest,
			rawError:       common.ErrBadRequest,
		},
		{
			name:           "Internal",
			input:          posi,
			expected:       ranks.Ranks{},
			expectedError:  *common.NewErrorResponseBuilder(common.ErrInternal).SetStatus(http.StatusInternalServerError).SetMessage("internal error").Get(),
			expectedStatus: http.StatusInternalServerError,
			rawError:       common.ErrInternal,
		},
		{
			name:           "NotFound",
			input:          posi,
			expected:       ranks.Ranks{},
			expectedError:  *common.NewErrorResponseBuilder(common.ErrNotFound).SetStatus(http.StatusNotFound).SetMessage("guild id or owner id not found").Get(),
			expectedStatus: http.StatusNotFound,
			rawError:       common.ErrNotFound,
		},
		{
			name:           "Unexpected",
			input:          posi,
			expected:       ranks.Ranks{},
			expectedError:  *common.NewErrorResponseBuilder(errors.New("tnd")).SetStatus(http.StatusInternalServerError).SetMessage("internal error").Get(),
			expectedStatus: http.StatusInternalServerError,
			rawError:       errors.New("tnd"),
		},
	}

	mockRankService := mogs.MockRankService{}
	r := chi.NewRouter()
	rankc := NewRankController(&mockRankService, mogs.NewMockLogger())
	rankc.RegisterRoutes(r)
	for _, d := range data {
		mockRankService = mogs.MockRankService{}
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
		})
	}
}
