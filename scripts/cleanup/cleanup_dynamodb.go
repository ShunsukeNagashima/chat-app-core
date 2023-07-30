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

	result, err := svc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return err
	}

	for _, tableName := range result.TableNames {
		_, err := svc.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: tableName,
		})

		if err != nil {
			return err
		}

		log.Printf("Table %s deleted", *tableName)
	}

	log.Print("Dynamodb tables cleanup completed")

	return nil
}
