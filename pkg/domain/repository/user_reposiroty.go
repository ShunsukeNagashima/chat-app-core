package repository

import (
	"context"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

//go:generate mockery --name=UserRepository --output=mocks
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetMultiple(ctx context.Context, lastEvaluatedKey string, limit int) ([]*model.User, string, error)
	GetByID(ctx context.Context, userId string) (*model.User, error)
	BatchGetUsers(ctx context.Context, userIds []string) ([]*model.User, error)
}
