package routes

import (
	handlers "mecanica_xpto/internal/adapter/http/handlers"

	"github.com/gin-gonic/gin"
)

func addUserRoutes(rg *gin.RouterGroup, userHandler *handlers.UserHandler) {
	users := rg.Group(PathUsers)
	{
		users.GET("/", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUser)
		users.POST("/", userHandler.CreateUser)
		users.PATCH("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}
}
