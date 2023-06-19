package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetRoomByID(t *testing.T) {
	mockRepo := new(mocks.RoomRepository)
	mockRoom := &model.Room{
		RoomID:   "1",
		Name:     "Room1",
		RoomType: model.Public,
	}

	mockRepo.On("GetById", mock.Anything, mockRoom.RoomID).Return(mockRoom, nil)
	roomUsecase := NewRoomUsecase(mockRepo)

	room, err := roomUsecase.GetRoomByID(context.Background(), mockRoom.RoomID)

	assert.NoError(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, mockRoom.RoomID, room.RoomID)
	assert.Equal(t, mockRoom.Name, room.Name)
	assert.Equal(t, mockRoom.RoomType, room.RoomType)
	mockRepo.AssertExpectations(t)
}

func TestGetAllPublicRoom(t *testing.T) {
	mockRepo := new(mocks.RoomRepository)
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
	roomUsecase := NewRoomUsecase(mockRepo)

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
	mockRepo := new(mocks.RoomRepository)
	mockRoom := &model.Room{
		RoomID:   "1",
		Name:     "Room1",
		RoomType: model.Public,
	}

	mockRepo.On("GetByName", mock.Anything, mockRoom.Name).Return(nil, nil)

	mockRepo.On("Create", mock.Anything, mockRoom).Return(nil)

	roomUsecase := NewRoomUsecase(mockRepo)

	err := roomUsecase.CreateRoom(context.Background(), mockRoom)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateRoom_AlreadyExists(t *testing.T) {
	mockRepo := new(mocks.RoomRepository)
	mockRoom := &model.Room{
		RoomID:   "1",
		Name:     "Room1",
		RoomType: model.Public,
	}

	mockRepo.On("GetByName", mock.Anything, mockRoom.Name).Return(mockRoom, nil)

	roomUsecase := NewRoomUsecase(mockRepo)

	err := roomUsecase.CreateRoom(context.Background(), mockRoom)

	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Errorf("name of chat room '%s' is duplicated", mockRoom.Name).Error())
	mockRepo.AssertExpectations(t)
}

func TestDeleteRoom(t *testing.T) {
	mockRepo := new(mocks.RoomRepository)
	roomID := "1"

	mockRoom := &model.Room{
		RoomID:   roomID,
		Name:     "Room1",
		RoomType: model.Public,
	}

	mockRepo.On("GetById", mock.Anything, roomID).Return(mockRoom, nil)

	mockRepo.On("Delete", mock.Anything, roomID).Return(nil)

	roomUsecase := NewRoomUsecase(mockRepo)

	err := roomUsecase.DeleteRoom(context.Background(), roomID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteRoom_NotExist(t *testing.T) {
	mockRepo := new(mocks.RoomRepository)
	roomID := "1"

	mockRepo.On("GetById", mock.Anything, roomID).Return(nil, nil)

	roomUsecase := NewRoomUsecase(mockRepo)

	err := roomUsecase.DeleteRoom(context.Background(), roomID)

	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Errorf("room with the ID %s couldn't be found", roomID).Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateRoom(t *testing.T) {
	mockRepo := new(mocks.RoomRepository)
	roomID := "1"
	roomName := "Room1"

	mockRoom := &model.Room{
		RoomID:   roomID,
		Name:     roomName,
		RoomType: model.Public,
	}

	mockRepo.On("GetByName", mock.Anything, roomName).Return(nil, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Room")).Return(nil)

	roomUsecase := NewRoomUsecase(mockRepo)

	err := roomUsecase.UpdateRoom(context.Background(), mockRoom)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateRoom_DuplicateName(t *testing.T) {
	mockRepo := new(mocks.RoomRepository)
	roomID := "1"
	roomName := "Room1"

	mockRoom := &model.Room{
		RoomID:   roomID,
		Name:     roomName,
		RoomType: model.Public,
	}

	mockRepo.On("GetByName", mock.Anything, roomName).Return(&model.Room{RoomID: "2", Name: roomName, RoomType: model.Public}, nil)

	roomUsecase := NewRoomUsecase(mockRepo)

	err := roomUsecase.UpdateRoom(context.Background(), mockRoom)

	assert.Error(t, err)
	assert.EqualError(t, err, fmt.Errorf("name of chat room '%s' is duplicated", roomName).Error())
	mockRepo.AssertExpectations(t)
}
