package route

import (
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/infra/repository"
	"github.com/shunsukenagashima/chat-api/pkg/interface/controller"
	"github.com/shunsukenagashima/chat-api/pkg/usecase"
)

func RegisterRoutes(router *gin.Engine) {
	hub := model.NewHub()
	go hub.Run()

	sess, _ := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	})

	db := dynamodb.New(sess)

	rr := repository.NewRoomRepository(db)
	ru := usecase.NewRoomUsecase(rr)
	v := validator.New()

	v.RegisterValidation("alnumdash", isAlnumOrDash)

	hc := controller.NewHelloController()
	wsc := controller.NewWSController(hub)
	rc := controller.NewRoomController(ru, v)

	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/hello", hc.SayHello)
		apiGroup.GET("/rooms/:roomID", rc.GetRoomByID)
		apiGroup.GET("/rooms", rc.GetAllPublicRoom)
		apiGroup.POST("/rooms", rc.CreateRoom)
		apiGroup.PUT("/rooms/:roomID", rc.UpdateRoom)
		apiGroup.DELETE("/rooms/:roomID", rc.DeleteRoom)
	}

	router.GET("/ws", wsc.HandleConnection)
}

func isAlnumOrDash(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(fl.Field().String())
}
