package services

import (
	"errors"
	"jiansuo/database"
	"jiansuo/models"
	"jiansuo/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SearchParams struct {
	Query      string `form:"q" json:"q"`
	CategoryID uint   `form:"category_id" json:"category_id"`
	Status     *int   `form:"status" json:"status"`
	Page       int    `form:"page" json:"page"`
	PageSize   int    `form:"page_size" json:"page_size"`
	SortBy     string `form:"sort_by" json:"sort_by"`
	SortOrder  string `form:"sort_order" json:"sort_order"`
	SessionID  string `form:"session_id" json:"session_id"`
}

type SearchResponse struct {
	List       []models.SearchResult `json:"list"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	Tokens     []string              `json:"tokens"`
	TotalPages int                   `json:"total_pages"`
}

func Search(params SearchParams, userIP, userAgent string) (*SearchResponse, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}

	query := strings.TrimSpace(params.Query)
	tokens := utils.Tokenize(query)

	db := database.DB.Model(&models.Keyword{}).Preload("Category")

	db = db.Where("status = ?", 1)

	if params.CategoryID > 0 {
		db = db.Where("category_id = ?", params.CategoryID)
	}
	if params.Status != nil {
		db = db.Where("status = ?", *params.Status)
	}

	var allKeywords []models.Keyword
	if err := db.Find(&allKeywords).Error; err != nil {
		return nil, err
	}

	var scoredResults []models.SearchResult
	for _, kw := range allKeywords {
		score := 0.0
		if query != "" {
			score = utils.CalculateScore(kw.Title, kw.Content, kw.Tags, tokens, query)
			if score <= 0 {
				continue
			}
		} else {
			score = float64(kw.Sort) / 10
		}

		score += float64(kw.ViewCount) * 0.01

		highlightTokens := tokens
		if query != "" {
			highlightTokens = append(highlightTokens, query)
		}

		scoredResults = append(scoredResults, models.SearchResult{
			Keyword:            kw,
			TitleHighlighted:   utils.Highlight(kw.Title, highlightTokens, "mark"),
			ContentHighlighted: utils.Highlight(kw.Content, highlightTokens, "mark"),
			TagsHighlighted:    utils.Highlight(kw.Tags, highlightTokens, "mark"),
			Score:              score,
		})
	}

	if len(scoredResults) == 0 {
		scoredResults = []models.SearchResult{}
	}

	_ = params.SortBy
	switch strings.ToLower(params.SortOrder) {
	case "asc":
		for i := 0; i < len(scoredResults); i++ {
			for j := i + 1; j < len(scoredResults); j++ {
				if scoredResults[i].Score > scoredResults[j].Score {
					scoredResults[i], scoredResults[j] = scoredResults[j], scoredResults[i]
				}
			}
		}
	default:
		for i := 0; i < len(scoredResults); i++ {
			for j := i + 1; j < len(scoredResults); j++ {
				if scoredResults[i].Score < scoredResults[j].Score {
					scoredResults[i], scoredResults[j] = scoredResults[j], scoredResults[i]
				}
			}
		}
	}

	total := int64(len(scoredResults))
	start := (params.Page - 1) * params.PageSize
	end := start + params.PageSize

	var paginatedList []models.SearchResult
	if start < len(scoredResults) {
		if end > len(scoredResults) {
			end = len(scoredResults)
		}
		paginatedList = scoredResults[start:end]
	} else {
		paginatedList = []models.SearchResult{}
	}

	if query != "" || params.CategoryID > 0 {
		go recordSearchLog(query, params.CategoryID, int(total), userIP, userAgent, params.SessionID)
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &SearchResponse{
		List:       paginatedList,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		Tokens:     tokens,
		TotalPages: totalPages,
	}, nil
}

func GetSuggestions(query string, categoryID uint, limit int) []string {
	query = strings.TrimSpace(query)
	if query == "" {
		return []string{}
	}
	if limit <= 0 {
		limit = 10
	}

	db := database.DB.Model(&models.Keyword{}).Where("status = ?", 1)
	if categoryID > 0 {
		db = db.Where("category_id = ?", categoryID)
	}

	var titles []string
	db.Pluck("title", &titles)

	var tags []string
	var allTags []string
	database.DB.Model(&models.Keyword{}).Where("status = ?", 1).Pluck("tags", &tags)
	for _, t := range tags {
		if t != "" {
			tagParts := strings.Split(t, ",")
			for _, tp := range tagParts {
				tp = strings.TrimSpace(tp)
				if tp != "" {
					allTags = append(allTags, tp)
				}
			}
		}
	}

	allKeywords := append(titles, allTags...)
	suggestions := utils.GenerateSuggestions(allKeywords, query, limit)

	if suggestions == nil {
		suggestions = []string{}
	}

	return suggestions
}

func recordSearchLog(keyword string, categoryID uint, resultCount int, userIP, userAgent, sessionID string) {
	log := models.SearchLog{
		Keyword:     keyword,
		CategoryID:  categoryID,
		ResultCount: resultCount,
		UserIP:      userIP,
		UserAgent:   userAgent,
		SessionID:   sessionID,
		CreatedAt:   time.Now(),
	}
	database.DB.Create(&log)
}

func IncrementViewCount(keywordID uint) error {
	result := database.DB.Model(&models.Keyword{}).
		Where("id = ?", keywordID).
		UpdateColumn("view_count", gorm.Expr("view_count + 1"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("keyword not found")
	}
	return nil
}

func GetHotKeywords(limit int) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := database.DB.Table("search_logs").
		Select("keyword, COUNT(*) as search_count").
		Where("keyword != ''").
		Group("keyword").
		Order("search_count DESC").
		Limit(limit).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var keyword string
		var count int
		rows.Scan(&keyword, &count)
		results = append(results, map[string]interface{}{
			"keyword":      keyword,
			"search_count": count,
		})
	}

	return results, nil
}
