package handlers

import (
	"io"
	"net/http"

	"github.com/ashmit-singh-gogia/tts-backend/internal/database"
	"github.com/ashmit-singh-gogia/tts-backend/internal/models"
	"github.com/ashmit-singh-gogia/tts-backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB = database.InitDB()

func HandleTextUpload(c *gin.Context) {
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
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading file")
		return
	}
	text := string(content)
	entry := models.History{}
	result := DB.Create(&entry)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	finalFileName, err := services.GetFinalAudio(text, int(entry.ID))
	if err != nil {
		c.JSON(500, gin.H{"error": "Some Internal error"})
		return
	}
	DB.Model(&entry).Update("Filename", finalFileName)

	publicURL := "/" + finalFileName
	c.JSON(200, gin.H{
		"status":    "success",
		"audio_url": publicURL,
		"id":        entry.ID,
	})
}
