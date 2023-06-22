package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllRoomsByUserID(t *testing.T) {
	mockRoomUserRepo := new(mocks.RoomUserRepository)
	mockUserRepo := new(mocks.UserRepository)
	mockRooms := []*model.Room{
		{
			RoomID:   "1",
			Name:     "room-1",
			RoomType: model.Public,
		},
		{
			RoomID:   "2",
			Name:     "room-2",
			RoomType: model.Public,
		},
	}

	mockRoomUserRepo.On("GetAllRoomsByUserID", mock.Anything, mock.Anything).Return(mockRooms, nil)

	roomUserUsecase := NewRoomUserUsecase(mockRoomUserRepo, mockUserRepo)

	rooms, err := roomUserUsecase.GetAllRoomsByUserID(context.Background(), "1")

	assert.NoError(t, err)
	assert.NotEmpty(t, rooms)
	assert.Equal(t, len(mockRooms), len(rooms))
	for i, room := range rooms {
		assert.Equal(t, mockRooms[i].RoomID, room.RoomID)
		assert.Equal(t, mockRooms[i].Name, room.Name)
		assert.Equal(t, mockRooms[i].RoomType, room.RoomType)
	}
	mockRoomUserRepo.AssertExpectations(t)
}

func TestAddUserToRoom(t *testing.T) {
	mockRoomUserRepo := new(mocks.RoomUserRepository)
	mockUserRepo := new(mocks.UserRepository)
	mockUser := &model.User{
		UserID:   "1",
		Username: "user-1",
		Email:    "user-1@example.com",
	}

	testCases := []struct {
		name        string
		user        *model.User
		userID      string
		roomID      string
		expectedErr error
	}{
		{
			name:        "Success",
			user:        mockUser,
			userID:      mockUser.UserID,
			roomID:      "1",
			expectedErr: nil,
		},
		{
			name:        "UserNotFound",
			user:        nil,
			userID:      "invalid_user_id",
			roomID:      "1",
			expectedErr: errors.New("user not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo.On("GetByID", mock.Anything, tc.userID).Return(tc.user, tc.expectedErr)
			if tc.user != nil {
				mockRoomUserRepo.On("AddUserToRoom", mock.Anything, tc.roomID, tc.userID).Return(nil)
			}

			roomUserUsecase := NewRoomUserUsecase(mockRoomUserRepo, mockUserRepo)

			err := roomUserUsecase.AddUserToRoom(context.Background(), tc.roomID, tc.userID)

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				mockUserRepo.AssertExpectations(t)
				mockRoomUserRepo.AssertExpectations(t)
			}
		})
	}
}

func TestRemoveUserFromRoom(t *testing.T) {
	mockRoomUserRepo := new(mocks.RoomUserRepository)
	mockUserRepo := new(mocks.UserRepository)

	roomID := "1"
	userID := "1"

	mockRoomUserRepo.On("RemoveUserFromRoom", mock.Anything, roomID, userID).Return(nil)

	roomUserUsecase := NewRoomUserUsecase(mockRoomUserRepo, mockUserRepo)

	err := roomUserUsecase.RemoveUserFromRoom(context.Background(), roomID, userID)

	assert.NoError(t, err)
	mockRoomUserRepo.AssertExpectations(t)
}

func TestAddUsersToRoom(t *testing.T) {
	mockRoomUserRepo := new(mocks.RoomUserRepository)
	mockUserRepo := new(mocks.UserRepository)
	mockUsers := []*model.User{
		{
			UserID:   "1",
			Username: "user-1",
			Email:    "user-1@example.com",
		},
		{
			UserID:   "2",
			Username: "user-2",
			Email:    "user-2@example.com",
		},
		{
			UserID:   "3",
			Username: "user-3",
			Email:    "user-3@example.com",
		},
	}

	testCases := []struct {
		name        string
		userIDs     []string
		roomID      string
		expectedErr error
	}{
		{
			name:        "Success",
			userIDs:     []string{mockUsers[0].UserID, mockUsers[1].UserID, mockUsers[2].UserID},
			roomID:      "1",
			expectedErr: nil,
		},
		{
			name:        "OneofUsersNotFound",
			userIDs:     []string{mockUsers[0].UserID, mockUsers[1].UserID, "invalid_user_id"},
			roomID:      "1",
			expectedErr: errors.New("user not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			for i, userID := range tc.userIDs {
				if userID == "invalid_user_id" {
					mockUserRepo.On("GetByID", mock.Anything, userID).Return(nil, errors.New("user not found"))
				} else {
					mockUserRepo.On("GetByID", mock.Anything, userID).Return(mockUsers[i], nil)
				}
			}
			mockRoomUserRepo.On("AddUsersToRoom", mock.Anything, tc.roomID, tc.userIDs).Return(nil)

			roomUserUsecase := NewRoomUserUsecase(mockRoomUserRepo, mockUserRepo)

			err := roomUserUsecase.AddUsersToRoom(context.Background(), tc.roomID, tc.userIDs)

			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				mockUserRepo.AssertExpectations(t)
				mockRoomUserRepo.AssertExpectations(t)
			}
		})
	}
}
