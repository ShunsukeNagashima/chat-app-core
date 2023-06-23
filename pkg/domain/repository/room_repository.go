package repository

import (
	"context"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

//go:generate mockery --name=RoomRepository --output=mocks
type RoomRepository interface {
	GetById(ctx context.Context, roomID string) (*model.Room, error)
	GetByName(ctx context.Context, name string) (*model.Room, error)
	GetAllPublic(ctx context.Context) ([]*model.Room, error)
	CreateAndAddUser(ctx context.Context, room *model.Room, ownerID string) error
	Delete(ctx context.Context, roomID string) error
	Update(ctx context.Context, room *model.Room) error
}
