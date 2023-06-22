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
	mockUsecase := new(mocks.RoomUsecase)
	validator := newTestValidator()

	uc := NewRoomController(mockUsecase, validator)

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
			},
			mockReturn:   nil,
			expectedCode: http.StatusCreated,
		},
		{
			name: "Invalid roomType",
			reqBody: map[string]string{
				"name":     "chat_room",
				"roomType": "invalid",
			},
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid name",
			reqBody: map[string]string{
				"name":     "invalid%room",
				"roomType": "public",
			},
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Missing roomType",
			reqBody: map[string]string{
				"name": "chat_room",
			},
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Missing name",
			reqBody: map[string]string{
				"roomType": "public",
			},
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tc.reqBody)

			request, _ := http.NewRequest(http.MethodPost, "/rooms", bytes.NewBuffer(reqBody))
			response := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(response)
			ctx.Request = request

			mockUsecase.On("CreateRoom", mock.Anything, mock.Anything).Return(tc.mockReturn)

			uc.CreateRoom(ctx)

			assert.Equal(t, tc.expectedCode, response.Code)

			if tc.expectedCode == http.StatusCreated {
				mockUsecase.AssertExpectations(t)
				var responseBody map[string]interface{}
				json.Unmarshal(response.Body.Bytes(), &responseBody)

				result, _ := responseBody["result"].(map[string]interface{})
				t.Log(result)
				assert.Equal(t, tc.reqBody["name"], result["name"])
				assert.Equal(t, tc.reqBody["roomType"], result["room_type"])
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
		roomID       string
		mockReturn   *model.Room
		expectedErr  error
		expectedCode int
	}{
		{
			name:         "Success",
			roomID:       "1",
			mockReturn:   &mockRoom,
			expectedErr:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid roomID",
			roomID:       "invalid",
			mockReturn:   nil,
			expectedErr:  errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/rooms/"+tc.roomID, nil)
			response := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(response)
			ctx.Params = gin.Params{{Key: "roomID", Value: tc.roomID}}
			ctx.Request = request

			mockUsecase.On("GetRoomByID", mock.Anything, tc.roomID).Return(tc.mockReturn, tc.expectedErr)

			uc.GetRoomByID(ctx)

			assert.Equal(t, tc.expectedCode, response.Code)

			if tc.expectedCode == http.StatusOK {
				mockUsecase.AssertExpectations(t)
				var responseBody map[string]interface{}
				json.Unmarshal(response.Body.Bytes(), &responseBody)

				var resultRoom model.Room
				roomBytes, _ := json.Marshal(responseBody["result"])
				json.Unmarshal(roomBytes, &resultRoom)

				assert.Equal(t, mockRoom, resultRoom)
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
				var responseBody map[string]interface{}
				json.Unmarshal(response.Body.Bytes(), &responseBody)

				var resultRooms []*model.Room
				roomBytes, _ := json.Marshal(responseBody["result"])
				json.Unmarshal(roomBytes, &resultRooms)

				assert.Equal(t, mockRooms, resultRooms)
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
		roomID       string
		expectedErr  error
		expectedCode int
	}{
		{
			name:         "Success",
			roomID:       "1",
			expectedErr:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid roomID",
			roomID:       "invalid",
			expectedErr:  errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodDelete, "/rooms/"+tc.roomID, nil)
			response := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(response)
			ctx.Params = gin.Params{{Key: "roomID", Value: tc.roomID}}
			ctx.Request = request

			mockUsecase.On("DeleteRoom", mock.Anything, tc.roomID).Return(tc.expectedErr)

			uc.DeleteRoom(ctx)

			assert.Equal(t, tc.expectedCode, response.Code)

			if tc.expectedCode == http.StatusOK {
				mockUsecase.AssertExpectations(t)
				var responseBody map[string]interface{}
				json.Unmarshal(response.Body.Bytes(), &responseBody)

				assert.Equal(t, "room deleted successfully", responseBody["result"])
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
		roomID       string
		reqBody      map[string]string
		expectedErr  error
		expectedCode int
	}{
		{
			name:   "Success",
			roomID: "1",
			reqBody: map[string]string{
				"name":     "chat_room_update",
				"roomType": "public",
			},
			expectedErr:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name:   "Invalid roomType",
			roomID: "1",
			reqBody: map[string]string{
				"name":     "chat_room_update",
				"roomType": "invalid",
			},
			expectedErr:  nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:   "Invalid name",
			roomID: "1",
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

			request, _ := http.NewRequest(http.MethodPut, "/rooms/"+tc.roomID, bytes.NewBuffer(reqBody))
			response := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(response)
			ctx.Params = gin.Params{{Key: "roomID", Value: tc.roomID}}
			ctx.Request = request

			mockUsecase.On("UpdateRoom", mock.Anything, mock.Anything).Return(tc.expectedErr)

			uc.UpdateRoom(ctx)

			assert.Equal(t, tc.expectedCode, response.Code)

			if tc.expectedCode == http.StatusOK {
				mockUsecase.AssertExpectations(t)
				var responseBody map[string]interface{}
				json.Unmarshal(response.Body.Bytes(), &responseBody)

				assert.Equal(t, "room updated successfully", responseBody["result"])
			}
		})
	}

}

func newTestValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("alnumdash", isAlnumOrDash)
	return v
}

func isAlnumOrDash(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(fl.Field().String())
}
