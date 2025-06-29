package main

import (
	"ecommerce-price-tracker/internal/db"
	"ecommerce-price-tracker/internal/queue"
	"ecommerce-price-tracker/internal/routes/api"
	"ecommerce-price-tracker/internal/routes/web"
	"ecommerce-price-tracker/internal/services"
	"ecommerce-price-tracker/internal/services/scraper"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"strconv"
)

func main() {
	job := flag.String("job", "", "run a one-off job and exit (e.g. seed-selectors)")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not loaded", "err", err)
	}
	GinMode := os.Getenv("GIN_MODE")
	if GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.Default()

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		slog.Error("setting trusted proxies", "err", err)
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
		slog.Error("required database environment variables are not set",
			"required", "DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT")
		return
	}

	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	db.InitWithDSN(dsn)

	if *job != "" {
		runJob(*job)
		return
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}
	if err := queue.Init(rabbitURL); err != nil {
		slog.Error("connecting to rabbitmq", "err", err)
		return
	}
	go func() {
		if err := services.ConsumeEmailAlerts(); err != nil {
			slog.Error("email consumer stopped", "err", err)
		}
	}()

	updateIn, _ := strconv.Atoi(os.Getenv("UPDATE_IN"))
	go services.RunUpdaterJob(updateIn)
	if GinMode == "debug" {
		go func() {
			err := services.EmailAll(db.GetDB())
			if err != nil {
				slog.Error("email send failed", "err", err)
			}
		}()
		go services.UpdateAll()
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "3000"
	}

	err := r.Run("localhost:" + serverPort)
	if err != nil {
		slog.Error("starting server", "err", err)
		return
	}
}

func runJob(name string) {
	switch name {
	case "seed-selectors":
		n, err := scraper.SeedSelectors(db.GetDB())
		if err != nil {
			fmt.Printf("seed-selectors failed: %v\n", err)
			return
		}
		fmt.Printf("seed-selectors: created %d row(s)\n", n)
	default:
		fmt.Printf("unknown job: %q\n", name)
	}
}
