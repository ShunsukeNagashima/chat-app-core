package scripts

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CleanUpDynamodb() error {
	sess, _ := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	})

	svc := dynamodb.New(sess)

	tableNames := []string{"Users", "Messages", "RoomUsers", "Rooms", "Likes", "Readby"}

	for _, tableName := range tableNames {
		_, err := svc.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(tableName),
		})

		if err != nil {
			return err
		}
	}

	log.Print("Dynamodb tables cleanup completed")

	return nil
}
