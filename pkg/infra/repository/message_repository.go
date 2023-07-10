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
		TableName:              aws.String(mr.dbName),
		Limit:                  aws.Int64(int64(limit)),
		KeyConditionExpression: aws.String("roomId = :r"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				S: aws.String(roomId),
			},
		},
	}

	if lastEvaluatedKey != "" {
		input.ExclusiveStartKey = map[string]*dynamodb.AttributeValue{
			"roomId": {
				S: aws.String(roomId),
			},
			"createdAt": {
				S: aws.String(lastEvaluatedKey),
			},
		}
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
		nextKey = *result.LastEvaluatedKey["createdAt"].S
	}

	return messages, nextKey, nil
}

func (mr *MessageRepositoryImpl) GetByID(ctx context.Context, roomId, messageId string) (*model.Message, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(mr.dbName),
		IndexName:              aws.String("MessageIdIndex"),
		KeyConditionExpression: aws.String("roomId = :r and messageId = :m"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				S: aws.String(roomId),
			},
			":m": {
				S: aws.String(messageId),
			},
		},
	}

	result, err := mr.db.Query(input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, apperror.NewNotFoundErr("Message", "Message'ID: "+messageId)
	}

	var message *model.Message
	if err := dynamodbattribute.UnmarshalMap(result.Items[0], &message); err != nil {
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
		log.Printf("%s", err.Error())
	}

	return nil
}

func (mr *MessageRepositoryImpl) Update(ctx context.Context, roomId, messageId, newContent string) error {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(mr.dbName),
		IndexName:              aws.String("MessageIdIndex"),
		KeyConditionExpression: aws.String("roomId = :r and messageId = :m"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				S: aws.String(roomId),
			},
			":m": {
				S: aws.String(messageId),
			},
		},
	}

	result, err := mr.db.Query(queryInput)
	if err != nil {
		return err
	}

	updateInput := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c": {
				S: aws.String(newContent),
			},
		},
		TableName:        aws.String(mr.dbName),
		Key:              result.Items[0],
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set content = :c"),
	}

	_, err = mr.db.UpdateItem(updateInput)
	if err != nil {
		return err
	}

	if err := mr.er.Update(ctx, "messages", messageId, strings.NewReader(fmt.Sprintf(`{"doc": {"content": "%s"}}`, newContent))); err != nil {
		log.Printf("%s", err.Error())
	}

	return nil
}

func (mr *MessageRepositoryImpl) Delete(ctx context.Context, roomId, messageId string) error {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(mr.dbName),
		IndexName:              aws.String("MessageIdIndex"),
		KeyConditionExpression: aws.String("roomId = :r and messageId = :m"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				S: aws.String(roomId),
			},
			":m": {
				S: aws.String(messageId),
			},
		},
	}

	result, err := mr.db.Query(queryInput)
	if err != nil {
		return err
	}

	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(mr.dbName),
		Key:       result.Items[0],
	}

	_, err = mr.db.DeleteItem(deleteInput)
	if err != nil {
		return err
	}

	if err := mr.er.Delete(ctx, "messages", messageId); err != nil {
		log.Printf("%s", err.Error())
	}

	return nil
}
