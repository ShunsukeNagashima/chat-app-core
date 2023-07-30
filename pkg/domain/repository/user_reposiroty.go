package repository

import (
	"context"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

//go:generate mockery --name=UserRepository --output=mocks
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, userId string) (*model.User, error)
	SearchUsers(ctx context.Context, query, nextKey string, size int) ([]*model.User, string, error)
	BatchGetUsers(ctx context.Context, userIds []string) ([]*model.User, error)
}
