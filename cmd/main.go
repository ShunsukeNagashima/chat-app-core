package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	firebase "firebase.google.com/go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/infra/auth"
	"github.com/shunsukenagashima/chat-api/pkg/infra/repository"
	"github.com/shunsukenagashima/chat-api/pkg/interface/controller"
	"github.com/shunsukenagashima/chat-api/pkg/interface/route"
	"github.com/shunsukenagashima/chat-api/pkg/usecase"
	"google.golang.org/api/option"
)

type AppSecret struct {
	FirebaseCredentials string `json:"firebase_credentials"`
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalf("Server failed to run with %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	router := gin.Default()

	controllers, err := initializeControllers(ctx)
	if err != nil {
		return err
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	router.Use(cors.New(corsConfig))

	route.RegisterRoutes(router, controllers)

	return router.Run()
}

func initializeControllers(ctx context.Context) (*controller.Controllers, error) {
	hm := model.NewRoomHubManager()

	db, err := initializeDynamodbClient()
	if err != nil {
		return nil, err
	}

	svc, err := initializeSecretManagerClient()
	if err != nil {
		return nil, err
	}

	env := os.Getenv("APP_ENV")
	secretName := fmt.Sprintf("app/%s/AppSecrets", env)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return nil, err
	}

	log.Println("result.SecretString: ", *result.SecretString)
	opt := option.WithCredentialsJSON([]byte(*result.SecretString))
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	fa := auth.NewFirebaseAuth(client)

	rr := repository.NewRoomRepository(db)
	rur := repository.NewRoomUserRepository(db)
	ur := repository.NewUserRepository(db)
	mr := repository.NewMessageRepository(db)

	ru := usecase.NewRoomUsecase(rr, ur)
	ruu := usecase.NewRoomUserUsecase(rur, ur, rr)
	uu := usecase.NewUserUsecase(ur, fa)
	mu := usecase.NewMessageUsecase(mr)

	v := validator.New()

	if err := v.RegisterValidation("alnumdash", isAlnumOrDash); err != nil {
		return nil, err
	}

	controllers := &controller.Controllers{
		HelloController:    controller.NewHelloController(),
		WSController:       controller.NewWSController(hm),
		RoomController:     controller.NewRoomController(ru, v),
		RoomUserController: controller.NewRoomUserController(ruu, v),
		UserController:     controller.NewUserController(uu, v),
		MessageController:  controller.NewMessageController(mu, v),
	}

	return controllers, nil
}

func initializeDynamodbClient() (*dynamodb.DynamoDB, error) {
	if os.Getenv("APP_ENV") == "local" {
		sess, err := session.NewSession(&aws.Config{
			Region:   aws.String("ap-northeast-1"),
			Endpoint: aws.String("http://dynamodb-local:8000"),
		})
		if err != nil {
			return nil, err
		}
		return dynamodb.New(sess), nil
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	})
	if err != nil {
		return nil, err
	}

	return dynamodb.New(sess), nil
}

func initializeSecretManagerClient() (*secretsmanager.SecretsManager, error) {
	if os.Getenv("APP_ENV") == "local" {
		log.Println("local mode")
		sess, err := session.NewSession(&aws.Config{
			Region:   aws.String("ap-northeast-1"),
			Endpoint: aws.String("http://localstack:4566"),
		})
		if err != nil {
			return nil, err
		}
		return secretsmanager.New(sess), nil
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	})
	if err != nil {
		return nil, err
	}
	return secretsmanager.New(sess), nil
}

func isAlnumOrDash(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(fl.Field().String())
}
