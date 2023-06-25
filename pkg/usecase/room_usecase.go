package usecase

import (
	"context"
	"fmt"

	"github.com/shunsukenagashima/chat-api/pkg/apperror"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/repository"
	"github.com/shunsukenagashima/chat-api/pkg/domain/usecase"
)

type RoomUsecaseImpl struct {
	roomRepo repository.RoomRepository
	userRepo repository.UserRepository
}

func NewRoomUsecase(roomRepo repository.RoomRepository, userRepo repository.UserRepository) usecase.RoomUsecase {
	return &RoomUsecaseImpl{
		roomRepo,
		userRepo,
	}
}

func (ru *RoomUsecaseImpl) GetRoomByID(ctx context.Context, roomId string) (*model.Room, error) {
	return ru.roomRepo.GetById(ctx, roomId)
}

func (ru *RoomUsecaseImpl) GetAllPublicRoom(ctx context.Context) ([]*model.Room, error) {
	return ru.roomRepo.GetAllPublic(ctx)
}

func (ru *RoomUsecaseImpl) CreateRoom(ctx context.Context, room *model.Room, ownerID string) error {
	existingRoom, err := ru.roomRepo.GetByName(ctx, room.Name)
	if err != nil {
		return err
	}
	if existingRoom != nil {
		return apperror.NewAlreadyExistsErr("Room", "RoomName: "+room.Name)
	}

	owner, err := ru.userRepo.GetByID(ctx, ownerID)
	if err != nil {
		return err
	}
	if owner == nil {
		return apperror.NewNotFoundErr("User", "UserID: "+ownerID)
	}

	if err := ru.roomRepo.CreateAndAddUser(ctx, room, ownerID); err != nil {
		return err
	}

	return nil
}

func (ru *RoomUsecaseImpl) DeleteRoom(ctx context.Context, roomId string) error {
	room, err := ru.roomRepo.GetById(ctx, roomId)
	if err != nil {
		return err
	}
	if room == nil {
		return fmt.Errorf("room with the ID '%s' couldn't be found", roomId)
	}

	if err := ru.roomRepo.Delete(ctx, roomId); err != nil {
		return err
	}

	return nil
}

func (ru *RoomUsecaseImpl) UpdateRoom(ctx context.Context, room *model.Room) error {
	existingRoom, err := ru.roomRepo.GetByName(ctx, room.Name)
	if err != nil {
		return err
	}
	if existingRoom != nil && existingRoom.RoomID != room.RoomID {
		return fmt.Errorf("name of chat room '%s' is duplicated", room.Name)
	}

	if err := ru.roomRepo.Update(ctx, room); err != nil {
		return err
	}

	return nil
}
