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
	userID := ctx.Param("userID")

	rooms, err := rc.roomUserUsecase.GetAllRoomsByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": rooms})
}

func (rc *RoomUserController) RemoveUserFromRoom(ctx *gin.Context) {
	roomID := ctx.Param("roomID")
	userID := ctx.Param("userID")

	if err := rc.roomUserUsecase.RemoveUserFromRoom(ctx, roomID, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success to remove the user from the room"})
}

func (rc *RoomUserController) AddUsersToRoom(ctx *gin.Context) {
	var req struct {
		RoomID  string   `json:"roomID" validate:"required"`
		UserIDs []string `json:"userIDs" validate:"required"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := rc.roomUserUsecase.AddUsersToRoom(ctx, req.RoomID, req.UserIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "success to add the users to the room"})
}
