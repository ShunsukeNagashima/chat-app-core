package repository

import (
	"context"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

//go:generate mockery --name=MessageRepository --output=mocks
type MessageRepository interface {
	GetAllMessagesByRoomID(ctx context.Context, roomId string) ([]*model.Message, error)
	Create(ctx context.Context, message *model.Message) error
	Update(ctx context.Context, messageId, newContent string) error
	Delete(ctx context.Context, messageId string) error
}
