package router

import (
	"example-loan.com/loan-app/app/config"
	"github.com/gin-gonic/gin"
)

func InitRouter(initApp *config.Initialization) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	api := router.Group("/api")
	api.POST("/loan/request", initApp.LoanController.RequestLoan)
	api.POST("/loan/approve", initApp.LoanController.ApproveLoan)

	return router
}
