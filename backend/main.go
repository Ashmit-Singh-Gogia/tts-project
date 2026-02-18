package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
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
type History struct {
	gorm.Model
	Filename string `json:"filename"` // e.g., "audio/105_final.mp3"
}

// Helper Functions //
func CustomSplit(text string, maxLength int) []string {
	ans := []string{}

	for len(text) > 0 {
		text = strings.TrimSpace(text) // Always clean up leading spaces

		if len(text) <= maxLength {
			ans = append(ans, text)
			return ans
		}

		// 1. Look at our window
		searchArea := text[:maxLength]

		// 2. Try to find a Period
		cutPoint := strings.LastIndex(searchArea, ".")

		// 3. If no period, try to find a Space (so we don't break words)
		if cutPoint == -1 {
			cutPoint = strings.LastIndex(searchArea, " ")
		}

		// 4. Perform the cut
		if cutPoint == -1 {
			// Absolute fallback: No period AND no space? Hard cut at maxLength.
			ans = append(ans, text[:maxLength])
			text = text[maxLength:]
		} else {
			ans = append(ans, text[:cutPoint+1])
			text = text[cutPoint+1:]
		}
	}
	return ans

}

func mergeAudio(files []string, output string) error {
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		// Write the binary data to the master file
		out.Write(data)
	}
	return nil
}

func main() {
	db, err := gorm.Open(sqlite.Open("tts.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	db.AutoMigrate(&TTSHistory{})
	db.AutoMigrate(&History{})
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
		uniqueID := fmt.Sprintf("%d", entry.ID)
		finalFileName := "audio/" + uniqueID + "_final.mp3"

		// 5. Run the Smart Chunker Logic
		chunks := CustomSplit(text, 200)
		speech := htgotts.Speech{Folder: "audio", Language: voices.English, Handler: &handlers.Native{}}

		var chunkFiles []string

		for i, chunk := range chunks {
			// Name the chunk: "105_part_0"
			chunkName := fmt.Sprintf("%s_part_%d", uniqueID, i)

			speech.CreateSpeechFile(chunk, chunkName)

			// Add the extension manually because the library adds it secretly!
			chunkFiles = append(chunkFiles, "audio/"+chunkName+".mp3")
		}

		// 6. Merge the chunks
		err = mergeAudio(chunkFiles, finalFileName)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to merge audio"})
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
