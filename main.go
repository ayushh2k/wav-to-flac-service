// main.go
package main

import (
	"log"

	"wav-to-flac-service/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Define routes
	r.POST("/convert", handlers.ConvertWavToFlac)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
