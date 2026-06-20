package services

import (
	"jiansuo/database"
	"jiansuo/models"
	"time"
)

type SearchLogQueryParams struct {
	Keyword   string `form:"keyword"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	SessionID string `form:"session_id"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
}

type SearchLogListResponse struct {
	List       []models.SearchLog `json:"list"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

func ListSearchLogs(params SearchLogQueryParams) (*SearchLogListResponse, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 200 {
		params.PageSize = 200
	}

	db := database.DB.Model(&models.SearchLog{})

	if params.Keyword != "" {
		db = db.Where("keyword LIKE ?", "%"+params.Keyword+"%")
	}
	if params.SessionID != "" {
		db = db.Where("session_id = ?", params.SessionID)
	}
	if params.StartDate != "" {
		if t, err := time.Parse("2006-01-02", params.StartDate); err == nil {
			db = db.Where("created_at >= ?", t)
		}
	}
	if params.EndDate != "" {
		if t, err := time.Parse("2006-01-02", params.EndDate); err == nil {
			t = t.Add(24 * time.Hour)
			db = db.Where("created_at < ?", t)
		}
	}

	var total int64
	db.Count(&total)

	var list []models.SearchLog
	offset := (params.Page - 1) * params.PageSize
	err := db.Order("created_at desc").Offset(offset).Limit(params.PageSize).Find(&list).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &SearchLogListResponse{
		List:       list,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

type SearchStatistics struct {
	TotalSearches   int64                    `json:"total_searches"`
	TodaySearches   int64                    `json:"today_searches"`
	UniqueKeywords  int64                    `json:"unique_keywords"`
	TotalKeywords   int64                    `json:"total_keywords"`
	TotalCategories int64                    `json:"total_categories"`
	HotKeywords     []map[string]interface{} `json:"hot_keywords"`
	DailyTrend      []map[string]interface{} `json:"daily_trend"`
}

func GetSearchStatistics() (*SearchStatistics, error) {
	stats := &SearchStatistics{}

	database.DB.Model(&models.SearchLog{}).Count(&stats.TotalSearches)

	today := time.Now().Truncate(24 * time.Hour)
	database.DB.Model(&models.SearchLog{}).Where("created_at >= ?", today).Count(&stats.TodaySearches)

	database.DB.Model(&models.SearchLog{}).Distinct("keyword").Where("keyword != ''").Count(&stats.UniqueKeywords)

	database.DB.Model(&models.Keyword{}).Count(&stats.TotalKeywords)
	database.DB.Model(&models.Category{}).Count(&stats.TotalCategories)

	hotKeywords, _ := GetHotKeywords(10)
	stats.HotKeywords = hotKeywords

	rows, err := database.DB.Table("search_logs").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Group("DATE(created_at)").
		Order("date DESC").
		Limit(30).
		Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var date string
			var count int
			rows.Scan(&date, &count)
			stats.DailyTrend = append(stats.DailyTrend, map[string]interface{}{
				"date":  date,
				"count": count,
			})
		}
	}

	return stats, nil
}
