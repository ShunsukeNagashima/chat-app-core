package scripts

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func SetupMessages() error {
	sess, _ := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	})

	svc := dynamodb.New(sess)

	// messages テーブルの作成
	_, err := svc.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("roomID"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("createdAt"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("roomID"),
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
		TableName: aws.String("messages"),
	})
	if err != nil {
		return err
	}

	// テストデータの投入
	_, err = svc.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"messageID": {
				S: aws.String("sampleMessageID"),
			},
			"content": {
				S: aws.String("sample message"),
			},
			"createdAt": {
				S: aws.String("2023-06-01T12:00:00Z"),
			},
			"roomID": {
				S: aws.String("sampleRoomID"),
			},
			"senderID": {
				S: aws.String("sampleSenderID"),
			},
		},
		TableName: aws.String("messages"),
	})
	if err != nil {
		return err
	}
	return nil
}
