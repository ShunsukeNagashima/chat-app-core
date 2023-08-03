package usecase

import (
	"context"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

//go:generate mockery --name=UserUsecase --output=mocks
type UserUsecase interface {
	CreateUser(ctx context.Context, user *model.User, idToken string) error
	GetMultipleUsers(ctx context.Context, lastEvaluatedKey string, limit int) ([]*model.User, string, error)
	GetUserByID(ctx context.Context, userId string) (*model.User, error)
	BatchGetUsers(ctx context.Context, userIds []string) ([]*model.User, error)
}
