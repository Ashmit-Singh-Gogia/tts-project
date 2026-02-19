package models

import "gorm.io/gorm"

type History struct {
	gorm.Model
	Filename string `json:"filename"` // e.g., "audio/105_final.mp3"
}
