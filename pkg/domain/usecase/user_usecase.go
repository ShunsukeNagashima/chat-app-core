package usecase

import (
	"context"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

//go:generate mockery --name=UserUsecase --output=mocks
type UserUsecase interface {
	CreateUser(ctx context.Context, user *model.User, idToken string) error
	GetUserByID(ctx context.Context, userId string) (*model.User, error)
	SearchUsers(ctx context.Context, query string, from, size int) ([]*model.User, error)
}
