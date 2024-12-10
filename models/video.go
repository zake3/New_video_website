package models

import (
	"time"
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Title       string    `json:"title"`
	Description string    `json:"description"`
	FilePath    string    `json:"file_path"`
	Thumbnail   string    `json:"thumbnail"`
	UserID      uint      `json:"user_id"`
	ViewCount   int       `json:"view_count"`
	UploadDate  time.Time `json:"upload_date"`
}