package usecase

import (
	"context"
	"errors"

	"github.com/shunsukenagashima/chat-api/pkg/apperror"
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

func (mu *MessageUsecaseImpl) GetAllMessagesByRoomID(ctx context.Context, roomId string) ([]*model.Message, error) {
	return mu.messageRepo.GetAllByRoomID(ctx, roomId)
}

func (mu *MessageUsecaseImpl) CreateMessage(ctx context.Context, message *model.Message) error {
	return mu.messageRepo.Create(ctx, message)
}

func (mu *MessageUsecaseImpl) UpdateMessage(ctx context.Context, messageId, newContent string) error {
	_, err := mu.messageRepo.GetByID(ctx, messageId)
	if err != nil {
		var notFoundErr *apperror.NotFoundErr
		if !errors.As(err, &notFoundErr) {
			return err
		}
	}

	return mu.messageRepo.Update(ctx, messageId, newContent)
}

func (mu *MessageUsecaseImpl) DeleteMessage(ctx context.Context, messageId string) error {
	_, err := mu.messageRepo.GetByID(ctx, messageId)
	if err != nil {
		var notFoundErr *apperror.NotFoundErr
		if !errors.As(err, &notFoundErr) {
			return err
		}
	}

	return mu.messageRepo.Delete(ctx, messageId)
}
