package database

import (
	"log"

	"github.com/ashmit-singh-gogia/tts-backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// We export this DB variable so other packages can use it later if needed
var DB *gorm.DB

func InitDB() *gorm.DB {
	var err error
	DB, err = gorm.Open(sqlite.Open("tts.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	err = DB.AutoMigrate(&models.History{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return DB
}
