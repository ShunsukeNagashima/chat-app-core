package usecase

import (
	"context"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

//go:generate mockery --name=RoomUsecase --output=mocks
type RoomUsecase interface {
	GetRoomByID(ctx context.Context, roomId string) (*model.Room, error)
	GetAllPublicRoom(ctx context.Context) ([]*model.Room, error)
	CreateRoom(ctx context.Context, room *model.Room, ownerId string) error
	DeleteRoom(ctx context.Context, roomId string) error
	UpdateRoom(ctx context.Context, room *model.Room) error
}
