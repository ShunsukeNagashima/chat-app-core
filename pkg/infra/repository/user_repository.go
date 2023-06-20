package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/repository"
)

type UserRepositoryImpl struct {
	db     *dynamodb.DynamoDB
	dbName string
}

func NewUserRepository(db *dynamodb.DynamoDB) repository.UserRepository {
	return &UserRepositoryImpl{
		db,
		"Users",
	}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.dbName),
		Item:      item,
	}

	_, err = r.db.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, userID string) (*model.User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.dbName),
		Key: map[string]*dynamodb.AttributeValue{
			"userID": {
				S: aws.String(userID),
			},
		},
	}

	result, err := r.db.GetItem(input)

	if err != nil {
		return nil, err
	}

	if len(result.Item) == 0 {
		return nil, fmt.Errorf("user with ID %s not found", userID)
	}

	var user model.User
	if err := dynamodbattribute.UnmarshalMap(result.Item, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
