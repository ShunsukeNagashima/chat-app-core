package usecase

import (
	"context"
	"fmt"

	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/repository"
	"github.com/shunsukenagashima/chat-api/pkg/domain/usecase"
	"github.com/shunsukenagashima/chat-api/pkg/infra/auth"
)

type UserUsecaseImpl struct {
	repo         repository.UserRepository
	firebaseAuth auth.FirebaseAuthenticator
}

func NewUserUsecase(repo repository.UserRepository, firebaseAuth auth.FirebaseAuthenticator) usecase.UserUsecase {
	return &UserUsecaseImpl{
		repo,
		firebaseAuth,
	}
}

func (uu *UserUsecaseImpl) CreateUser(ctx context.Context, user *model.User) error {
	token, err := uu.firebaseAuth.GetFirebaseUser(ctx, user.UserID)
	if err != nil {
		return err
	}
	if token.UID != user.UserID {
		return fmt.Errorf("provided user ID does not match the user ID in Firebase token")
	}

	return uu.repo.Create(ctx, user)
}

func (uu *UserUsecaseImpl) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	return uu.repo.GetByID(ctx, userID)
}
