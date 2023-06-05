package controllers

import "github.com/gin-gonic/gin"

type HelloController struct{}

func NewHelloController() *HelloController {
	return &HelloController{}
}

func (hc *HelloController) SayHello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}
