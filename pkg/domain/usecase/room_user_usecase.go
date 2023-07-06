package usecase

import (
	"context"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

//go:generate mockery --name=RoomUserUsecase --output=mocks
type RoomUserUsecase interface {
	GetAllRoomsByUserID(ctx context.Context, userId string) ([]*model.Room, error)
	RemoveUserFromRoom(ctx context.Context, roomId, userId string) error
	AddUsersToRoom(ctx context.Context, roomId string, userIds []string) error
}
