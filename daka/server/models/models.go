package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"size:255;not null" json:"name"`
	Description     string         `gorm:"type:text" json:"description"`
	CycleType       string         `gorm:"size:50;not null" json:"cycle_type"`
	CountdownSeconds int           `gorm:"not null" json:"countdown_seconds"`
	StartTime       string         `gorm:"size:20" json:"start_time"`
	EndTime         string         `gorm:"size:20" json:"end_time"`
	Color           string         `gorm:"size:20;default:#6366f1" json:"color"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type Record struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	TaskID     uint           `gorm:"not null;index" json:"task_id"`
	Task       Task           `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	RecordDate string         `gorm:"size:20;not null;index" json:"record_date"`
	Status     string         `gorm:"size:50;not null" json:"status"`
	CheckInTime *time.Time    `json:"check_in_time"`
	IsMakeup   bool           `gorm:"default:false" json:"is_makeup"`
	Note       string         `gorm:"type:text" json:"note"`
	Duration   int            `gorm:"default:0" json:"duration"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type TaskCreateRequest struct {
	Name              string `json:"name" binding:"required"`
	Description       string `json:"description"`
	CycleType         string `json:"cycle_type" binding:"required"`
	CountdownSeconds  int    `json:"countdown_seconds" binding:"required,min=1"`
	StartTime         string `json:"start_time"`
	EndTime           string `json:"end_time"`
	Color             string `json:"color"`
	IsActive          bool   `json:"is_active"`
}

type TaskUpdateRequest struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	CycleType         string `json:"cycle_type"`
	CountdownSeconds  int    `json:"countdown_seconds" binding:"omitempty,min=1"`
	StartTime         string `json:"start_time"`
	EndTime           string `json:"end_time"`
	Color             string `json:"color"`
	IsActive          *bool  `json:"is_active"`
}

type CheckInRequest struct {
	TaskID   uint   `json:"task_id" binding:"required"`
	IsMakeup bool   `json:"is_makeup"`
	Note     string `json:"note"`
	Duration int    `json:"duration"`
}

type TaskWithStatus struct {
	Task
	TodayRecord *Record `json:"today_record"`
}
