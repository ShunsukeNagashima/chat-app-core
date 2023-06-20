package auth

import (
	"context"

	"firebase.google.com/go/auth"
)

//go:generate mockery --name=FirebaseAuthenticator --output=mocks
type FirebaseAuthenticator interface {
	GetFirebaseUser(ctx context.Context, uid string) (*auth.Token, error)
}

type FirebaseAuth struct {
	client *auth.Client
}

func NewFirebaseAuth(authClient *auth.Client) *FirebaseAuth {
	return &FirebaseAuth{
		client: authClient,
	}
}

func (fa *FirebaseAuth) GetFirebaseUser(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := fa.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}
