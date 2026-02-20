package routes

import (
	handlers "mecanica_xpto/internal/adapter/http/handlers"

	"github.com/gin-gonic/gin"
)

func addCustomerRoutes(rg *gin.RouterGroup, customerHandler *handlers.CustomerHandler) {
	customersRoutes := rg.Group(PathCustomers)
	{
		customersRoutes.GET("/:id", customerHandler.GetCustomerByID)
		customersRoutes.GET("/document/:document", customerHandler.GetCustomer)
		customersRoutes.GET("/full/:id", customerHandler.GetFullCustomer)
		customersRoutes.POST("", customerHandler.CreateCustomer)
		customersRoutes.PATCH("/:id", customerHandler.UpdateCustomer)
		customersRoutes.DELETE("/:id", customerHandler.DeleteCustomer)
		customersRoutes.GET("/", customerHandler.ListCustomer)
	}
}
