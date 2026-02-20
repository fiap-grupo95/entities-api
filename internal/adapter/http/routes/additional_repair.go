//go:build legacy
// +build legacy

package routes

import (
	handlers "mecanica_xpto/internal/adapter/http/handlers"

	"github.com/gin-gonic/gin"
)

func addAdditionalRepairRoutes(rg *gin.RouterGroup, additionalRepair *handlers.AdditionalRepairHandler) {
	serviceOrdersRoutes := rg.Group(PathAdditionalRepair)
	{
		serviceOrdersRoutes.POST("", additionalRepair.CreateAdditionalRepair)
		serviceOrdersRoutes.GET("/:id", additionalRepair.GetAdditionalRepair)
		serviceOrdersRoutes.PATCH("/:id/add", additionalRepair.AddPartSupplyAndService)
		serviceOrdersRoutes.PATCH("/:id/remove", additionalRepair.RemovePartSupplyAndService)
		serviceOrdersRoutes.PATCH("/:id/customer_approval", additionalRepair.CustomerApproval)
	}
}
