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
	mockRoomRepo := new(mocks.RoomRepository)
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

	mockRoomUsers := []*model.RoomUser{
		{
			RoomID: "1",
			UserID: "1",
		},
		{
			RoomID: "2",
			UserID: "1",
		},
	}

	mockRoomUserRepo.On("GetAllRoomsByUserID", mock.Anything, mock.Anything).Return(mockRoomUsers, nil)
	mockRoomRepo.On("GetByID", mock.Anything, mock.Anything).Return(mockRooms[0], nil).Once()
	mockRoomRepo.On("GetByID", mock.Anything, mock.Anything).Return(mockRooms[1], nil).Once()

	roomUserUsecase := NewRoomUserUsecase(mockRoomUserRepo, mockUserRepo, mockRoomRepo)

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

func TestRemoveUserFromRoom(t *testing.T) {
	mockRoomUserRepo := new(mocks.RoomUserRepository)
	mockUserRepo := new(mocks.UserRepository)
	mockRoomRepo := new(mocks.RoomRepository)

	roomId := "1"
	userId := "1"

	mockRoomUserRepo.On("RemoveUserFromRoom", mock.Anything, roomId, userId).Return(nil)

	roomUserUsecase := NewRoomUserUsecase(mockRoomUserRepo, mockUserRepo, mockRoomRepo)

	err := roomUserUsecase.RemoveUserFromRoom(context.Background(), roomId, userId)

	assert.NoError(t, err)
	mockRoomUserRepo.AssertExpectations(t)
}

func TestAddUsersToRoom(t *testing.T) {
	mockRoomUserRepo := new(mocks.RoomUserRepository)
	mockUserRepo := new(mocks.UserRepository)
	mockRoomRepo := new(mocks.RoomRepository)
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

	mockRoom := &model.Room{
		RoomID:   "1",
		Name:     "room-1",
		RoomType: model.Public,
	}

	testCases := []struct {
		name        string
		userIDs     []string
		roomId      string
		expectedErr error
	}{
		{
			name:        "Success",
			userIDs:     []string{mockUsers[0].UserID, mockUsers[1].UserID, mockUsers[2].UserID},
			roomId:      "1",
			expectedErr: nil,
		},
		{
			name:        "OneofUsersNotFound",
			userIDs:     []string{mockUsers[0].UserID, mockUsers[1].UserID, "invalid_user_id"},
			roomId:      "1",
			expectedErr: errors.New("user not found"),
		},
		{
			name:        "RoomNotFound",
			userIDs:     []string{mockUsers[0].UserID, mockUsers[1].UserID, mockUsers[2].UserID},
			roomId:      "invalid_room_id",
			expectedErr: errors.New("room not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			for i, userId := range tc.userIDs {
				if userId == "invalid_user_id" {
					mockUserRepo.On("GetByID", mock.Anything, userId).Return(nil, errors.New("user not found"))
				} else {
					mockUserRepo.On("GetByID", mock.Anything, userId).Return(mockUsers[i], nil)
				}
			}
			mockRoomUserRepo.On("AddUsersToRoom", mock.Anything, tc.roomId, tc.userIDs).Return(nil)

			if tc.roomId == "invalid_room_id" {
				mockRoomRepo.On("GetByID", mock.Anything, tc.roomId).Return(nil, errors.New("room not found"))
			} else {
				mockRoomRepo.On("GetByID", mock.Anything, tc.roomId).Return(mockRoom, nil)
			}

			roomUserUsecase := NewRoomUserUsecase(mockRoomUserRepo, mockUserRepo, mockRoomRepo)

			err := roomUserUsecase.AddUsersToRoom(context.Background(), tc.roomId, tc.userIDs)

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
