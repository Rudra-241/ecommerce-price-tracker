package services

import (
	"fmt"
	"time"
)

func RunUpdaterJob(h int) {
	intervalHours := h
	fmt.Printf("Running updater job in every: %d\n", intervalHours)
	interval := time.Duration(intervalHours) * time.Hour
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		updateProducts()
	}
}

func updateProducts() {
	UpdateAll()
	fmt.Println("Products updated at:", time.Now().Format(time.RFC1123))

}
