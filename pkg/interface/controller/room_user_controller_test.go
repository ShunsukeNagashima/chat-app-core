package controller

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllRoomsByUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.RoomUserUsecase)
	validator := validator.New()

	uc := NewRoomUserController(mockUsecase, validator)

	mockRooms := []*model.Room{
		{
			RoomID:   "1",
			Name:     "room-1",
			RoomType: model.Public,
		},
		{
			RoomID:   "2",
			Name:     "room-2",
			RoomType: model.Private,
		},
	}

	testCases := []struct {
		name         string
		userId       string
		mockReturn   []*model.Room
		expectedErr  error
		expectedCode int
	}{
		{
			name:         "Success",
			userId:       "1",
			mockReturn:   mockRooms,
			expectedErr:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid UserID",
			userId:       "invalid_userID",
			mockReturn:   nil,
			expectedErr:  errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsecase.On("GetAllRoomsByUserID", mock.Anything, tc.userId).Return(tc.mockReturn, tc.expectedErr)

			_, ctx, response := prepareRequestAndContext(http.MethodGet, "users/"+tc.userId+"/rooms", gin.Params{{Key: "userId", Value: tc.userId}}, nil)

			uc.GetAllRoomsByUserID(ctx)

			checkResponseRooms(t, tc.expectedErr, tc.expectedCode, response, tc.mockReturn)
		})
	}
}

func TestRemoveUserFromRoom(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.RoomUserUsecase)
	validator := validator.New()

	uc := NewRoomUserController(mockUsecase, validator)

	testCases := []struct {
		name         string
		roomId       string
		userId       string
		expectedErr  error
		expectedCode int
	}{
		{
			name:         "Success",
			roomId:       "1",
			userId:       "1",
			expectedErr:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid RoomID",
			roomId:       "invalid_roomID",
			userId:       "1",
			expectedErr:  errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Invalid UserID",
			roomId:       "1",
			userId:       "invalid_userID",
			expectedErr:  errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsecase.On("RemoveUserFromRoom", mock.Anything, tc.roomId, tc.userId).Return(tc.expectedErr)

			_, ctx, response := prepareRequestAndContext(http.MethodDelete, "rooms/"+tc.roomId+"/users/"+tc.userId, gin.Params{{Key: "roomId", Value: tc.roomId}, {Key: "userId", Value: tc.userId}}, nil)

			uc.RemoveUserFromRoom(ctx)

			checkResponseMessage(t, tc.expectedErr, tc.expectedCode, response, "success to remove the user from the room")
		})
	}
}

func TestAddUsersToRoom(t *testing.T) {
	gin.SetMode(gin.TestMode)
	validator := validator.New()

	testCases := []struct {
		name         string
		roomId       string
		userIDs      []string
		expectedErr  error
		expectedCode int
	}{
		{
			name:         "Success",
			roomId:       "1",
			userIDs:      []string{"1", "2"},
			expectedErr:  nil,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "Invalid UserID",
			roomId:       "1",
			userIDs:      []string{"invalid_userID", "2"},
			expectedErr:  errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsecase := new(mocks.RoomUserUsecase)

			mockUsecase.On("AddUsersToRoom", mock.Anything, tc.roomId, tc.userIDs).Return(tc.expectedErr)

			uc := NewRoomUserController(mockUsecase, validator)

			req := struct {
				UserIDs []string `json:"userIDs"`
			}{
				UserIDs: tc.userIDs,
			}

			reqBody, _ := json.Marshal(req)

			_, ctx, response := prepareRequestAndContext(http.MethodPost, "rooms/"+tc.roomId+"/users", gin.Params{{Key: "roomId", Value: tc.roomId}}, bytes.NewBuffer(reqBody))

			uc.AddUsersToRoom(ctx)

			checkResponseMessage(t, tc.expectedErr, tc.expectedCode, response, "success to add the users to the room")
		})
	}
}

func prepareRequestAndContext(method, url string, params gin.Params, body io.Reader) (*http.Request, *gin.Context, *httptest.ResponseRecorder) {
	request, _ := http.NewRequest(method, url, body)
	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)
	ctx.Params = params
	ctx.Request = request
	return request, ctx, response
}

func checkResponseMessage(t *testing.T, expectedErr error, expectedCode int, response *httptest.ResponseRecorder, expectedMessage string) {
	assert.Equal(t, expectedCode, response.Code)

	var responseBody map[string]string

	if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
		t.Fatal(err)
	}

	if expectedErr == nil {
		assert.Equal(t, expectedMessage, responseBody["result"])
	} else {
		assert.Equal(t, expectedErr.Error(), responseBody["error"])
	}
}

func checkResponseRooms(t *testing.T, expectedErr error, expectedCode int, response *httptest.ResponseRecorder, expectedRooms []*model.Room) {
	if expectedErr == nil {
		assert.Equal(t, expectedCode, response.Code)
		var result struct {
			Result []*model.Room `json:"result"`
		}
		if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expectedRooms, result.Result)
	} else {
		assert.Equal(t, expectedCode, response.Code)
		var result struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expectedErr.Error(), result.Message)
	}
}
