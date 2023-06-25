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

func TestGetRoomByID(t *testing.T) {
	mockRoomRepo := new(mocks.RoomRepository)
	mockUserRepo := new(mocks.UserRepository)
	mockRoom := &model.Room{
		RoomID:   "1",
		Name:     "Room1",
		RoomType: model.Public,
	}

	mockRoomRepo.On("GetById", mock.Anything, mockRoom.RoomID).Return(mockRoom, nil)
	roomUsecase := NewRoomUsecase(mockRoomRepo, mockUserRepo)

	room, err := roomUsecase.GetRoomByID(context.Background(), mockRoom.RoomID)

	assert.NoError(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, mockRoom.RoomID, room.RoomID)
	assert.Equal(t, mockRoom.Name, room.Name)
	assert.Equal(t, mockRoom.RoomType, room.RoomType)
	mockRoomRepo.AssertExpectations(t)
}

func TestGetAllPublicRoom(t *testing.T) {
	mockRepo := new(mocks.RoomRepository)
	mockUserRepo := new(mocks.UserRepository)
	mockRooms := []*model.Room{
		{
			RoomID:   "1",
			Name:     "Room1",
			RoomType: model.Public,
		},
		{
			RoomID:   "2",
			Name:     "Room2",
			RoomType: model.Public,
		},
	}

	mockRepo.On("GetAllPublic", mock.Anything).Return(mockRooms, nil)
	roomUsecase := NewRoomUsecase(mockRepo, mockUserRepo)

	rooms, err := roomUsecase.GetAllPublicRoom(context.Background())

	assert.NoError(t, err)
	assert.NotEmpty(t, rooms)
	assert.Equal(t, len(mockRooms), len(rooms))
	for i, room := range rooms {
		assert.Equal(t, mockRooms[i].RoomID, room.RoomID)
		assert.Equal(t, mockRooms[i].Name, room.Name)
		assert.Equal(t, mockRooms[i].RoomType, room.RoomType)
	}
	mockRepo.AssertExpectations(t)
}

func TestCreateRoom(t *testing.T) {
	mockRoom := &model.Room{
		RoomID:   "1",
		Name:     "Room1",
		RoomType: model.Public,
	}

	mockUser := &model.User{
		UserID:   "1",
		Username: "user-1",
		Email:    "user-1@example.com",
	}

	testCases := []struct {
		name                   string
		room                   *model.Room
		ownerID                string
		mockRoomRepoReturn     *model.Room
		mockUserRepoReturn     *model.User
		mockRoomUserRepoReturn error
		expectedErr            error
	}{
		{
			name:                   "Success",
			room:                   mockRoom,
			ownerID:                "1",
			mockRoomRepoReturn:     nil,
			mockUserRepoReturn:     mockUser,
			mockRoomUserRepoReturn: nil,
			expectedErr:            nil,
		},
		{
			name:                   "Invalid OwnerID",
			room:                   mockRoom,
			ownerID:                "invalid_ownerID",
			mockRoomRepoReturn:     nil,
			mockUserRepoReturn:     nil,
			mockRoomUserRepoReturn: nil,
			expectedErr:            errors.New("user with the ID 'invalid_ownerID' couldn't be found"),
		},
		{
			name:                   "Duplicated Room Name",
			room:                   mockRoom,
			ownerID:                "1",
			mockRoomRepoReturn:     mockRoom,
			mockUserRepoReturn:     nil,
			mockRoomUserRepoReturn: nil,
			expectedErr:            errors.New("name of chat room 'Room1' is duplicated"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRoomRepo := new(mocks.RoomRepository)
			mockUserRepo := new(mocks.UserRepository)

			mockRoomRepo.On("GetByName", mock.Anything, tc.room.Name).Return(tc.mockRoomRepoReturn, nil)
			mockUserRepo.On("GetByID", mock.Anything, tc.ownerID).Return(tc.mockUserRepoReturn, nil)
			mockRoomRepo.On("CreateAndAddUser", mock.Anything, tc.room, tc.ownerID).Return(nil)

			roomUsecase := NewRoomUsecase(mockRoomRepo, mockUserRepo)

			err := roomUsecase.CreateRoom(context.Background(), tc.room, tc.ownerID)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				mockRoomRepo.AssertExpectations(t)
				mockUserRepo.AssertExpectations(t)
			}
		})
	}
}

func TestDeleteRoom(t *testing.T) {
	roomId := "1"

	mockRoom := &model.Room{
		RoomID:   roomId,
		Name:     "Room1",
		RoomType: model.Public,
	}

	testCases := []struct {
		name               string
		roomId             string
		mockRoomRepoReturn *model.Room
		expectedErr        error
	}{
		{
			name:               "Success",
			roomId:             roomId,
			mockRoomRepoReturn: mockRoom,
			expectedErr:        nil,
		},
		{
			name:               "Invalid RoomID",
			roomId:             "invalid_roomID",
			mockRoomRepoReturn: nil,
			expectedErr:        errors.New("room with the ID 'invalid_roomID' couldn't be found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRoomRepo := new(mocks.RoomRepository)
			mockUserRepo := new(mocks.UserRepository)

			mockRoomRepo.On("GetById", mock.Anything, tc.roomId).Return(tc.mockRoomRepoReturn, nil)
			mockRoomRepo.On("Delete", mock.Anything, tc.roomId).Return(nil)

			roomUsecase := NewRoomUsecase(mockRoomRepo, mockUserRepo)

			err := roomUsecase.DeleteRoom(context.Background(), tc.roomId)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				mockRoomRepo.AssertExpectations(t)
			}
		})
	}
}

func TestUpdateRoom(t *testing.T) {
	mockRoom := &model.Room{
		RoomID:   "1",
		Name:     "Room1",
		RoomType: model.Public,
	}

	testCases := []struct {
		name                string
		room                *model.Room
		expectedErr         error
		mockGetByNameReturn *model.Room
	}{
		{
			name: "Success",
			room: &model.Room{
				RoomID:   mockRoom.RoomID,
				Name:     mockRoom.Name + "updated",
				RoomType: model.Private,
			},
			expectedErr:         nil,
			mockGetByNameReturn: nil,
		},
		{
			name: "Duplicated Room Name With Different RoomID",
			room: &model.Room{
				RoomID:   "2",
				Name:     mockRoom.Name,
				RoomType: model.Private,
			},
			expectedErr:         errors.New("name of chat room 'Room1' is duplicated"),
			mockGetByNameReturn: mockRoom,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRoomRepo := new(mocks.RoomRepository)
			mockUserRepo := new(mocks.UserRepository)

			mockRoomRepo.On("GetByName", mock.Anything, tc.room.Name).Return(tc.mockGetByNameReturn, nil)
			mockRoomRepo.On("Update", mock.Anything, tc.room).Return(nil)

			roomUsecase := NewRoomUsecase(mockRoomRepo, mockUserRepo)

			err := roomUsecase.UpdateRoom(context.Background(), tc.room)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				mockRoomRepo.AssertExpectations(t)
			}
		})
	}
}
