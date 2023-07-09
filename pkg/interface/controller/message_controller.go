package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/usecase"
)

type MessageController struct {
	messageUsecase usecase.MessageUsecase
	validator      *validator.Validate
}

func NewMessageController(messageUsecase usecase.MessageUsecase, validator *validator.Validate) *MessageController {
	return &MessageController{
		messageUsecase,
		validator,
	}
}

func (mc *MessageController) GetMessagesByRoomID(ctx *gin.Context) {
	roomId := ctx.Param("roomId")
	lastEvaluatedKey := ctx.Query("lastEvaluatedKey")
	limit := ctx.DefaultQuery("limit", "10")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	result, nextKey, err := mc.messageUsecase.GetMessagesByRoomID(ctx.Request.Context(), roomId, lastEvaluatedKey, limitInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result":  result,
		"nextKey": nextKey,
	})
}

func (mc *MessageController) CreateMessage(ctx *gin.Context) {
	roomId := ctx.Param("roomId")

	var req struct {
		UserID  string `json:"userId" validate:"required"`
		Content string `json:"content" validate:"required,min=1"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := mc.validator.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := &model.Message{
		RoomID:  roomId,
		UserID:  req.UserID,
		Content: req.Content,
	}

	if err := mc.messageUsecase.CreateMessage(ctx.Request.Context(), message); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"result": message})
}

func (mc *MessageController) UpdateMessage(ctx *gin.Context) {
	roomId := ctx.Param("roomId")
	messageId := ctx.Param("messageId")

	var req struct {
		Content string `json:"content" validate:"required,min=1"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := mc.validator.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := mc.messageUsecase.UpdateMessage(ctx.Request.Context(), roomId, messageId, req.Content); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": "message updated successfully"})
}

func (mc *MessageController) DeleteMessage(ctx *gin.Context) {
	roomId := ctx.Param("roomId")
	messageId := ctx.Param("messageId")

	if err := mc.messageUsecase.DeleteMessage(ctx.Request.Context(), roomId, messageId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": "message deleted successfully"})
}
