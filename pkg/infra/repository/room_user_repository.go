package repository

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/repository"
)

type RoomUserRepositoryImpl struct {
	db     *dynamodb.DynamoDB
	dbName string
}

func NewRoomUserRepository(db *dynamodb.DynamoDB) repository.RoomUserRepository {
	return &RoomUserRepositoryImpl{
		db,
		"RoomUsers",
	}
}

func (r *RoomUserRepositoryImpl) GetAllRoomsByUserID(ctx context.Context, userId string) ([]*model.RoomUser, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String("RoomUsers"),
		IndexName: aws.String("UserIDIndex"),
		KeyConditions: map[string]*dynamodb.Condition{
			"userId": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(userId),
					},
				},
			},
		},
	}

	result, err := r.db.Query(input)
	if err != nil {
		return nil, err
	}
	var roomUsers []*model.RoomUser
	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &roomUsers); err != nil {
		return nil, err
	}

	return roomUsers, nil
}

func (r *RoomUserRepositoryImpl) RemoveUserFromRoom(ctx context.Context, roomId, userId string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.dbName),
		Key: map[string]*dynamodb.AttributeValue{
			"roomId": {
				S: aws.String(roomId),
			},
			"userId": {
				S: aws.String(userId),
			},
		},
	}

	_, err := r.db.DeleteItem(input)
	return err
}

func (r *RoomUserRepositoryImpl) AddUsersToRoom(ctx context.Context, roomId string, userIDs []string) error {
	writeRequests := make([]*dynamodb.TransactWriteItem, len(userIDs))
	for i, userId := range userIDs {
		writeRequests[i] = &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				TableName: aws.String(r.dbName),
				Item: map[string]*dynamodb.AttributeValue{
					"roomId": {
						S: aws.String(roomId),
					},
					"userId": {
						S: aws.String(userId),
					},
				},
			},
		}
	}
	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: writeRequests,
	}
	_, err := r.db.TransactWriteItems(input)
	return err
}
