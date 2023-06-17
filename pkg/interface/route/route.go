package route

import (
	"github.com/gin-gonic/gin"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/interface/controller"
)

func RegisterRoutes(router *gin.Engine) {
	hub := model.NewHub()
	go hub.Run()

	hc := controller.NewHelloController()
	wsc := controller.NewWSController(hub)

	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/hello", hc.SayHello)
	}

	router.GET("/ws", wsc.HandleConnection)
}
