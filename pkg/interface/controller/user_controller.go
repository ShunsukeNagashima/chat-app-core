package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/usecase"
)

type UserController struct {
	userUsecase usecase.UserUsecase
	validator   *validator.Validate
}

func NewUserController(userUsecase usecase.UserUsecase, validator *validator.Validate) *UserController {
	return &UserController{
		userUsecase,
		validator,
	}
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var req struct {
		Name  string `json:"name" validate:"required,min=1,max=30"`
		Email string `json:"email" validate:"required,email"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.validator.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &model.User{
		UserID:   uuid.New().String(),
		Username: req.Name,
		Email:    req.Email,
	}

	if err := uc.userUsecase.CreateUser(ctx.Request.Context(), user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": user})
}

func (uc *UserController) GetUserByID(ctx *gin.Context) {
	userID := ctx.Param("userID")

	result, err := uc.userUsecase.GetUserByID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": result})
}
