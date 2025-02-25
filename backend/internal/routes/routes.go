package routes

import (
	"ecommerce-price-tracker/internal/handlers"
	middleware "ecommerce-price-tracker/pkg/middlewares"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	api.POST("/register", handlers.CreateUser)
	api.POST("/login", handlers.LoginUser)
	api.POST("/testing", func(ctx *gin.Context) {
		uid, _ := ctx.Get("userID")
		email, _ := ctx.Get("email")
		role, _ := ctx.Get("role")
		fmt.Print(ctx.Get("userID"))
		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("this is just a test, you're %v %s %v ", uid, email, role),
		})
	})
}
