package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		UserID  string `json:"userId" validate:"required"`
		Name    string `json:"name" validate:"required,min=1,max=30"`
		Email   string `json:"email" validate:"required,email"`
		IDToken string `json:"idToken" validate:"required"`
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
		UserID:   req.UserID,
		Username: req.Name,
		Email:    req.Email,
	}

	if err := uc.userUsecase.CreateUser(ctx.Request.Context(), user, req.IDToken); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": user})
}

func (uc *UserController) GetUserByID(ctx *gin.Context) {
	userId := ctx.Param("userId")

	result, err := uc.userUsecase.GetUserByID(ctx.Request.Context(), userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": result})
}

func (uc *UserController) SearchUsers(ctx *gin.Context) {
	query := ctx.Query("query")
	fromStr := ctx.DefaultQuery("from", "0")
	sizeStr := ctx.DefaultQuery("size", "20")

	from, err := strconv.Atoi(fromStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := uc.userUsecase.SearchUsers(ctx.Request.Context(), query, from, size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"result": users})
}
