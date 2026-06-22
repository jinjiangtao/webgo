package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:64;not null;uniqueIndex" json:"username"`
	Color     string    `gorm:"size:16;not null" json:"color"`
	CreatedAt time.Time `json:"createdAt"`
}

type Whiteboard struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:128;not null" json:"name"`
	Background string        `gorm:"size:32;not null;default:'#ffffff'" json:"background"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Operation struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	WhiteboardID uint   `gorm:"not null;index" json:"whiteboardId"`
	UserID       uint   `gorm:"not null;index" json:"userId"`
	Type         string `gorm:"size:32;not null" json:"type"`
	Data         string `gorm:"type:text;not null" json:"data"`
	Timestamp    int64  `gorm:"not null;index" json:"timestamp"`
}

type Snapshot struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	WhiteboardID uint      `gorm:"not null;index" json:"whiteboardId"`
	Name         string    `gorm:"size:128;not null" json:"name"`
	Operations   string    `gorm:"type:text;not null" json:"operations"`
	CreatedAt    time.Time `json:"createdAt"`
}
