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

func (r *RoomUserRepositoryImpl) GetAllRoomsByUserID(ctx context.Context, userID string) ([]*model.Room, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String("RoomUsers"),
		IndexName: aws.String("UserIDIndex"),
		KeyConditions: map[string]*dynamodb.Condition{
			"userID": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(userID),
					},
				},
			},
		},
	}

	result, err := r.db.Query(input)
	if err != nil {
		return nil, err
	}

	var rooms []*model.Room
	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *RoomUserRepositoryImpl) AddUserToRoom(ctx context.Context, roomID, userID string) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.dbName),
		Item: map[string]*dynamodb.AttributeValue{
			"roomID": {
				S: aws.String(roomID),
			},
			"userID": {
				S: aws.String(userID),
			},
		},
	}

	_, err := r.db.PutItem(input)
	return err
}

func (r *RoomUserRepositoryImpl) RemoveUserFromRoom(ctx context.Context, roomID, userID string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.dbName),
		Key: map[string]*dynamodb.AttributeValue{
			"roomID": {
				S: aws.String(roomID),
			},
			"userID": {
				S: aws.String(userID),
			},
		},
	}

	_, err := r.db.DeleteItem(input)
	return err
}

func (r *RoomUserRepositoryImpl) AddUsersToRoom(ctx context.Context, roomID string, userIDs []string) error {
	writeRequests := make([]*dynamodb.TransactWriteItem, len(userIDs))
	for i, userID := range userIDs {
		writeRequests[i] = &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				TableName: aws.String(r.dbName),
				Item: map[string]*dynamodb.AttributeValue{
					"roomID": {
						S: aws.String(roomID),
					},
					"userID": {
						S: aws.String(userID),
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
