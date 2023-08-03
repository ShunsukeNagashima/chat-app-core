package repository

import (
	"context"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

//go:generate mockery --name=RoomUserRepository --output=mocks
type RoomUserRepository interface {
	GetAllRoomsByUserID(ctx context.Context, userId string) ([]*model.RoomUser, error)
	GetUsersByRoomID(ctx context.Context, roomId, lastEvaluatedKey string, limit int) ([]*model.RoomUser, string, error)
	RemoveUserFromRoom(ctx context.Context, roomId, userId string) error
	AddUsersToRoom(ctx context.Context, roomId string, userIDs []string) error
}
