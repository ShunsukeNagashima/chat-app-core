package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shunsukenagashima/chat-api/pkg/clock"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetMessagesByRoomID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.MessageUsecase)
	validator := newTestValidator()
	clock := clock.FixedClocker{}

	roomId := "1"
	var mockMessages []*model.Message
	for i := 0; i < 10; i++ {
		mockMessages = append(mockMessages, &model.Message{
			MessageID: strconv.Itoa(i),
			RoomID:    roomId,
			UserID:    "1",
			Content:   "Hello " + strconv.Itoa(i),
			CreatedAt: clock.Now().Add(time.Duration(i) * time.Minute),
		})
	}
	request, _ := http.NewRequest(http.MethodGet, "messages"+roomId, nil)
	response := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(response)
	ctx.Params = gin.Params{{Key: "roomId", Value: roomId}}
	ctx.Request = request

	mockUsecase.On("GetMessagesByRoomID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockMessages, "10", nil)

	mc := NewMessageController(mockUsecase, validator)

	mc.GetMessagesByRoomID(ctx)

	assert.Equal(t, http.StatusOK, response.Code)
	mockUsecase.AssertExpectations(t)

	var result struct {
		Result  []*model.Message `json:"result"`
		NextKey string           `json:"nextKey"`
	}

	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, mockMessages, result.Result)
}

func TestCreateMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	validator := newTestValidator()

	teatCases := []struct {
		name         string
		reqBody      map[string]string
		mockReturn   error
		expectedCode int
	}{

		{
			name: "Success",
			reqBody: map[string]string{
				"userId":  "1",
				"content": "Hello",
			},
			mockReturn:   nil,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "Invalid Body",
			reqBody:      map[string]string{},
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Missing userId",
			reqBody: map[string]string{
				"content": "Hello",
			},
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Missing content",
			reqBody: map[string]string{
				"userId": "1",
			},
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Empty Content",
			reqBody: map[string]string{
				"userId":  "1",
				"content": "",
			},
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Create Failed",
			reqBody: map[string]string{
				"userId":  "1",
				"content": "Hello",
			},
			mockReturn:   errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range teatCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsecase := new(mocks.MessageUsecase)
			mockUsecase.On("CreateMessage", mock.Anything, mock.Anything).Return(tc.mockReturn)

			reqBody, err := json.Marshal(tc.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			request, _ := http.NewRequest(http.MethodPost, "messages", bytes.NewBuffer(reqBody))
			response := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(response)
			ctx.Params = gin.Params{{Key: "roomId", Value: "1"}}
			ctx.Request = request

			mc := NewMessageController(mockUsecase, validator)

			mc.CreateMessage(ctx)

			assert.Equal(t, tc.expectedCode, response.Code)
		})
	}
}

func TestUpdateMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	validator := newTestValidator()

	teatCases := []struct {
		name         string
		reqBody      map[string]string
		mockReturn   error
		expectedCode int
	}{

		{
			name: "Success",
			reqBody: map[string]string{
				"content": "Hello",
			},
			mockReturn:   nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid Body",
			reqBody:      map[string]string{},
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Missing content",
			reqBody: map[string]string{
				"userId": "1",
			},
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Empty Content",
			reqBody: map[string]string{
				"content": "",
			},
			mockReturn:   nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Update Failed",
			reqBody: map[string]string{
				"content": "Hello",
			},
			mockReturn:   errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range teatCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsecase := new(mocks.MessageUsecase)
			mockUsecase.On("UpdateMessage", mock.Anything, mock.Anything, mock.Anything).Return(tc.mockReturn)

			reqBody, err := json.Marshal(tc.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			request, _ := http.NewRequest(http.MethodPut, "messages", bytes.NewBuffer(reqBody))
			response := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(response)
			ctx.Params = gin.Params{{Key: "messageId", Value: "1"}}
			ctx.Request = request

			mc := NewMessageController(mockUsecase, validator)

			mc.UpdateMessage(ctx)

			assert.Equal(t, tc.expectedCode, response.Code)
		})
	}
}

func TestDeleteMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	validator := newTestValidator()

	teatCases := []struct {
		name         string
		mockReturn   error
		expectedCode int
	}{

		{
			name:         "Success",
			mockReturn:   nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Delete Failed",
			mockReturn:   errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range teatCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsecase := new(mocks.MessageUsecase)
			mockUsecase.On("DeleteMessage", mock.Anything, mock.Anything).Return(tc.mockReturn)

			request, _ := http.NewRequest(http.MethodDelete, "messages", nil)
			response := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(response)
			ctx.Params = gin.Params{{Key: "messageId", Value: "1"}}
			ctx.Request = request

			mc := NewMessageController(mockUsecase, validator)

			mc.DeleteMessage(ctx)

			assert.Equal(t, tc.expectedCode, response.Code)
		})
	}
}
