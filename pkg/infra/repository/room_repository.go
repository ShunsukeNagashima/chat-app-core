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

type RoomRepository struct {
	db     *dynamodb.DynamoDB
	dbName string
}

func NewRoomRepository(db *dynamodb.DynamoDB) repository.RoomRepository {
	return &RoomRepository{
		db,
		"Rooms",
	}
}

func (r *RoomRepository) GetById(ctx context.Context, roomID string) (*model.Room, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.dbName),
		Key: map[string]*dynamodb.AttributeValue{
			"roomID": {
				S: aws.String(roomID),
			},
		},
	}

	result, err := r.db.GetItem(input)

	if err != nil {
		return nil, err
	}

	if len(result.Item) == 0 {
		return nil, fmt.Errorf("room with ID %s not found", roomID)
	}

	var room model.Room
	if err := dynamodbattribute.UnmarshalMap(result.Item, &room); err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *RoomRepository) GetByName(ctx context.Context, name string) (*model.Room, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String("rooms"),
		IndexName: aws.String("NameIndex"),
		KeyConditions: map[string]*dynamodb.Condition{
			"name": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(name),
					},
				},
			},
		},
	}

	result, err := r.db.QueryWithContext(ctx, input)

	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	room := new(model.Room)
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &room)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r *RoomRepository) GetAllPublic(ctx context.Context) ([]*model.Room, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String(r.dbName),
		FilterExpression: aws.String("room_type = :public"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":public": {
				S: aws.String("Public"),
			},
		},
	}

	result, err := r.db.Scan(input)
	if err != nil {
		return nil, err
	}

	var rooms []*model.Room
	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *RoomRepository) Create(ctx context.Context, room *model.Room) error {
	item, err := dynamodbattribute.MarshalMap(room)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Rooms"),
		Item:      item,
	}

	_, err = r.db.PutItem(input)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoomRepository) Delete(ctx context.Context, roomID string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("Rooms"),
		Key: map[string]*dynamodb.AttributeValue{
			"roomID": {
				S: aws.String(roomID),
			},
		},
	}

	_, err := r.db.DeleteItem(input)
	if err != nil {
		return err
	}
	return nil
}

func (r *RoomRepository) Update(ctx context.Context, room *model.Room) error {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#N": aws.String("name"),
			"#T": aws.String("room_type"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":n": {
				S: aws.String(room.Name),
			},
			":t": {
				S: aws.String(string(room.RoomType)),
			},
		},
		TableName: aws.String(r.dbName),
		Key: map[string]*dynamodb.AttributeValue{
			"roomID": {
				S: aws.String(room.RoomID),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("SET #N = :n, #T = :t"),
	}

	_, err := r.db.UpdateItem(input)
	if err != nil {
		return err
	}

	return nil
}
