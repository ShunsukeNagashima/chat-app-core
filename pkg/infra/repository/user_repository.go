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
	"github.com/shunsukenagashima/chat-api/pkg/apperror"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/repository"
)

type UserRepositoryImpl struct {
	db     *dynamodb.DynamoDB
	es     *elasticsearch.Client
	dbName string
}

type EsSearchResponse struct {
	Hits struct {
		Hits []struct {
			Source *model.User `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func NewUserRepository(db *dynamodb.DynamoDB, es *elasticsearch.Client) repository.UserRepository {
	return &UserRepositoryImpl{
		db,
		es,
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

	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "users",
		DocumentID: user.UserID,
		Body:       bytes.NewReader(userJson),
		Refresh:    "true",
	}

	retryCount := 3
	var lastErr error

	for i := 0; i < retryCount; i++ {
		res, err := req.Do(ctx, r.es)
		if err != nil {
			log.Printf("Error when trying to create index for user: %v", err)
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
		log.Printf("Failed to create index for user after %d attempts", retryCount)
	}

	return nil
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

func (r *UserRepositoryImpl) SearchUsers(ctx context.Context, query string, from, size int) ([]*model.User, error) {
	var buf bytes.Buffer
	queryMap := map[string]interface{}{
		"from": from,
		"size": size,
		"query": map[string]interface{}{
			"match_phrase_prefix": map[string]interface{}{
				"userName": query,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(queryMap); err != nil {
		return nil, err
	}

	res, err := r.es.Search(
		r.es.Search.WithContext(ctx),
		r.es.Search.WithIndex("users"),
		r.es.Search.WithBody(&buf),
		r.es.Search.WithTrackTotalHits(true),
		r.es.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error getting response: %s", res.String())
	}

	var esResponse EsSearchResponse
	if err := json.NewDecoder(res.Body).Decode(&esResponse); err != nil {
		return nil, err
	}

	users := make([]*model.User, len(esResponse.Hits.Hits))
	for i, hit := range esResponse.Hits.Hits {
		users[i] = hit.Source
	}

	return users, nil
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
