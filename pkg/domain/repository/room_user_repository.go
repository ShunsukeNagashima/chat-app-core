package repository

import "github.com/shunsukenagashima/chat-api/pkg/domain/model"

type RoomUserRepository interface {
	GetAllByUserID(userID string) ([]*model.Room, error)
}
