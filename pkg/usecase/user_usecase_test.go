package usecase

import (
	"context"
	"testing"

	"firebase.google.com/go/auth"
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

	mockRepo.On("Create", mock.Anything, mockUser).Return(nil)
	mockAuth.On("GetFirebaseUser", mock.Anything, mockUser.UserID).Return(&auth.Token{UID: "1"}, nil)

	userUsecase := NewUserUsecase(mockRepo, mockAuth)
	err := userUsecase.CreateUser(context.Background(), mockUser)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

func TestCreateUser_WhenIDIsntMatch(t *testing.T) {
	mockRepo := new(repoMocks.UserRepository)
	mockAuth := new(authMocks.FirebaseAuthenticator)
	mockUser := &model.User{
		UserID:   "1",
		Username: "user-1",
		Email:    "user-1@example.com",
	}

	mockAuth.On("GetFirebaseUser", mock.Anything, mockUser.UserID).Return(&auth.Token{UID: "2"}, nil)

	userUsecase := NewUserUsecase(mockRepo, mockAuth)
	err := userUsecase.CreateUser(context.Background(), mockUser)

	assert.Error(t, err)
	assert.EqualError(t, err, "provided user ID does not match the user ID in Firebase token")
	mockAuth.AssertExpectations(t)
}
