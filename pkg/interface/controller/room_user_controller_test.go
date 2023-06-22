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
		userID       string
		mockReturn   []*model.Room
		expectedErr  error
		expectedCode int
	}{
		{
			name:         "Success",
			userID:       "1",
			mockReturn:   mockRooms,
			expectedErr:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid UserID",
			userID:       "invalid_userID",
			mockReturn:   nil,
			expectedErr:  errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsecase.On("GetAllRoomsByUserID", mock.Anything, tc.userID).Return(tc.mockReturn, tc.expectedErr)

			_, ctx, response := prepareRequestAndContext(http.MethodGet, "users/"+tc.userID+"/rooms", gin.Params{{Key: "userID", Value: tc.userID}}, nil)

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
		roomID       string
		userID       string
		expectedErr  error
		expectedCode int
	}{
		{
			name:         "Success",
			roomID:       "1",
			userID:       "1",
			expectedErr:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid RoomID",
			roomID:       "invalid_roomID",
			userID:       "1",
			expectedErr:  errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Invalid UserID",
			roomID:       "1",
			userID:       "invalid_userID",
			expectedErr:  errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsecase.On("RemoveUserFromRoom", mock.Anything, tc.roomID, tc.userID).Return(tc.expectedErr)

			_, ctx, response := prepareRequestAndContext(http.MethodDelete, "rooms/"+tc.roomID+"/users/"+tc.userID, gin.Params{{Key: "roomID", Value: tc.roomID}, {Key: "userID", Value: tc.userID}}, nil)

			uc.RemoveUserFromRoom(ctx)

			checkResponseMessage(t, tc.expectedErr, tc.expectedCode, response, "success to remove the user from the room")
		})
	}
}

func TestAddUsersToRoom(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.RoomUserUsecase)
	validator := validator.New()

	uc := NewRoomUserController(mockUsecase, validator)

	testCases := []struct {
		name         string
		roomID       string
		userIDs      []string
		expectedErr  error
		expectedCode int
	}{
		{
			name:         "Success",
			roomID:       "1",
			userIDs:      []string{"1", "2"},
			expectedErr:  nil,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "Invalid RoomID",
			roomID:       "invalid_roomID",
			userIDs:      []string{"1", "2"},
			expectedErr:  errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Invalid UserID",
			roomID:       "1",
			userIDs:      []string{"invalid_userID", "2"},
			expectedErr:  errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsecase.On("AddUsersToRoom", mock.Anything, tc.roomID, tc.userIDs).Return(tc.expectedErr)

			req := struct {
				RoomID  string   `json:"roomID"`
				UserIDs []string `json:"userIDs"`
			}{
				RoomID:  tc.roomID,
				UserIDs: tc.userIDs,
			}

			reqBody, _ := json.Marshal(req)

			_, ctx, response := prepareRequestAndContext(http.MethodPost, "rooms/"+tc.roomID+"/users", gin.Params{{Key: "roomID", Value: tc.roomID}}, bytes.NewBuffer(reqBody))

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
	if expectedErr == nil {
		assert.Equal(t, expectedCode, response.Code)
		var result struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expectedMessage, result.Message)
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
