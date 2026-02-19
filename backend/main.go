package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/voices"
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
	router := gin.Default()
	router.Use(cors.Default())
	router.Static("/audio", "./audio")
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
		speech := htgotts.Speech{Folder: "audio", Language: voices.English}
		fileName := strconv.Itoa(int(newRecord.ID))
		_, err := speech.CreateSpeechFile(newRecord.Text, fileName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		publicURL := fmt.Sprintf("/audio/%s.mp3", fileName)

		c.JSON(http.StatusOK, gin.H{
			"status":    "saved",
			"audio_url": publicURL,
		})
	})

	router.POST("/upload", func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "Upload failed: "+err.Error())
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			c.String(http.StatusInternalServerError, "Error opening file")
			return
		}
		defer file.Close() // Always close the file when done!
		content, err := io.ReadAll(file)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading file")
			return
		}
		text := string(content)
		entry := History{}
		result := db.Create(&entry)
		if result.Error != nil {
			c.JSON(500, gin.H{"error": "Database error"})
			return
		}

		// 7. UPDATE: Save the final filename to the History table
		db.Model(&entry).Update("Filename", finalFileName)

		// 8. Respond to Frontend
		// We strip "audio/" for the URL so the browser can load it
		publicURL := "/" + finalFileName
		c.JSON(200, gin.H{
			"status":    "success",
			"audio_url": publicURL,
			"id":        entry.ID,
		})
	})
	router.Run(":8080")
}
