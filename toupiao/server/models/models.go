package models

import (
	"time"
)

type Activity struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `gorm:"not null" json:"start_time"`
	EndTime     time.Time `gorm:"not null" json:"end_time"`
	VoteType    string    `gorm:"not null;default:'single'" json:"vote_type"`
	MaxChoices  int       `gorm:"default:1" json:"max_choices"`
	Status      int       `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	Options     []Option  `gorm:"foreignKey:ActivityID" json:"options"`
}

type Option struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	ActivityID uint   `gorm:"not null;index" json:"activity_id"`
	Name       string `gorm:"not null" json:"name"`
	ImageURL   string `json:"image_url"`
	VoteCount  int    `gorm:"default:0" json:"vote_count"`
}

type VoteRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ActivityID uint      `gorm:"not null;index" json:"activity_id"`
	OptionID   uint      `gorm:"not null;index" json:"option_id"`
	IP         string    `gorm:"not null;index" json:"ip"`
	Cookie     string    `gorm:"index" json:"cookie"`
	CreatedAt  time.Time `json:"created_at"`
	OptionName string    `json:"option_name"`
}
