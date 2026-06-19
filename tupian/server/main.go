package main

import (
	"log"

	"image-editor-server/database"
	"image-editor-server/routes"
)

func main() {
	database.Init()

	r := routes.SetupRouter()

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
