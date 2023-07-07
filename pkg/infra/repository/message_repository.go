package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/repository"
)

type MessageRepositoryImpl struct {
	db     *dynamodb.DynamoDB
	es     *elasticsearch.Client
	dbName string
}

func NewMessageRepository(db *dynamodb.DynamoDB, es *elasticsearch.Client) repository.MessageRepository {
	return &MessageRepositoryImpl{
		db,
		es,
		"Messages",
	}
}

func (r *MessageRepositoryImpl) GetAllMessagesByRoomID(ctx context.Context, roomId string) ([]*model.Message, error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(r.dbName),
		IndexName: aws.String("RoomIdIndex"),
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

	result, err := r.db.Query(input)
	if err != nil {
		return nil, err
	}

	var messages []*model.Message
	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *MessageRepositoryImpl) Create(ctx context.Context, message *model.Message) error {
	item, err := dynamodbattribute.MarshalMap(message)
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

	messageJson, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "messages",
		DocumentID: message.MessageID,
		Body:       bytes.NewReader(messageJson),
		Refresh:    "true",
	}

	retryCount := 3
	var lastErr error
	for i := 0; i < retryCount; i++ {
		res, err := req.Do(ctx, r.es)
		if err != nil {
			log.Printf("Error getting response: %s", err)
			lastErr = fmt.Errorf("Error getting response: %s", err)
			continue
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusCreated {
			log.Printf("Unexpected status code returned %d", res.StatusCode)
			lastErr = fmt.Errorf("unexpected status code: %d", res.StatusCode)
			continue
		}

		lastErr = nil
		break
	}

	if lastErr != nil {
		return fmt.Errorf("failed to create index for message: %s", lastErr)
	}

	return nil
}

func (r *MessageRepositoryImpl) Update(ctx context.Context, message *model.Message) error {
	item, err := dynamodbattribute.MarshalMap(message)
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

	messageJson, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req := esapi.UpdateRequest{
		Index:      "messages",
		DocumentID: message.MessageID,
		Body:       bytes.NewReader(messageJson),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.es)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}

func (r *MessageRepositoryImpl) Delete(ctx context.Context, messageId string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.dbName),
		Key: map[string]*dynamodb.AttributeValue{
			"messageId": {
				S: aws.String(messageId),
			},
		},
	}

	_, err := r.db.DeleteItem(input)
	if err != nil {
		return err
	}

	req := esapi.DeleteRequest{
		Index:      "messages",
		DocumentID: messageId,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.es)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}
