package routes

import (
	"github.com/gin-gonic/gin"

	"internal-transfer-microservice/internal/controller"
)

// SetupAccountRoutes sets up the account routes
func SetupAccountRoutes(router *gin.Engine, accountController *controller.AccountController) {
	accountRoutes := router.Group("/api/v1/accounts")
	{
		accountRoutes.GET("/:id", accountController.GetAccount)
		accountRoutes.POST("", accountController.CreateAccount)
		accountRoutes.POST("/transfer", accountController.TransferMoney)
	}
}
