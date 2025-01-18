package main

import (
	"log"
	"receipt-processor/internal/handlers"
	"receipt-processor/internal/store"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initializing storage
	storage := store.NewMemoryStore()

	// Initializing router
	router := gin.Default()

	// Initializing handlers
	receiptHandler := handlers.NewReceiptHandler(storage)

	// Routes
	router.POST("/receipts/process", receiptHandler.ProcessReceipt)
	router.GET("/receipts/:id/points", receiptHandler.GetPoints)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
