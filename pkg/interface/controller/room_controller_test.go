package controller

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateRoom(t *testing.T) {
	gin.SetMode(gin.TestMode)
	validator := newTestValidator()

	testCases := []struct {
		name         string
		reqBody      map[string]string
		mockReturn   error
		expectedCode int
	}{
		{
			name: "Success",
			reqBody: map[string]string{
				"name":     "chat_room",
				"roomType": "public",
				"ownerId":  "1",
			},
			mockReturn:   nil,
			expectedCode: http.StatusCreated,
		},
		{
			name: "Invalid roomType",
			reqBody: map[string]string{
				"name":     "chat_room",
				"roomType": "invalid",
				"ownerId":  "1",
			},
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid name",
			reqBody: map[string]string{
				"name":     "invalid%room",
				"roomType": "public",
				"ownerId":  "1",
			},
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Missing roomType",
			reqBody: map[string]string{
				"name":    "chat_room",
				"ownerId": "1",
			},
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Missing name",
			reqBody: map[string]string{
				"roomType": "public",
				"ownerId":  "1",
			},
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Missing ownerId",
			reqBody: map[string]string{
				"name":     "chat_room",
				"roomType": "public",
			},
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsecase := new(mocks.RoomUsecase)

			reqBody, _ := json.Marshal(tc.reqBody)

			request, _ := http.NewRequest(http.MethodPost, "/rooms", bytes.NewBuffer(reqBody))
			response := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(response)
			ctx.Request = request

			mockUsecase.On("CreateRoom", mock.Anything, mock.Anything, mock.Anything).Return(tc.mockReturn)

			uc := NewRoomController(mockUsecase, validator)

			uc.CreateRoom(ctx)

			assert.Equal(t, tc.expectedCode, response.Code)

			if tc.expectedCode == http.StatusCreated {
				mockUsecase.AssertExpectations(t)
				var responseBody map[string]interface{}
				if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
					t.Fatal(err)
				}

				result, _ := responseBody["result"].(map[string]interface{})
				assert.Equal(t, tc.reqBody["name"], result["name"])
				assert.Equal(t, tc.reqBody["roomType"], result["roomType"])
			}
		})
	}
}

func TestGetRoomByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.RoomUsecase)
	validator := newTestValidator()

	uc := NewRoomController(mockUsecase, validator)

	mockRoom := model.Room{
		RoomID:   "1",
		Name:     "chat_room",
		RoomType: model.Public,
	}

	testCases := []struct {
		name         string
		roomId       string
		mockReturn   *model.Room
		expectedErr  error
		expectedCode int
	}{
		{
			name:         "Success",
			roomId:       "1",
			mockReturn:   &mockRoom,
			expectedErr:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid roomId",
			roomId:       "invalid",
			mockReturn:   nil,
			expectedErr:  errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/rooms/"+tc.roomId, nil)
			response := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(response)
			ctx.Params = gin.Params{{Key: "roomId", Value: tc.roomId}}
			ctx.Request = request

			mockUsecase.On("GetRoomByID", mock.Anything, tc.roomId).Return(tc.mockReturn, tc.expectedErr)

			uc.GetRoomByID(ctx)

			assert.Equal(t, tc.expectedCode, response.Code)

			if tc.expectedCode == http.StatusOK {
				mockUsecase.AssertExpectations(t)

				var result struct {
					Result model.Room `json:"result"`
				}

				if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, mockRoom, result.Result)
			}
		})
	}
}

func TestGetAllPublicRooms(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.RoomUsecase)
	validator := newTestValidator()

	uc := NewRoomController(mockUsecase, validator)

	mockRooms := []*model.Room{
		{
			RoomID:   "1",
			Name:     "chat_room_1",
			RoomType: model.Public,
		},
		{
			RoomID:   "2",
			Name:     "chat_room_2",
			RoomType: model.Private,
		},
	}

	testCases := []struct {
		name         string
		mockReturn   []*model.Room
		expectedErr  error
		expectedCode int
	}{
		{
			name:         "Success",
			mockReturn:   mockRooms,
			expectedErr:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Empty slice returns",
			mockReturn:   []*model.Room{},
			expectedErr:  nil,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/rooms", nil)
			response := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(response)
			ctx.Request = request

			mockUsecase.On("GetAllPublicRoom", mock.Anything).Return(tc.mockReturn, tc.expectedErr)

			uc.GetAllPublicRoom(ctx)

			assert.Equal(t, tc.expectedCode, response.Code)

			if tc.expectedCode == http.StatusOK {
				mockUsecase.AssertExpectations(t)

				var result struct {
					Result []*model.Room `json:"result"`
				}
				if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, mockRooms, result.Result)
			}
		})
	}
}

func TestDeleteRoom(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.RoomUsecase)
	validator := newTestValidator()

	uc := NewRoomController(mockUsecase, validator)

	testCases := []struct {
		name         string
		roomId       string
		expectedErr  error
		expectedCode int
	}{
		{
			name:         "Success",
			roomId:       "1",
			expectedErr:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid roomId",
			roomId:       "invalid",
			expectedErr:  errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodDelete, "/rooms/"+tc.roomId, nil)
			response := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(response)
			ctx.Params = gin.Params{{Key: "roomId", Value: tc.roomId}}
			ctx.Request = request

			mockUsecase.On("DeleteRoom", mock.Anything, tc.roomId).Return(tc.expectedErr)

			uc.DeleteRoom(ctx)

			assert.Equal(t, tc.expectedCode, response.Code)

			if tc.expectedCode == http.StatusOK {
				mockUsecase.AssertExpectations(t)

				var result struct {
					Result string `json:"result"`
				}

				if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, "room deleted successfully", result.Result)
			}
		})
	}
}

func TestUpdateRoom(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.RoomUsecase)
	validator := newTestValidator()

	uc := NewRoomController(mockUsecase, validator)

	testCases := []struct {
		name         string
		roomId       string
		reqBody      map[string]string
		expectedErr  error
		expectedCode int
	}{
		{
			name:   "Success",
			roomId: "1",
			reqBody: map[string]string{
				"name":     "chat_room_update",
				"roomType": "public",
			},
			expectedErr:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name:   "Invalid roomType",
			roomId: "1",
			reqBody: map[string]string{
				"name":     "chat_room_update",
				"roomType": "invalid",
			},
			expectedErr:  nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "Invalid name",
			roomId: "1",
			reqBody: map[string]string{
				"name":     "invalid%room",
				"roomType": "public",
			},
			expectedErr:  nil,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tc.reqBody)

			request, _ := http.NewRequest(http.MethodPut, "/rooms/"+tc.roomId, bytes.NewBuffer(reqBody))
			response := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(response)
			ctx.Params = gin.Params{{Key: "roomId", Value: tc.roomId}}
			ctx.Request = request

			mockUsecase.On("UpdateRoom", mock.Anything, mock.Anything).Return(tc.expectedErr)

			uc.UpdateRoom(ctx)

			assert.Equal(t, tc.expectedCode, response.Code)

			if tc.expectedCode == http.StatusOK {
				mockUsecase.AssertExpectations(t)

				var result struct {
					Result string `json:"result"`
				}
				if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, "room updated successfully", result.Result)
			}
		})
	}
}

func newTestValidator() *validator.Validate {
	v := validator.New()
	if err := v.RegisterValidation("alnumdash", isAlnumOrDash); err != nil {
		panic(err)
	}
	return v
}

func isAlnumOrDash(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(fl.Field().String())
}
