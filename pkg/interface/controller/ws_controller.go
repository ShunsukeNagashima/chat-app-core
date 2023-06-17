package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSController struct {
	Hub *model.Hub
}

func NewWSController(hub *model.Hub) *WSController {
	return &WSController{
		Hub: hub,
	}
}

func (wc *WSController) HandleConnection(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set webscoket upgrade: %+v", err)
		return
	}

	client := model.NewClient(conn, wc.Hub)

	wc.Hub.Register <- client

	go client.Write()
	go client.Read()
}
