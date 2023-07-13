package repository

import (
	"context"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

//go:generate mockery --name=MessageRepository --output=mocks
type MessageRepository interface {
	GetMessagesByRoomID(ctx context.Context, roomId, lastEvaluatedKey string, limit int) ([]*model.Message, string, error)
	GetByID(ctx context.Context, roomId, messageId string) (*model.Message, error)
	Create(ctx context.Context, message *model.Message) error
	Update(ctx context.Context, roomId, messageId, newContent string) error
	Delete(ctx context.Context, roomId, messageId string) error
}
