package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shunsukenagashima/chat-api/pkg/domain/usecase"
)

type RoomUserController struct {
	roomUserUsecase usecase.RoomUserUsecase
	validator       *validator.Validate
}

func NewRoomUserController(roomUserUsecase usecase.RoomUserUsecase, validator *validator.Validate) *RoomUserController {
	return &RoomUserController{
		roomUserUsecase,
		validator,
	}
}

func (rc *RoomUserController) GetAllRoomsByUserID(ctx *gin.Context) {
	userId := ctx.Param("userId")

	rooms, err := rc.roomUserUsecase.GetAllRoomsByUserID(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": rooms})
}

func (rc *RoomUserController) RemoveUserFromRoom(ctx *gin.Context) {
	roomId := ctx.Param("roomId")
	userId := ctx.Param("userId")

	if err := rc.roomUserUsecase.RemoveUserFromRoom(ctx, roomId, userId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": "success to remove the user from the room"})
}

func (rc *RoomUserController) AddUsersToRoom(ctx *gin.Context) {
	roomId := ctx.Param("roomId")

	var req struct {
		UserIDs []string `json:"userIds" validate:"required"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := rc.roomUserUsecase.AddUsersToRoom(ctx, roomId, req.UserIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"result": "success to add the users to the room"})
}
