package web

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetUpWebRoutes(r *gin.Engine) {
	r.Static("/assets", "public/assets")
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.File("public/index.html")
	})
}
