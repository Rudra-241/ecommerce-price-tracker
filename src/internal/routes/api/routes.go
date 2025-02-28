package api

import (
	"ecommerce-price-tracker/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupAPIRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.Use(middlewares.AuthMiddleware())

	UserRoutes(api)
	ProductRoutes(api)
}
