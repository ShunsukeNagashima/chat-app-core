package repository

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/shunsukenagashima/chat-api/pkg/apperror"
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

func (r *UserRepositoryImpl) GetMultiple(ctx context.Context, lastEvaluatedKey string, limit int) ([]*model.User, string, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.dbName),
		Limit:     aws.Int64(int64(limit)),
	}

	if lastEvaluatedKey != "" {
		input.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"userId": {
				S: aws.String(lastEvaluatedKey),
			},
		}
	}

	result, err := r.db.Scan(input)
	if err != nil {
		return nil, "", err
	}

	var users []*model.User
	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &users); err != nil {
		return nil, "", err
	}

	var nextKey string
	if result.LastEvaluatedKey != nil {
		nextKey = *result.LastEvaluatedKey["userId"].S
	}

	return users, nextKey, nil
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, userId string) (*model.User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.dbName),
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {
				S: aws.String(userId),
			},
		},
	}

	result, err := r.db.GetItem(input)

	if err != nil {
		return nil, err
	}

	if len(result.Item) == 0 {
		return nil, apperror.NewNotFoundErr("User", "UserID: "+userId)
	}

	var user model.User
	if err := dynamodbattribute.UnmarshalMap(result.Item, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryImpl) BatchGetUsers(ctx context.Context, userIds []string) ([]*model.User, error) {
	var keys []map[string]*dynamodb.AttributeValue
	for _, userId := range userIds {
		keys = append(keys, map[string]*dynamodb.AttributeValue{
			"userId": {
				S: aws.String(userId),
			},
		})
	}

	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			r.dbName: {
				Keys: keys,
			},
		},
	}

	result, err := r.db.BatchGetItem(input)
	if err != nil {
		return nil, err
	}

	var users []*model.User
	for _, item := range result.Responses[r.dbName] {
		var user model.User
		if err := dynamodbattribute.UnmarshalMap(item, &user); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}
