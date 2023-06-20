package usecase

import (
	"context"
	"fmt"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/repository"
	"github.com/shunsukenagashima/chat-api/pkg/domain/usecase"
)

type RoomUsecaseImpl struct {
	repo repository.RoomRepository
}

func NewRoomUsecase(repo repository.RoomRepository) usecase.RoomUsecase {
	return &RoomUsecaseImpl{
		repo,
	}
}

func (ru *RoomUsecaseImpl) GetRoomByID(ctx context.Context, roomID string) (*model.Room, error) {
	return ru.repo.GetById(ctx, roomID)
}

func (ru *RoomUsecaseImpl) GetAllPublicRoom(ctx context.Context) ([]*model.Room, error) {
	return ru.repo.GetAllPublic(ctx)
}

func (ru *RoomUsecaseImpl) CreateRoom(ctx context.Context, room *model.Room) error {
	existingRoom, err := ru.repo.GetByName(ctx, room.Name)
	if err != nil {
		return err
	}
	if existingRoom != nil {
		return fmt.Errorf("name of chat room '%s' is duplicated", room.Name)
	}

	if err := ru.repo.Create(ctx, room); err != nil {
		return err
	}

	return nil
}

func (ru *RoomUsecaseImpl) DeleteRoom(ctx context.Context, roomID string) error {
	room, err := ru.repo.GetById(ctx, roomID)
	if err != nil {
		return err
	}
	if room == nil {
		return fmt.Errorf("room with the ID %s couldn't be found", roomID)
	}

	if err := ru.repo.Delete(ctx, roomID); err != nil {
		return err
	}

	return nil
}

func (ru *RoomUsecaseImpl) UpdateRoom(ctx context.Context, room *model.Room) error {
	existingRoom, err := ru.repo.GetByName(ctx, room.Name)
	if err != nil {
		return err
	}
	if existingRoom != nil && existingRoom.RoomID != room.RoomID {
		return fmt.Errorf("name of chat room '%s' is duplicated", room.Name)
	}

	if err := ru.repo.Update(ctx, room); err != nil {
		return err
	}

	return nil
}
