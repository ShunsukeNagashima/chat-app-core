package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/usecase"
)

type RoomController struct {
	roomUsecase usecase.RoomUsecase
	validator   *validator.Validate
}

func NewRoomController(roomUsecase usecase.RoomUsecase, validator *validator.Validate) *RoomController {
	return &RoomController{
		roomUsecase,
		validator,
	}
}

func (rc *RoomController) GetRoomByID(ctx *gin.Context) {
	roomId := ctx.Param("roomId")

	result, err := rc.roomUsecase.GetRoomByID(ctx.Request.Context(), roomId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"result": result})
}

func (rc *RoomController) GetAllPublicRoom(ctx *gin.Context) {
	result, err := rc.roomUsecase.GetAllPublicRoom(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"result": result})
}

func (rc *RoomController) CreateRoom(ctx *gin.Context) {
	var req struct {
		Name     string `json:"name" validate:"required,min=1,max=30,alnumdash"`
		RoomType string `json:"roomType" validate:"required,oneof=public private"`
		OwnerID  string `json:"ownerID" validate:"required"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := rc.validator.Struct(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomType, err := model.ParseRoomType(req.RoomType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room := &model.Room{
		RoomID:   uuid.New().String(),
		Name:     req.Name,
		RoomType: roomType,
	}

	if err := rc.roomUsecase.CreateRoom(ctx, room, req.OwnerID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"result": room})
}

func (rc *RoomController) DeleteRoom(ctx *gin.Context) {
	roomId := ctx.Param("roomId")

	if err := rc.roomUsecase.DeleteRoom(ctx.Request.Context(), roomId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": "room deleted successfully"})
}

func (rc *RoomController) UpdateRoom(ctx *gin.Context) {
	roomId := ctx.Param("roomId")

	var req struct {
		Name     string `json:"name" validate:"required,min=1,max=30,alnumdash"`
		RoomType string `json:"roomType" validate:"required,oneof=public private"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := rc.validator.Struct(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomType, err := model.ParseRoomType(req.RoomType)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	room := &model.Room{
		RoomID:   roomId,
		Name:     req.Name,
		RoomType: roomType,
	}

	if err := rc.roomUsecase.UpdateRoom(ctx, room); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"result": "room updated successfully"})
}
