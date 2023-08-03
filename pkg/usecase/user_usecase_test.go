package usecase

import (
	"context"
	"strconv"
	"testing"

	"firebase.google.com/go/auth"
	"github.com/shunsukenagashima/chat-api/pkg/apperror"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	repoMocks "github.com/shunsukenagashima/chat-api/pkg/domain/repository/mocks"
	authMocks "github.com/shunsukenagashima/chat-api/pkg/infra/auth/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserByID(t *testing.T) {
	mockRepo := new(repoMocks.UserRepository)
	mockAuth := new(authMocks.FirebaseAuthenticator)
	mockUser := &model.User{
		UserID:   "1",
		Username: "user-1",
		Email:    "user-1@example.com",
	}

	mockRepo.On("GetByID", mock.Anything, mockUser.UserID).Return(mockUser, nil)
	userUsecase := NewUserUsecase(mockRepo, mockAuth)

	user, err := userUsecase.GetUserByID(context.Background(), mockUser.UserID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, mockUser.UserID, user.UserID)
	assert.Equal(t, mockUser.Username, user.Username)
	assert.Equal(t, mockUser.Email, user.Email)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(repoMocks.UserRepository)
	mockAuth := new(authMocks.FirebaseAuthenticator)
	mockUser := &model.User{
		UserID:   "1",
		Username: "user-1",
		Email:    "user-1@example.com",
	}

	idToken := "test_id_token"

	mockRepo.On("Create", mock.Anything, mockUser).Return(nil)
	mockRepo.On("GetByID", mock.Anything, mockUser.UserID).Return(nil, &apperror.NotFoundErr{})
	mockAuth.On("GetFirebaseUser", mock.Anything, idToken).Return(&auth.Token{UID: "1"}, nil)

	userUsecase := NewUserUsecase(mockRepo, mockAuth)
	err := userUsecase.CreateUser(context.Background(), mockUser, idToken)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	// uncomment this line when NAT Gateway is ready
	// mockAuth.AssertExpectations(t)
}

// uncomment this test when NAT Gateway is ready
// func TestCreateUser_WhenIDIsntMatch(t *testing.T) {
// 	mockRepo := new(repoMocks.UserRepository)
// 	mockAuth := new(authMocks.FirebaseAuthenticator)
// 	mockUser := &model.User{
// 		UserID:   "1",
// 		Username: "user-1",
// 		Email:    "user-1@example.com",
// 	}

// 	idToken := "test_id_token"

// 	mockAuth.On("GetFirebaseUser", mock.Anything, idToken).Return(&auth.Token{UID: "2"}, nil)

// 	userUsecase := NewUserUsecase(mockRepo, mockAuth)
// 	err := userUsecase.CreateUser(context.Background(), mockUser, idToken)

// 	assert.Error(t, err)
// 	assert.EqualError(t, err, "provided user ID does not match the user ID in Firebase token")
// 	mockAuth.AssertExpectations(t)
// }

func TestBatchGetUsers(t *testing.T) {
	mockRepo := new(repoMocks.UserRepository)
	mockAuth := new(authMocks.FirebaseAuthenticator)

	var mockUsers []*model.User

	for i := 1; i <= 3; i++ {
		iStr := strconv.Itoa(i)
		mockUser := &model.User{
			UserID:   "id-" + iStr,
			Username: "user-" + iStr,
			Email:    "user-" + iStr + "@example.com",
		}
		mockUsers = append(mockUsers, mockUser)
	}

	mockRepo.On("BatchGetUsers", mock.Anything, []string{"1", "2", "3"}).Return(mockUsers, nil)

	userUsecase := NewUserUsecase(mockRepo, mockAuth)

	users, err := userUsecase.BatchGetUsers(context.Background(), []string{"1", "2", "3"})

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, len(mockUsers), len(users))
	mockRepo.AssertExpectations(t)
}
