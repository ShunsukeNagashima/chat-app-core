package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"regexp"

	firebase "firebase.google.com/go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/elastic/go-elasticsearch/v8"
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
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}

	router.Use(cors.New(corsConfig))

	route.RegisterRoutes(router, controllers)

	return router.Run()
}

func initializeControllers(ctx context.Context) (*controller.Controllers, error) {
	hm := model.NewRoomHubManager()

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://dynamodb-local:8000"),
	})
	if err != nil {
		return nil, err
	}
	db := dynamodb.New(sess)

	credentials, err := os.ReadFile("/run/secrets/firebase-credentials")
	if err != nil {
		return nil, err
	}
	opt := option.WithCredentialsJSON(credentials)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"https://elasticsearch:9200",
		},
		Username: "elastic",
		Password: "password",
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	fa := auth.NewFirebaseAuth(client)

	rr := repository.NewRoomRepository(db)
	rur := repository.NewRoomUserRepository(db)
	ur := repository.NewUserRepository(db, es)

	ru := usecase.NewRoomUsecase(rr, ur)
	ruu := usecase.NewRoomUserUsecase(rur, ur, rr)
	uu := usecase.NewUserUsecase(ur, fa)

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
	}

	return controllers, nil
}

func isAlnumOrDash(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(fl.Field().String())
}
