package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetUpWebRoutes(r *gin.Engine) {
	r.Static("/styles", "public/styles")
	r.Static("/assets", "public/assets")
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
	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", nil)
	})
}
