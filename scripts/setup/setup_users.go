package scripts

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

func SetupUsers() (*model.User, error) {
	tableName := "Users"

	sess, _ := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	})

	svc := dynamodb.New(sess)

	// usersテーブルの作成
	_, err := svc.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("userId"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("userId"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}

	testUser := &model.User{
		UserID:   "Ko9BmAGyeBSP0w3WAnf83eg31rU2",
		Username: "Sample User",
		Email:    "sample-user@example.com",
	}

	// テストデータの投入
	_, err = svc.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"userId": {
				S: aws.String(testUser.UserID),
			},
			"userName": {
				S: aws.String(testUser.Username),
			},
			"email": {
				S: aws.String(testUser.Email),
			},
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}
	return testUser, nil
}
