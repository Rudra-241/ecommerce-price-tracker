package main

import (
	"ecommerce-price-tracker/internal/db"
	"ecommerce-price-tracker/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routes.SetupRoutes(r)

	db.InitWithDSN("host=localhost user=rudraaa password=admin dbname=webscraper port=5432 sslmode=disable")

	r.Run("localhost:3000")
}
