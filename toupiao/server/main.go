package main

import (
	"log"
	"voting-system/config"
	"voting-system/routers"
)

func main() {
	config.InitDB()

	r := routers.SetupRouter()

	log.Println("Server starting on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
