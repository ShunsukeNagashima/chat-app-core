package repository

import (
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

func (r *RoomUserRepositoryImpl) GetAllByUserID(userID string) ([]*model.Room, error) {
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
