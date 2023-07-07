package controller

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.UserUsecase)
	validator := validator.New()

	uc := NewUserController(mockUsecase, validator)

	mockUser := &model.User{
		UserID:   "1",
		Username: "user-1",
		Email:    "user-1@example.com",
	}

	testCases := []struct {
		name        string
		userId      string
		mockReturn  *model.User
		mockError   error
		isErrorCase bool
	}{
		{
			name:        "Success",
			userId:      "1",
			mockReturn:  mockUser,
			mockError:   nil,
			isErrorCase: false,
		},
		{
			name:        "Invalid UserID",
			userId:      "invalid_userID",
			mockReturn:  nil,
			mockError:   errors.New("some error"),
			isErrorCase: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsecase.On("GetUserByID", mock.Anything, tc.userId).Return(tc.mockReturn, tc.mockError)

			request, _ := http.NewRequest(http.MethodGet, "/users/"+mockUser.UserID, nil)
			response := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(response)
			ctx.Params = gin.Params{{Key: "userId", Value: tc.userId}}
			ctx.Request = request

			uc.GetUserByID(ctx)

			if tc.isErrorCase {
				assert.Equal(t, http.StatusInternalServerError, response.Code)
			} else {
				assert.Equal(t, http.StatusOK, response.Code)

				var responseBody map[string]interface{}
				if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
					t.Fatal(err)
				}

				var resultUser model.User
				userBytes, _ := json.Marshal(responseBody["result"])
				if err := json.Unmarshal(userBytes, &resultUser); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, *mockUser, resultUser)
				mockUsecase.AssertExpectations(t)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	validator := validator.New()

	testCases := []struct {
		name         string
		reqBody      map[string]string
		expectedCode int
	}{
		{
			name: "Success",
			reqBody: map[string]string{
				"userId":  "1",
				"name":    "user-1",
				"email":   "user-1@example.com",
				"idToken": "test_id_token",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Missing name",
			reqBody: map[string]string{
				"userId":  "2",
				"email":   "user-2@example.com",
				"idToken": "test_id_token",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Missing email",
			reqBody: map[string]string{
				"userId":  "3",
				"name":    "user-3",
				"idToken": "test_id_token",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid email",
			reqBody: map[string]string{
				"userId":  "4",
				"name":    "user-4",
				"email":   "invalid_email",
				"idToken": "test_id_token",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUsecase := new(mocks.UserUsecase)

			if tc.expectedCode == http.StatusOK {
				mockUsecase.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			}

			uc := NewUserController(mockUsecase, validator)

			reqBody, _ := json.Marshal(tc.reqBody)
			request, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
			response := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(response)
			ctx.Request = request

			uc.CreateUser(ctx)

			assert.Equal(t, tc.expectedCode, response.Code)

			if tc.expectedCode == http.StatusOK {
				mockUsecase.AssertExpectations(t)
				var responseBody map[string]interface{}
				if err := json.Unmarshal(response.Body.Bytes(), &responseBody); err != nil {
					t.Fatal(err)
				}

				result, _ := responseBody["result"].(map[string]interface{})

				t.Log(result)
				assert.Equal(t, tc.reqBody["name"], result["userName"])
				assert.Equal(t, tc.reqBody["email"], result["email"])
			}
		})
	}
}

func TestSearchUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(mocks.UserUsecase)
	validator := validator.New()

	var mockUsers []*model.User
	for i := 1; i <= 3; i++ {
		mockUsers = append(mockUsers, &model.User{
			UserID:   strconv.Itoa(i),
			Username: "user-" + strconv.Itoa(i),
			Email:    "user-" + strconv.Itoa(i) + "@example.com",
		})
	}

	request, _ := http.NewRequest(http.MethodPost, "/users/search?query=user&from=0&size=10", nil)
	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)
	ctx.Request = request

	mockUsecase.On("SearchUsers", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockUsers, nil)

	uc := NewUserController(mockUsecase, validator)

	uc.SearchUsers(ctx)

	assert.Equal(t, http.StatusOK, response.Code)

	var result struct {
		Result []*model.User `json:"result"`
	}

	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, mockUsers, result.Result)
	mockUsecase.AssertExpectations(t)
}
