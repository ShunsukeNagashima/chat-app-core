package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/repository"
	"github.com/shunsukenagashima/chat-api/pkg/domain/usecase"
)

type RoomUserUsecaseImpl struct {
	roomUserRepo repository.RoomUserRepository
	userRepo     repository.UserRepository
}

func NewRoomUserUsecase(roomUserRepo repository.RoomUserRepository, userRepo repository.UserRepository) usecase.RoomUserUsecase {
	return &RoomUserUsecaseImpl{
		roomUserRepo,
		userRepo,
	}
}

func (ru *RoomUserUsecaseImpl) GetAllRoomsByUserID(ctx context.Context, userID string) ([]*model.Room, error) {
	roomUsers, err := ru.roomUserRepo.GetAllRoomsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all rooms by user ID: %w", err)
	}

	return roomUsers, nil
}

func (ru *RoomUserUsecaseImpl) AddUserToRoom(ctx context.Context, roomID, userID string) error {
	user, err := ru.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	if err := ru.roomUserRepo.AddUserToRoom(ctx, roomID, userID); err != nil {
		return fmt.Errorf("failed to add the user to the room: %w", err)
	}
	return nil
}

func (ru *RoomUserUsecaseImpl) RemoveUserFromRoom(ctx context.Context, roomID, userID string) error {
	if err := ru.roomUserRepo.RemoveUserFromRoom(ctx, roomID, userID); err != nil {
		return fmt.Errorf("failed to remove the user from the room: %w", err)
	}
	return nil
}

func (ru *RoomUserUsecaseImpl) AddUsersToRoom(ctx context.Context, roomID string, userIDs []string) error {
	for _, userID := range userIDs {
		_, err := ru.userRepo.GetByID(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to fetch the user with ID %s: %w", userID, err)
		}
	}

	if err := ru.roomUserRepo.AddUsersToRoom(ctx, roomID, userIDs); err != nil {
		return fmt.Errorf("failed to add the users to the room: %w", err)
	}

	return nil
}
