package scripts

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func SetupMessages() error {
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
		LocalSecondaryIndexes: []*dynamodb.LocalSecondaryIndex{
			{
				IndexName: aws.String("MessageIdIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("roomId"),
						KeyType:       aws.String("HASH"),
					},
					{
						AttributeName: aws.String("messageId"),
						KeyType:       aws.String("RANGE"),
					},
				},
				Projection: &dynamodb.Projection{
					ProjectionType: aws.String("ALL"),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	// テストデータの投入
	_, err = svc.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"messageId": {
				S: aws.String("sampleMessageID"),
			},
			"content": {
				S: aws.String("sample message"),
			},
			"createdAt": {
				S: aws.String("2023-06-01T12:00:00Z"),
			},
			"roomId": {
				S: aws.String("sampleRoomID"),
			},
			"senderID": {
				S: aws.String("sampleSenderID"),
			},
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		return err
	}
	return nil
}
