package main

import (
	"context"
	"log"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/infra/repository"
	"github.com/shunsukenagashima/chat-api/pkg/interface/controller"
	"github.com/shunsukenagashima/chat-api/pkg/interface/route"
	"github.com/shunsukenagashima/chat-api/pkg/usecase"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalf("Server failed to run with %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	router := gin.Default()

	controllers, err := initializeControllers()
	if err != nil {
		return err
	}

	route.RegisterRoutes(router, controllers)

	return router.Run()
}

func initializeControllers() (*controller.Controllers, error) {
	hm := model.NewHubManager()

	sess, _ := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	})
	db := dynamodb.New(sess)

	rr := repository.NewRoomRepository(db)
	ru := usecase.NewRoomUsecase(rr)
	v := validator.New()

	v.RegisterValidation("alnumdash", isAlnumOrDash)

	controllers := &controller.Controllers{
		HelloController: controller.NewHelloController(),
		WSController:    controller.NewWSController(hm),
		RoomController:  controller.NewRoomController(ru, v),
	}

	return controllers, nil
}

func isAlnumOrDash(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(fl.Field().String())
}
