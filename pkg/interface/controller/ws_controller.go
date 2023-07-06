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
	HubManager *model.RoomHubManager
}

func NewWSController(hubManager *model.RoomHubManager) *WSController {
	return &WSController{
		HubManager: hubManager,
	}
}

func (wc *WSController) HandleRoomConnection(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Failed to set webscoket upgrade: %+v", err)
		return
	}

	roomId := ctx.Param("roomId")
	if roomId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "roomId is required"})
		return
	}

	hub, exists := wc.HubManager.GetRoomHub(roomId)
	if !exists {
		hub = wc.HubManager.CreateRoomHub(roomId)
	}
	client := model.NewClient(conn, hub)

	hub.RegisterClient(client)

	go client.Write()
	go client.Read()
}

func (wc *WSController) HandleGlobalConnection(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Failed to set webscoket upgrade: %+v", err)
		return
	}

	globalHub := model.GetGlobalHubInstance()

	client := model.NewClient(conn, globalHub)

	globalHub.RegisterClient(client)

	go client.Write()
	go client.Read()
}
