package repository

import "github.com/shunsukenagashima/chat-api/pkg/domain/model"

type RoomsRepository interface {
	Get(roomID string) (*model.Room, error)
	GetAllPublic() ([]*model.Room, error)
	Create(room *model.Room) error
	Update(room *model.Room) error
	Delete(roomID string) error
}
