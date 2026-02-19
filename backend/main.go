package main

import (
	"log"

	"github.com/ashmit-singh-gogia/tts-backend/internal/database"
	"github.com/ashmit-singh-gogia/tts-backend/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Initialize the Database (This fills the "DB" box we talked about!)
	database.InitDB()
	log.Println("Database initialized successfully")

	// 2. Create the Gin Router
	r := gin.Default()
	r.Static("/audio", "./audio")
	// 3. Setup CORS so your React frontend can talk to it without errors
	r.Use(cors.Default())

	// 4. Register Routes (Connecting the route to your new Handler)
	r.POST("/upload", handlers.HandleTextUpload)

	// 5. Start the Server
	log.Println("Server is running on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
