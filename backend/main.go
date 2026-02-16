package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TTSHistory struct {
	gorm.Model
	Text string
}

func main() {

	// github.com/mattn/go-sqlite3
	db, err := gorm.Open(sqlite.Open("tts.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	db.AutoMigrate(&TTSHistory{})
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run(":8080")
}
