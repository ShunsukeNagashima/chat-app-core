package usecase

import (
	"context"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

//go:generate mockery --name=MessageUsecase --output=mocks
type MessageUsecase interface {
	GetAllMessagesByRoomID(ctx context.Context, roomId string) ([]*model.Message, error)
	CreateMessage(ctx context.Context, message *model.Message) error
	UpdateMessage(ctx context.Context, messageId, newContent string) error
	DeleteMessage(ctx context.Context, messageId string) error
}
