package main

import (
	"github.com/ashmit-singh-gogia/tts-backend/internal/database"
	"github.com/ashmit-singh-gogia/tts-backend/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TTSHistory struct {
	gorm.Model
	Text string
}
type TTSRequest struct {
	Text string `json:"text"`
}

// Helper Functions //

func main() {
	database.InitDB()
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:5173",
		"https://tts-project-pjgs4a2q4-ashmit-singh-gogias-projects.vercel.app", // Your exact live URL!
	}
	router.Use(cors.New(config))
	router.Static("/audio", "./audio")
	router.POST("/upload", func(c *gin.Context) {
		handlers.HandleTextUpload(c, database.DB)
	})
	router.Run(":8080")
}
