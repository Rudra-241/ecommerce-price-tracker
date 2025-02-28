package web

import (
	"ecommerce-price-tracker/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetUpWebRoutes(r *gin.Engine) {
	r.Static("/styles", "public/styles")
	r.Static("/assets", "public/assets")
	r.Static("/scripts", "public/scripts")

	r.LoadHTMLGlob("public/*.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})
	r.GET("/dashboard", middlewares.AuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", nil)
	})
}
