package usecase

import (
	"context"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

//go:generate mockery --name=RoomUserUsecase --output=mocks
type RoomUserUsecase interface {
	GetAllRoomsByUserID(ctx context.Context, userID string) ([]*model.Room, error)
	AddUserToRoom(ctx context.Context, roomID, userID string) error
	RemoveUserFromRoom(ctx context.Context, roomID, userID string) error
	AddUsersToRoom(ctx context.Context, roomID string, userIDs []string) error
}
