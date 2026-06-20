package services

import (
	"errors"
	"jiansuo/database"
	"jiansuo/models"
	"time"
)

func CreateCategory(cat *models.Category) error {
	if cat.Name == "" {
		return errors.New("category name cannot be empty")
	}
	cat.CreatedAt = time.Now()
	cat.UpdatedAt = time.Now()
	return database.DB.Create(cat).Error
}

func GetCategoryByID(id uint) (*models.Category, error) {
	var cat models.Category
	err := database.DB.First(&cat, id).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func UpdateCategory(id uint, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	result := database.DB.Model(&models.Category{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}

func DeleteCategory(id uint) error {
	var count int64
	database.DB.Model(&models.Keyword{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("cannot delete category with associated keywords")
	}

	result := database.DB.Delete(&models.Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}

func ListCategories(includeDisabled bool) ([]models.Category, error) {
	var list []models.Category
	db := database.DB
	if !includeDisabled {
		db = db.Where("status = ?", 1)
	}
	err := db.Order("id asc").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func SetCategoryStatus(id uint, status int) error {
	return UpdateCategory(id, map[string]interface{}{"status": status})
}
