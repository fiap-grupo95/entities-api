//go:build legacy
// +build legacy

package routes

import (
	handlers "mecanica_xpto/internal/adapter/http/handlers"

	"github.com/gin-gonic/gin"
)

func addPaymentRoutes(rg *gin.RouterGroup, paymentHandler *handlers.PaymentHandler) {

	payments := rg.Group(PathPayments)
	{
		payments.GET("/:id", paymentHandler.GetPaymentByID)
		payments.GET("/", paymentHandler.ListPayments)
		payments.POST("/", paymentHandler.CreatePayment)
	}
}
