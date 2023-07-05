package scripts

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

func SetupUsers() ([]*model.User, error) {
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

	var users []*model.User

	firebaseUser := &model.User{
		UserID:   "Ko9BmAGyeBSP0w3WAnf83eg31rU2",
		Username: "Sample User",
		Email:    "sample-user@example.com",
	}

	users = append(users, firebaseUser)

	// テストデータの投入
	_, err = svc.PutItem(&dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"userId": {
				S: aws.String(firebaseUser.UserID),
			},
			"userName": {
				S: aws.String(firebaseUser.Username),
			},
			"email": {
				S: aws.String(firebaseUser.Email),
			},
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}

	for i := 1; i <= 10; i++ {
		userId := uuid.New().String()
		var userName string
		if i%2 == 0 {
			userName = "Sample User Even" + strconv.Itoa(i)
		} else {
			userName = "Sample User Odd" + strconv.Itoa(i)
		}
		email := "sample-user-" + strconv.Itoa(i) + "@example.com"
		// テストデータの投入
		_, err = svc.PutItem(&dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"userId": {
					S: aws.String(userId),
				},
				"userName": {
					S: aws.String(userName),
				},
				"email": {
					S: aws.String(email),
				},
			},
			TableName: aws.String(tableName),
		})
		if err != nil {
			return nil, err
		}
		user := &model.User{
			UserID:   userId,
			Username: userName,
			Email:    email,
		}
		users = append(users, user)
	}

	return users, nil
}
