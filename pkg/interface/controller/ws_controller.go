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
	HubManager *model.HubManager
}

func NewWSController(hubManager *model.HubManager) *WSController {
	return &WSController{
		HubManager: hubManager,
	}
}

func (wc *WSController) HandleConnection(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set webscoket upgrade: %+v", err)
		return
	}

	roomID := c.Param("roomID")

	hub := wc.HubManager.GetHub(roomID)
	client := model.NewClient(conn, hub)

	hub.Register <- client

	go client.Write()
	go client.Read()
}
