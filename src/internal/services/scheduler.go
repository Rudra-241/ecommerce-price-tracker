package services

import (
	"ecommerce-price-tracker/internal/db"
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
		//updateProducts()
		emailAll()
	}
}

func updateProducts() {
	UpdateAll()
	fmt.Println("Products updated at:", time.Now().Format(time.RFC1123))

}

func emailAll() {
	dbb := db.GetDB()
	err := EmailAll(dbb)
	fmt.Println("EmailAll:", err)
	if err != nil {
		fmt.Println("Unable to send emails:", err)
	}
	fmt.Println("Emailed recepients at: ", time.Now().Format(time.RFC1123))
}
