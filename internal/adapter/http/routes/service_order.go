//go:build legacy
// +build legacy

package routes

import (
	handlers "mecanica_xpto/internal/adapter/http/handlers"

	"github.com/gin-gonic/gin"
)

func addServiceOrderRoutes(rg *gin.RouterGroup, serviceOrderHandler *handlers.ServiceOrderHandler) {
	serviceOrdersRoutes := rg.Group(PathServiceOrders)
	{
		serviceOrdersRoutes.GET("/:id", serviceOrderHandler.GetServiceOrder)
		serviceOrdersRoutes.POST("", serviceOrderHandler.CreateServiceOrder)
		serviceOrdersRoutes.PATCH("/:id/diagnosis", serviceOrderHandler.UpdateServiceOrderDiagnosis)
		serviceOrdersRoutes.PATCH("/:id/estimate", serviceOrderHandler.UpdateServiceOrderEstimate)
		serviceOrdersRoutes.PATCH("/:id/execution", serviceOrderHandler.UpdateServiceOrderExecution)
		serviceOrdersRoutes.PATCH("/:id/delivery", serviceOrderHandler.UpdateServiceOrderDelivery)
		serviceOrdersRoutes.GET("/", serviceOrderHandler.ListServiceOrders)
	}
}
