package usecase

import (
	"context"
	"errors"

	"github.com/shunsukenagashima/chat-api/pkg/apperror"
	"github.com/shunsukenagashima/chat-api/pkg/clock"
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

func (uu *UserUsecaseImpl) CreateUser(ctx context.Context, user *model.User, idToken string) error {
	// TODO: NAT Gateway is required to use Firebase Auth
	// token, err := uu.firebaseAuth.GetFirebaseUser(ctx, idToken)
	// if err != nil {
	// 	return err
	// }
	// if token.UID != user.UserID {
	// 	return fmt.Errorf("provided user ID does not match the user ID in Firebase token")
	// }

	_, err := uu.repo.GetByID(ctx, user.UserID)
	if err != nil {
		var notFoundErr *apperror.NotFoundErr
		if !errors.As(err, &notFoundErr) {
			return err
		}
	}

	clock := clock.RealClocker{}
	user.CreatedAt = clock.Now()

	return uu.repo.Create(ctx, user)
}

func (uu *UserUsecaseImpl) GetMultipleUsers(ctx context.Context, lastEvaluatedKey string, limit int) ([]*model.User, string, error) {
	return uu.repo.GetMultiple(ctx, lastEvaluatedKey, limit)
}

func (uu *UserUsecaseImpl) GetUserByID(ctx context.Context, userId string) (*model.User, error) {
	return uu.repo.GetByID(ctx, userId)
}

func (uu *UserUsecaseImpl) BatchGetUsers(ctx context.Context, userIds []string) ([]*model.User, error) {
	return uu.repo.BatchGetUsers(ctx, userIds)
}
