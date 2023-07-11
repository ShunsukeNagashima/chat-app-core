package scripts

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

func SetupRoomUsers(roomIDs []string, users []*model.User) error {
	tableName := "RoomUsers"

	sess, _ := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	})

	svc := dynamodb.New(sess)

	// likes テーブルの作成
	_, err := svc.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("roomId"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("userId"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("roomId"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("userId"),
				KeyType:       aws.String("RANGE"),
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("UserIDIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("userId"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("roomId"),
						KeyType:       aws.String("RANGE"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(10),
					WriteCapacityUnits: aws.Int64(10),
				},
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

	// テストデータの投入
	for _, roomId := range roomIDs {
		_, err = svc.PutItem(&dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"roomId": {
					S: aws.String(roomId),
				},
				"userId": {
					S: aws.String(users[0].UserID),
				},
			},
			TableName: aws.String(tableName),
		})
		if err != nil {
			return err
		}
	}

	for _, user := range users {
		_, err = svc.PutItem(&dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"roomId": {
					S: aws.String(roomIDs[0]),
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
