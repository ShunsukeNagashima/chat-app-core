package route

import (
	"github.com/gin-gonic/gin"
	"github.com/shunsukenagashima/chat-api/pkg/interface/controller"
)

func RegisterRoutes(router *gin.Engine, controllers *controller.Controllers) {
	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/hello", controllers.HelloController.SayHello)
		apiGroup.GET("/rooms/:roomId", controllers.RoomController.GetRoomByID)
		apiGroup.GET("/rooms", controllers.RoomController.GetAllPublicRoom)
		apiGroup.POST("/rooms", controllers.RoomController.CreateRoom)
		apiGroup.PUT("/rooms/:roomId", controllers.RoomController.UpdateRoom)
		apiGroup.DELETE("/rooms/:roomId", controllers.RoomController.DeleteRoom)
		apiGroup.GET("/users/:userId", controllers.UserController.GetUserByID)
		apiGroup.POST("/users", controllers.UserController.CreateUser)
		apiGroup.GET("/users/:userId/rooms", controllers.RoomUserController.GetAllRoomsByUserID)
		apiGroup.GET("/users/batch", controllers.UserController.BatchGetUsers)
		apiGroup.GET("/rooms/:roomId/users", controllers.RoomUserController.GetUsersByRoomID)
		apiGroup.DELETE("/rooms/:roomId/users/:userId", controllers.RoomUserController.RemoveUserFromRoom)
		apiGroup.POST("/rooms/:roomId/users", controllers.RoomUserController.AddUsersToRoom)
		apiGroup.GET("/rooms/:roomId/messages", controllers.MessageController.GetMessagesByRoomID)
		apiGroup.POST("/rooms/:roomId/messages", controllers.MessageController.CreateMessage)
		apiGroup.PUT("/rooms/:roomId/messages/:messageId", controllers.MessageController.UpdateMessage)
		apiGroup.DELETE("/rooms/:roomId/messages/:messageId", controllers.MessageController.DeleteMessage)
	}

	router.GET("/ws/:roomId", controllers.WSController.HandleRoomConnection)
	router.GET("/ws", controllers.WSController.HandleGlobalConnection)
}
