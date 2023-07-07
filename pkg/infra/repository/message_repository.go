package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/shunsukenagashima/chat-api/pkg/apperror"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/repository"
)

type MessageRepositoryImpl struct {
	db     *dynamodb.DynamoDB
	er     repository.ElasticsearchRepository
	dbName string
}

func NewMessageRepository(db *dynamodb.DynamoDB, er repository.ElasticsearchRepository) repository.MessageRepository {
	return &MessageRepositoryImpl{
		db,
		er,
		"Messages",
	}
}

func (mr *MessageRepositoryImpl) GetMessagesByRoomID(ctx context.Context, roomId, lastEvaluatedKey string, limit int) ([]*model.Message, string, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(mr.dbName),
		IndexName: aws.String("RoomIdIndex"),
		Limit:     aws.Int64(int64(limit)),
		ExclusiveStartKey: map[string]*dynamodb.AttributeValue{
			"messageId": {
				S: aws.String(lastEvaluatedKey),
			},
		},
		KeyConditions: map[string]*dynamodb.Condition{
			"roomId": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(roomId),
					},
				},
			},
		},
	}

	result, err := mr.db.Query(input)
	if err != nil {
		return nil, "", err
	}

	var messages []*model.Message
	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &messages); err != nil {
		return nil, "", err
	}

	var nextKey string
	if result.LastEvaluatedKey != nil {
		nextKey = *result.LastEvaluatedKey["messageId"].S
	}

	return messages, nextKey, nil
}

func (mr *MessageRepositoryImpl) GetByID(ctx context.Context, messageId string) (*model.Message, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(mr.dbName),
		Key: map[string]*dynamodb.AttributeValue{
			"messageId": {
				S: aws.String(messageId),
			},
		},
	}

	result, err := mr.db.GetItem(input)
	if err != nil {
		return nil, err
	}

	if len(result.Item) == 0 {
		return nil, apperror.NewNotFoundErr("Message", "Message'ID: "+messageId)
	}

	var message *model.Message
	if err := dynamodbattribute.UnmarshalMap(result.Item, &message); err != nil {
		return nil, err
	}

	return message, nil
}

func (mr *MessageRepositoryImpl) Create(ctx context.Context, message *model.Message) error {
	item, err := dynamodbattribute.MarshalMap(message)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(mr.dbName),
		Item:      item,
	}

	_, err = mr.db.PutItem(input)
	if err != nil {
		return err
	}

	messageJson, err := json.Marshal(message)
	if err != nil {
		return err
	}

	if err := mr.er.Create(ctx, "messages", message.MessageID, bytes.NewReader(messageJson)); err != nil {
		log.Printf(err.Error())
	}

	return nil
}

func (mr *MessageRepositoryImpl) Update(ctx context.Context, messageId, newContent string) error {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c": {
				S: aws.String(newContent),
			},
		},
		TableName: aws.String(mr.dbName),
		Key: map[string]*dynamodb.AttributeValue{
			"messageId": {
				S: aws.String(messageId),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set content = :c"),
	}

	_, err := mr.db.UpdateItem(input)
	if err != nil {
		return err
	}

	if err := mr.er.Update(ctx, "messages", messageId, strings.NewReader(fmt.Sprintf(`{"doc": {"content": "%s"}}`, newContent))); err != nil {
		log.Printf(err.Error())
	}

	return nil
}

func (mr *MessageRepositoryImpl) Delete(ctx context.Context, messageId string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(mr.dbName),
		Key: map[string]*dynamodb.AttributeValue{
			"messageId": {
				S: aws.String(messageId),
			},
		},
	}

	_, err := mr.db.DeleteItem(input)
	if err != nil {
		return err
	}

	if err := mr.er.Delete(ctx, "messages", messageId); err != nil {
		log.Printf(err.Error())
	}

	return nil
}
