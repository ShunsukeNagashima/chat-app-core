package route

import (
	"github.com/gin-gonic/gin"
	"github.com/shunsukenagashima/chat-api/pkg/interface/controller"
)

func RegisterRoutes(router *gin.Engine, controllers *controller.Controllers) {
	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/hello", controllers.HelloController.SayHello)
		apiGroup.GET("/rooms/:roomID", controllers.RoomController.GetRoomByID)
		apiGroup.GET("/rooms", controllers.RoomController.GetAllPublicRoom)
		apiGroup.POST("/rooms", controllers.RoomController.CreateRoom)
		apiGroup.PUT("/rooms/:roomID", controllers.RoomController.UpdateRoom)
		apiGroup.DELETE("/rooms/:roomID", controllers.RoomController.DeleteRoom)
		apiGroup.GET("/users/:userID", controllers.UserController.GetUserByID)
		apiGroup.POST("/users", controllers.UserController.CreateUser)
		apiGroup.GET("/users/:userID/rooms", controllers.RoomUserController.GetAllRoomsByUserID)
		apiGroup.DELETE("/rooms/:roomID/users/:userID", controllers.RoomUserController.RemoveUserFromRoom)
		apiGroup.POST("/rooms/:roomID/users", controllers.RoomUserController.AddUsersToRoom)
	}

	router.GET("/ws", controllers.WSController.HandleConnection)
}
