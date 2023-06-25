package repository

import (
	"context"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

//go:generate mockery --name=RoomUserRepository --output=mocks
type RoomUserRepository interface {
	GetAllRoomsByUserID(ctx context.Context, userID string) ([]*model.RoomUser, error)
	RemoveUserFromRoom(ctx context.Context, roomID, userID string) error
	AddUsersToRoom(ctx context.Context, roomID string, userIDs []string) error
}
