package services

import (
	"encoding/csv"
	"errors"
	"io"
	"jiansuo/database"
	"jiansuo/models"
	"strconv"
	"strings"
	"time"
)

type KeywordQueryParams struct {
	Query      string `form:"q"`
	CategoryID uint   `form:"category_id"`
	Status     *int   `form:"status"`
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
	SortBy     string `form:"sort_by"`
	SortOrder  string `form:"sort_order"`
}

type KeywordListResponse struct {
	List       []models.Keyword `json:"list"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

func CreateKeyword(kw *models.Keyword) error {
	if kw.Title == "" {
		return errors.New("title cannot be empty")
	}
	if kw.Content == "" {
		return errors.New("content cannot be empty")
	}
	kw.CreatedAt = time.Now()
	kw.UpdatedAt = time.Now()
	return database.DB.Create(kw).Error
}

func GetKeywordByID(id uint) (*models.Keyword, error) {
	var kw models.Keyword
	err := database.DB.Preload("Category").First(&kw, id).Error
	if err != nil {
		return nil, err
	}
	return &kw, nil
}

func UpdateKeyword(id uint, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	result := database.DB.Model(&models.Keyword{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("keyword not found")
	}
	return nil
}

func DeleteKeyword(id uint) error {
	result := database.DB.Delete(&models.Keyword{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("keyword not found")
	}
	return nil
}

func ListKeywords(params KeywordQueryParams) (*KeywordListResponse, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	if params.PageSize > 200 {
		params.PageSize = 200
	}

	db := database.DB.Model(&models.Keyword{}).Preload("Category")

	query := strings.TrimSpace(params.Query)
	if query != "" {
		db = db.Where("title LIKE ? OR content LIKE ? OR tags LIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%")
	}
	if params.CategoryID > 0 {
		db = db.Where("category_id = ?", params.CategoryID)
	}
	if params.Status != nil {
		db = db.Where("status = ?", *params.Status)
	}

	var total int64
	db.Count(&total)

	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "sort"
	}
	sortOrder := params.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}
	orderClause := sortBy + " " + sortOrder
	if sortBy == "created_at" || sortBy == "updated_at" || sortBy == "sort" || sortBy == "view_count" || sortBy == "id" {
		db = db.Order(orderClause)
	} else {
		db = db.Order("sort desc, id desc")
	}

	var list []models.Keyword
	offset := (params.Page - 1) * params.PageSize
	err := db.Offset(offset).Limit(params.PageSize).Find(&list).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / params.PageSize
	if int(total)%params.PageSize > 0 {
		totalPages++
	}

	return &KeywordListResponse{
		List:       list,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

func BatchCreateKeywords(keywords []models.Keyword) (int, error) {
	if len(keywords) == 0 {
		return 0, errors.New("no keywords to import")
	}

	successCount := 0
	now := time.Now()

	for i := range keywords {
		keywords[i].CreatedAt = now
		keywords[i].UpdatedAt = now
		err := database.DB.Create(&keywords[i]).Error
		if err == nil {
			successCount++
		}
	}

	return successCount, nil
}

func ImportCSV(reader io.Reader) (int, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	if err != nil {
		return 0, err
	}

	if len(records) <= 1 {
		return 0, errors.New("empty CSV file")
	}

	var keywords []models.Keyword
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 2 {
			continue
		}

		kw := models.Keyword{
			Title:   strings.TrimSpace(record[0]),
			Content: strings.TrimSpace(record[1]),
		}
		if len(record) > 2 {
			if catID, err := strconv.Atoi(strings.TrimSpace(record[2])); err == nil {
				kw.CategoryID = uint(catID)
			}
		}
		if len(record) > 3 {
			kw.Tags = strings.TrimSpace(record[3])
		}
		if len(record) > 4 {
			if status, err := strconv.Atoi(strings.TrimSpace(record[4])); err == nil {
				kw.Status = status
			} else {
				kw.Status = 1
			}
		} else {
			kw.Status = 1
		}
		if len(record) > 5 {
			if sort, err := strconv.Atoi(strings.TrimSpace(record[5])); err == nil {
				kw.Sort = sort
			}
		}

		if kw.Title != "" && kw.Content != "" {
			keywords = append(keywords, kw)
		}
	}

	return BatchCreateKeywords(keywords)
}

func SetKeywordStatus(id uint, status int) error {
	return UpdateKeyword(id, map[string]interface{}{"status": status})
}
