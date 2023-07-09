package usecase

import (
	"context"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

//go:generate mockery --name=MessageUsecase --output=mocks
type MessageUsecase interface {
	GetMessagesByRoomID(ctx context.Context, roomId, lastEvaluatedKey string, limit int) ([]*model.Message, string, error)
	CreateMessage(ctx context.Context, message *model.Message) error
	UpdateMessage(ctx context.Context, roomId, messageId, newContent string) error
	DeleteMessage(ctx context.Context, roomId, messageId string) error
}
