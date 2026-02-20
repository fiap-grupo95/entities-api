package routes

import (
	handlers "mecanica_xpto/internal/adapter/http/handlers"

	"github.com/gin-gonic/gin"
)

func addPartsSupplyRoutes(rg *gin.RouterGroup, partsSupplyHandler *handlers.PartsSupplyHandler) {

	partsSupply := rg.Group(PathPartsSupply)
	{
		partsSupply.POST("/authorize-reserve", partsSupplyHandler.AuthorizeReserve)
		partsSupply.POST("/reserve", partsSupplyHandler.Reserve)
		partsSupply.POST("/release", partsSupplyHandler.Release)
		partsSupply.POST("/writeoff", partsSupplyHandler.WriteOff)
		partsSupply.GET("/:id", partsSupplyHandler.GetPartsSupplyByID)
		partsSupply.GET("/", partsSupplyHandler.ListPartsSupplies)
		partsSupply.POST("/", partsSupplyHandler.CreatePartsSupply)
		partsSupply.PUT("/:id", partsSupplyHandler.UpdatePartsSupply)
		partsSupply.DELETE("/:id", partsSupplyHandler.DeletePartsSupply)
	}
}
