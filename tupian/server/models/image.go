package models

import (
	"time"

	"gorm.io/gorm"
)

type Image struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UUID         string         `gorm:"uniqueIndex;size:36;not null" json:"uuid"`
	Filename     string         `gorm:"size:255;not null" json:"filename"`
	OriginalName string         `gorm:"size:255;not null" json:"original_name"`
	FilePath     string         `gorm:"size:500;not null" json:"file_path"`
	FileSize     int64          `json:"file_size"`
	ContentType  string         `gorm:"size:100" json:"content_type"`
	Width        int            `json:"width"`
	Height       int            `json:"height"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
