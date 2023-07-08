package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/shunsukenagashima/chat-api/pkg/clock"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/repository"
	"github.com/shunsukenagashima/chat-api/pkg/domain/usecase"
)

type MessageUsecaseImpl struct {
	messageRepo repository.MessageRepository
}

func NewMessageUsecase(messageRepo repository.MessageRepository) usecase.MessageUsecase {
	return &MessageUsecaseImpl{
		messageRepo: messageRepo,
	}
}

func (mu *MessageUsecaseImpl) GetMessagesByRoomID(ctx context.Context, roomId, lastEvaluatedKey string, limit int) ([]*model.Message, string, error) {
	return mu.messageRepo.GetMessagesByRoomID(ctx, roomId, lastEvaluatedKey, limit)
}

func (mu *MessageUsecaseImpl) CreateMessage(ctx context.Context, message *model.Message) error {
	clock := clock.RealClocker{}

	message.MessageID = uuid.New().String()
	message.CreatedAt = clock.Now()

	return mu.messageRepo.Create(ctx, message)
}

func (mu *MessageUsecaseImpl) UpdateMessage(ctx context.Context, messageId, newContent string) error {
	_, err := mu.messageRepo.GetByID(ctx, messageId)
	if err != nil {
		return err
	}

	return mu.messageRepo.Update(ctx, messageId, newContent)
}

func (mu *MessageUsecaseImpl) DeleteMessage(ctx context.Context, messageId string) error {
	_, err := mu.messageRepo.GetByID(ctx, messageId)
	if err != nil {
		return err
	}

	return mu.messageRepo.Delete(ctx, messageId)
}
