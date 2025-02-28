package main

import (
	"ecommerce-price-tracker/internal/db"
	"ecommerce-price-tracker/internal/routes/api"
	"ecommerce-price-tracker/internal/routes/web"
	"ecommerce-price-tracker/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
	GinMode := os.Getenv("GIN_MODE")
	if GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.Default()

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		fmt.Println("Error setting trusted proxies")
		return
	}
	api.SetupAPIRoutes(r)
	web.SetUpWebRoutes(r)

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		fmt.Println("Error: Required database environment variables are not set")
		fmt.Println("Please set DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, and DB_PORT")
		return
	}

	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	db.InitWithDSN(dsn)
	updateIn, _ := strconv.Atoi(os.Getenv("UPDATE_IN"))
	go services.RunUpdaterJob(updateIn)
	if GinMode == "debug" {
		go services.UpdateAll()
	}
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "3000"
	}

	err := r.Run("localhost:" + serverPort)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
}
