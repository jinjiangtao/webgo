package models

import (
	"time"
)

type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Description string    `gorm:"size:500" json:"description"`
	Status      int       `gorm:"default:1;index" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Keywords    []Keyword `gorm:"foreignKey:CategoryID" json:"keywords,omitempty"`
}

type Keyword struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:200;not null;index" json:"title"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	CategoryID  uint      `gorm:"index" json:"category_id"`
	Category    *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags        string    `gorm:"size:500" json:"tags"`
	Status      int       `gorm:"default:1;index" json:"status"`
	Sort        int       `gorm:"default:0;index" json:"sort"`
	ViewCount   int       `gorm:"default:0" json:"view_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SearchLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Keyword    string    `gorm:"size:200;index" json:"keyword"`
	UserIP     string    `gorm:"size:50" json:"user_ip"`
	UserAgent  string    `gorm:"size:500" json:"user_agent"`
	ResultCount int      `gorm:"default:0" json:"result_count"`
	CategoryID uint      `gorm:"index" json:"category_id,omitempty"`
	SessionID  string    `gorm:"size:100;index" json:"session_id"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
}

type SearchResult struct {
	Keyword      Keyword `json:"keyword"`
	TitleHighlighted   string `json:"title_highlighted"`
	ContentHighlighted string `json:"content_highlighted"`
	TagsHighlighted    string `json:"tags_highlighted"`
	Score        float64 `json:"score"`
}
