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
			{
				AttributeName: aws.String("createdAt"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("userId"),
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
		return nil, err
	}

	var users []*model.User
	clock := clock.FixedClocker{}
	firebaseUser := &model.User{
		UserID:    "qTF9aUAHNqNyi7R3sQtSGSRhTft1",
		Username:  "Sample User",
		Email:     "shun.mmks_n@icloud.com",
		ImageURL:  "https://images.unsplash.com/photo-1552053831-71594a27632d?fit=crop&w=500&h=500",
		CreatedAt: clock.Now(),
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
			"imageUrl": {
				S: aws.String(firebaseUser.ImageURL),
			},
			"createdAt": {
				S: aws.String(firebaseUser.CreatedAt.Format("2006-01-02T15:04:05Z")),
			},
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}

	for i := 1; i <= 30; i++ {
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
				"createdAt": {
					S: aws.String(clock.Now().Add(time.Duration(i) * time.Hour).Format("2006-01-02T15:04:05Z")),
				},
				"imageUrl": {
					S: aws.String("test-image-url"),
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
