package scripts

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"github.com/shunsukenagashima/chat-api/pkg/clock"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

func SetupMessages(users []*model.User, roomId string) error {
	tableName := "Messages"

	sess, _ := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	})

	svc := dynamodb.New(sess)

	// messages テーブルの作成
	_, err := svc.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("roomId"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("createdAt"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("roomId"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("createdAt"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		return err
	}

	clock := clock.FixedClocker{}

	for i, user := range users {
		// テストデータの投入
		_, err = svc.PutItem(&dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"messageId": {
					S: aws.String(uuid.New().String()),
				},
				"content": {
					S: aws.String("sample message" + strconv.Itoa(i)),
				},
				"createdAt": {
					S: aws.String(clock.Now().Add(time.Duration(i) * time.Minute).Format("2006-01-02T15:04:05Z")),
				},
				"roomId": {
					S: aws.String(roomId),
				},
				"userId": {
					S: aws.String(user.UserID),
				},
			},
			TableName: aws.String(tableName),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
