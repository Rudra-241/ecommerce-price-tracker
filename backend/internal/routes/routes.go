package routes

import (
	"ecommerce-price-tracker/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.Use(middlewares.AuthMiddleware())

	UserRoutes(api)
	ProductRoutes(api)
}
