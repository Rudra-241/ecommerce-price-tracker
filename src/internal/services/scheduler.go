package services

import (
	"ecommerce-price-tracker/internal/db"
	"log/slog"
	"time"
)

func RunUpdaterJob(h int) {
	intervalHours := h
	slog.Info("updater job scheduled", "interval_hours", intervalHours)
	interval := time.Duration(intervalHours) * time.Hour
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		updateProducts()
		emailAll()
	}
}

func updateProducts() {
	UpdateAll()
	slog.Info("products updated")
}

func emailAll() {
	dbb := db.GetDB()
	if err := EmailAll(dbb); err != nil {
		slog.Error("unable to queue email alerts", "err", err)
		return
	}
	slog.Info("email alerts queued")
}
