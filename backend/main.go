package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TTSHistory struct {
	gorm.Model
	Text string
}
type TTSRequest struct {
	Text string `json:"text"`
}

func main() {

	// github.com/mattn/go-sqlite3
	db, err := gorm.Open(sqlite.Open("tts.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	db.AutoMigrate(&TTSHistory{})
	db.AutoMigrate(&TTSRequest{})
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/history", func(c *gin.Context) {
		var history []TTSHistory
		result := db.Find(&history)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": result.Error})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"Data": history})
	})

	router.POST("/createRequest", func(c *gin.Context) {
		var json TTSRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newRecord := TTSHistory{
			Text: json.Text,
		}
		result := db.Create(&newRecord)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": result.Error})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "saved"})
	})
	router.Run(":8080")
}
