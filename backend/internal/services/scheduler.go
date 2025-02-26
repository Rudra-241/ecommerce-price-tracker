package services

import (
	"fmt"
	"time"
)

func RunUpdaterJob() {
	intervalHours := 2
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
