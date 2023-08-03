package usecase

import (
	"context"
	"fmt"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/repository"
	"github.com/shunsukenagashima/chat-api/pkg/domain/usecase"
)

type RoomUserUsecaseImpl struct {
	roomUserRepo repository.RoomUserRepository
	userRepo     repository.UserRepository
	roomRepo     repository.RoomRepository
}

func NewRoomUserUsecase(roomUserRepo repository.RoomUserRepository, userRepo repository.UserRepository, roomRepo repository.RoomRepository) usecase.RoomUserUsecase {
	return &RoomUserUsecaseImpl{
		roomUserRepo,
		userRepo,
		roomRepo,
	}
}

func (ru *RoomUserUsecaseImpl) GetAllRoomsByUserID(ctx context.Context, userId string) ([]*model.Room, error) {
	roomUsers, err := ru.roomUserRepo.GetAllRoomsByUserID(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get all rooms by user ID: %w", err)
	}

	rooms := []*model.Room{}
	for _, roomUser := range roomUsers {
		room, err := ru.roomRepo.GetByID(ctx, roomUser.RoomID)
		if err != nil {
			return nil, fmt.Errorf("failed to get room by ID: %w", err)
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (ru *RoomUserUsecaseImpl) GetUsersByRoomID(ctx context.Context, roomId, lastEvaluatedKey string, limit int) ([]*model.User, string, error) {
	roomUsersers, nextKey, err := ru.roomUserRepo.GetUsersByRoomID(ctx, roomId, lastEvaluatedKey, limit)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get users by room ID: %w", err)
	}

	var users []*model.User
	for _, roomUser := range roomUsersers {
		user, err := ru.userRepo.GetByID(ctx, roomUser.UserID)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get user by ID: %w", err)
		}
		users = append(users, user)
	}

	return users, nextKey, nil
}

func (ru *RoomUserUsecaseImpl) RemoveUserFromRoom(ctx context.Context, roomId, userId string) error {
	if err := ru.roomUserRepo.RemoveUserFromRoom(ctx, roomId, userId); err != nil {
		return fmt.Errorf("failed to remove the user from the room: %w", err)
	}
	return nil
}

func (ru *RoomUserUsecaseImpl) AddUsersToRoom(ctx context.Context, roomId string, userIDs []string) error {
	for _, userId := range userIDs {
		_, err := ru.userRepo.GetByID(ctx, userId)
		if err != nil {
			return fmt.Errorf("failed to fetch the user with ID %s: %w", userId, err)
		}
	}

	room, err := ru.roomRepo.GetByID(ctx, roomId)
	if err != nil {
		return fmt.Errorf("failed to fetch the room with ID %s: %w", roomId, err)
	}
	if room == nil {
		return fmt.Errorf("room with the ID %s couldn't be found", roomId)
	}

	if err := ru.roomUserRepo.AddUsersToRoom(ctx, roomId, userIDs); err != nil {
		return fmt.Errorf("failed to add the users to the room: %w", err)
	}

	return nil
}
