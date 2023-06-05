package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shunsukenagashima/chat-api/pkg/interface/controllers"
	"github.com/stretchr/testify/assert"
)

func TestHelloController(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	hc := controllers.NewHelloController()
	hc.SayHello(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message":"Hello World"}`, w.Body.String())
}
