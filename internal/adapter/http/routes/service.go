package routes

import (
	handlers "mecanica_xpto/internal/adapter/http/handlers"

	"github.com/gin-gonic/gin"
)

func addServiceRoutes(rg *gin.RouterGroup, serviceHandler *handlers.ServiceHandler) {

	service := rg.Group(PathService)
	{
		service.GET("/:id", serviceHandler.GetServiceByID)
		service.GET("/", serviceHandler.ListServices)
		service.POST("/", serviceHandler.CreateService)
		service.PUT("/:id", serviceHandler.UpdateService)
		service.DELETE("/:id", serviceHandler.DeleteService)
	}
}
