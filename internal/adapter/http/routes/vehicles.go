package routes

import (
	handlers "mecanica_xpto/internal/adapter/http/handlers"

	"github.com/gin-gonic/gin"
)

func addVehicleRoutes(rg *gin.RouterGroup, vehicleHandler *handlers.VehicleHandler) {
	vehicles := rg.Group(PathVehicles)
	{
		vehicles.GET("/", vehicleHandler.GetVehicles)
		vehicles.GET("/customer/:customerID", vehicleHandler.GetVehicleByCustomerID)
		vehicles.GET("/:id", vehicleHandler.GetVehicleByID)
		vehicles.GET("/plate/:plate", vehicleHandler.GetVehicleByPlate)
		vehicles.POST("/", vehicleHandler.CreateVehicle)
		vehicles.PATCH("/:id", vehicleHandler.UpdateVehicle)
		vehicles.DELETE("/:id", vehicleHandler.DeleteVehicle)
	}
}
