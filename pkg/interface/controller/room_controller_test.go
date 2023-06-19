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

func TestCreateRoom_WhenRequestBindingFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.RoomUsecase)
	v := newTestValidator()

	rc := NewRoomController(mockUsecase, v)

	reqBody, _ := json.Marshal(map[string]string{
		"invalid_key": "invalid_value",
	})

	request, _ := http.NewRequest(http.MethodPost, "/rooms", bytes.NewBuffer(reqBody))
	response := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(response)
	ctx.Request = request

	rc.CreateRoom(ctx)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestCreateRoom_WhenValidationFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.RoomUsecase)
	v := newTestValidator()

	rc := NewRoomController(mockUsecase, v)

	mockUsecase.On("CreateRoom", mock.Anything, mock.Anything).Return(errors.New("some error"))

	testCases := []struct {
		name    string
		reqBody map[string]string
	}{
		{
			name: "Missing roomType",
			reqBody: map[string]string{
				"name": "chat_room",
			},
		},
		{
			name: "Missing name",
			reqBody: map[string]string{
				"roomType": "public",
			},
		},
		{
			name: "Invalid name",
			reqBody: map[string]string{
				"name":     "invalid%room",
				"roomType": "public",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tc.reqBody)

			request, _ := http.NewRequest(http.MethodPost, "/rooms", bytes.NewBuffer(reqBody))
			response := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(response)
			ctx.Request = request

			rc.CreateRoom(ctx)

			assert.Equal(t, http.StatusBadRequest, response.Code)
		})
	}
}

func TestCreateRoom_WhenParseRoomTypeFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.RoomUsecase)
	v := newTestValidator()

	rc := NewRoomController(mockUsecase, v)

	reqBody, _ := json.Marshal(map[string]string{
		"name":     "chat_room",
		"roomType": "invalid",
	})

	request, _ := http.NewRequest(http.MethodPost, "/rooms", bytes.NewBuffer(reqBody))
	response := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(response)
	ctx.Request = request

	mockUsecase.On("CreateRoom", mock.Anything, mock.Anything).Return(nil)

	rc.CreateRoom(ctx)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestCreateRoom_WhenExecutionSucceeds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.RoomUsecase)
	v := newTestValidator()

	rc := NewRoomController(mockUsecase, v)

	reqBody, _ := json.Marshal(map[string]string{
		"name":     "chat_room",
		"roomType": "public",
	})

	request, _ := http.NewRequest(http.MethodPost, "/rooms", bytes.NewBuffer(reqBody))
	response := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(response)
	ctx.Request = request

	mockUsecase.On("CreateRoom", mock.Anything, mock.Anything).Return(nil)

	rc.CreateRoom(ctx)

	mockUsecase.AssertExpectations(t)

	assert.Equal(t, http.StatusOK, response.Code)
}

func TestGetRoomByID_WhenRoomExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.RoomUsecase)
	v := newTestValidator()

	rc := NewRoomController(mockUsecase, v)

	roomID := "1"
	mockRoom := model.Room{
		RoomID:   roomID,
		Name:     "chat_room",
		RoomType: model.Public,
	}

	mockUsecase.On("GetRoomByID", mock.Anything, roomID).Return(&mockRoom, nil)

	request, _ := http.NewRequest(http.MethodGet, "/rooms/"+roomID, nil)
	response := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(response)
	ctx.Params = gin.Params{{Key: "roomID", Value: roomID}}
	ctx.Request = request

	rc.GetRoomByID(ctx)

	assert.Equal(t, http.StatusOK, response.Code)

	var responseBody map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	var resultRoom model.Room
	roomBytes, _ := json.Marshal(responseBody["result"])
	json.Unmarshal(roomBytes, &resultRoom)

	assert.Equal(t, mockRoom, resultRoom)
	mockUsecase.AssertExpectations(t)
}

func TestGetRoomByID_InvalidRoomID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.RoomUsecase)
	v := newTestValidator()

	rc := NewRoomController(mockUsecase, v)

	roomID := "1"

	mockUsecase.On("GetRoomByID", mock.Anything, roomID).Return(nil, errors.New("some error"))

	request, _ := http.NewRequest(http.MethodGet, "/rooms/"+roomID, nil)
	response := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(response)
	ctx.Params = gin.Params{{Key: "roomID", Value: roomID}}
	ctx.Request = request

	rc.GetRoomByID(ctx)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}

func TestGetAllPublicRooms_WhenRoomExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.RoomUsecase)
	v := newTestValidator()

	rc := NewRoomController(mockUsecase, v)

	mockRooms := []*model.Room{
		{
			RoomID:   "1",
			Name:     "chat_room_1",
			RoomType: model.Public,
		},
		{
			RoomID:   "2",
			Name:     "chat_room_2",
			RoomType: model.Public,
		},
		{
			RoomID:   "3",
			Name:     "chat_room_3",
			RoomType: model.Public,
		},
	}

	mockUsecase.On("GetAllPublicRoom", mock.Anything).Return(mockRooms, nil)

	request, _ := http.NewRequest(http.MethodGet, "/rooms", nil)
	response := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(response)
	ctx.Request = request

	rc.GetAllPublicRoom(ctx)

	assert.Equal(t, http.StatusOK, response.Code)

	var responseBody map[string][]interface{}
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	var resultRooms []*model.Room
	roomBytes, _ := json.Marshal(responseBody["result"])

	json.Unmarshal(roomBytes, &resultRooms)

	assert.Equal(t, mockRooms, resultRooms)
	mockUsecase.AssertExpectations(t)
}

func TestGetAllPublicRooms_WhenEmptySliceReturns(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.RoomUsecase)
	v := newTestValidator()

	rc := NewRoomController(mockUsecase, v)

	mockRooms := []*model.Room{}

	mockUsecase.On("GetAllPublicRoom", mock.Anything).Return(mockRooms, nil)

	request, _ := http.NewRequest(http.MethodGet, "/rooms", nil)
	response := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(response)
	ctx.Request = request

	rc.GetAllPublicRoom(ctx)

	assert.Equal(t, http.StatusOK, response.Code)

	var responseBody map[string][]interface{}
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	var resultRooms []*model.Room
	roomBytes, _ := json.Marshal(responseBody["result"])

	json.Unmarshal(roomBytes, &resultRooms)

	assert.Equal(t, mockRooms, resultRooms)
	mockUsecase.AssertExpectations(t)
}

func TestDeleteRoom_WhenExecutionSucceeds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.RoomUsecase)
	v := newTestValidator()

	rc := NewRoomController(mockUsecase, v)

	roomID := "1"
	resultMsg := "room deleted successfully"

	mockUsecase.On("DeleteRoom", mock.Anything, roomID).Return(nil)

	request, _ := http.NewRequest(http.MethodDelete, "/rooms/"+roomID, nil)
	response := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(response)
	ctx.Params = gin.Params{{Key: "roomID", Value: roomID}}
	ctx.Request = request

	rc.DeleteRoom(ctx)

	assert.Equal(t, http.StatusOK, response.Code)

	var responseBody map[string]string
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	assert.Equal(t, resultMsg, responseBody["result"])
	mockUsecase.AssertExpectations(t)
}

func TestDeleteRoom_WhenRoomDoesntExist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.RoomUsecase)
	v := newTestValidator()

	rc := NewRoomController(mockUsecase, v)

	roomID := "1"

	mockUsecase.On("DeleteRoom", mock.Anything, roomID).Return(errors.New("room not found"))

	request, _ := http.NewRequest(http.MethodDelete, "/rooms/"+roomID, nil)
	response := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(response)
	ctx.Params = gin.Params{{Key: "roomID", Value: roomID}}
	ctx.Request = request

	rc.DeleteRoom(ctx)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}

func TestUpdateRoom_WhenExecutionSucceeds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.RoomUsecase)
	v := newTestValidator()

	rc := NewRoomController(mockUsecase, v)

	roomID := "1"
	resultMsg := "room updated successfully"

	reqBody, _ := json.Marshal(map[string]string{
		"name":     "chat_room_update",
		"roomType": "public",
	})

	request, _ := http.NewRequest(http.MethodPut, "/rooms/"+roomID, bytes.NewBuffer(reqBody))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(response)
	ctx.Params = gin.Params{{Key: "roomID", Value: roomID}}
	ctx.Request = request

	mockUsecase.On("UpdateRoom", mock.Anything, mock.Anything).Return(nil)

	rc.UpdateRoom(ctx)

	assert.Equal(t, http.StatusOK, response.Code)

	var responseBody map[string]string
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	assert.Equal(t, resultMsg, responseBody["result"])
	mockUsecase.AssertExpectations(t)
}

func newTestValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("alnumdash", isAlnumOrDash)
	return v
}

func isAlnumOrDash(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(fl.Field().String())
}
