package auth

import (
	"context"

	"firebase.google.com/go/auth"
)

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
