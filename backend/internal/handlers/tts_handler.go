package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ashmit-singh-gogia/tts-backend/internal/database"
	"github.com/ashmit-singh-gogia/tts-backend/internal/models"
	"github.com/ashmit-singh-gogia/tts-backend/internal/services"
	"github.com/gin-gonic/gin"
)

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
	entry := models.History{}
	result := database.DB.Create(&entry)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}
	text := ""
	ext := filepath.Ext(fileHeader.Filename)
	switch ext {
	case ".pdf":
		// 1. Get the system's safe temporary folder (e.g., /var/folders/... on Mac)
		tempDir := "./tmp"

		// Ensure the folder actually exists before we try to save to it
		err = os.MkdirAll(tempDir, os.ModePerm)
		if err != nil {
			fmt.Println("FOLDER CREATION ERROR:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create temp directory"})
			return
		}
		tempFilePath := filepath.Join(tempDir, fmt.Sprintf("%d_temp.pdf", entry.ID))

		// 3. Save the file there safely
		err = c.SaveUploadedFile(fileHeader, tempFilePath)
		if err != nil {
			fmt.Println("FILE SAVE ERROR:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save temporary file"})
			return
		}
		// 4. Pass that exact path to your extractor
		text, err = services.ExtractTextFromPDF(tempFilePath)
		if err != nil {
			fmt.Println("EXTRACTION ERROR:", err) // <--- ADD THIS LINE
			c.String(http.StatusInternalServerError, "Error reading file")
			return
		}
		// 5. Clean up the file IMMEDIATELY after extraction
		os.Remove(tempFilePath)
	case ".txt":
		content, err := io.ReadAll(file)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading file")
			return
		}
		text = string(content)
	default:
		// FAIL FAST: Unsupported file type
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type. Please upload .txt or .pdf"})
		return
	}

	finalFileName, err := services.GetFinalAudio(text, int(entry.ID))
	if err != nil {
		fmt.Println("TTS ERROR:", err) // <--- ADD THIS LINE
		c.JSON(500, gin.H{"error": "Some Internal error"})
		return
	}
	database.DB.Model(&entry).Update("Filename", finalFileName)

	publicURL := "/" + finalFileName
	c.JSON(200, gin.H{
		"status":    "success",
		"audio_url": publicURL,
		"id":        entry.ID,
	})
}
