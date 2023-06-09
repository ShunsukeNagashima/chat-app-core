package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/interface/controllers"
)

func RegisterRoutes(router *gin.Engine) {
	hub := model.NewHub()
	go hub.Run()

	hc := controllers.NewHelloController()
	wsc := controllers.NewWSController(hub)

	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/hello", hc.SayHello)
	}

	router.GET("/ws", wsc.HandleConnection)
}
