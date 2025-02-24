package services

import (
	"errors"
	"gorm.io/gorm"
	"pionira/internal/models"
)

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{db: db}
}

func (c CategoryService) List() ([]*models.CategoryModel, error) {
	var categories []*models.CategoryModel
	result := c.db.Find(&categories)
	if result.Error != nil {
		return nil, errors.New("failed to fetch all categories")
	}

	return categories, nil
}

func (c CategoryService) Create(newCategory models.CategoryModel) error {
	return nil
}
