package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shunsukenagashima/chat-api/pkg/interface/controllers"
)

func RegisterRoutes(router *gin.Engine) {

	hc := controllers.NewHelloController()

	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/hello", hc.SayHello)
	}
}
