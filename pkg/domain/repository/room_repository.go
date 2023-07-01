package repository

import (
	"context"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

//go:generate mockery --name=RoomRepository --output=mocks
type RoomRepository interface {
	GetByID(ctx context.Context, roomId string) (*model.Room, error)
	GetByName(ctx context.Context, name string) (*model.Room, error)
	GetAllPublic(ctx context.Context) ([]*model.Room, error)
	CreateAndAddUser(ctx context.Context, room *model.Room, ownerId string) error
	Delete(ctx context.Context, roomId string) error
	Update(ctx context.Context, room *model.Room) error
}
