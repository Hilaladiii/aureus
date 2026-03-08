package repository

import (
	"context"

	"github.com/Hilaladiii/aureus/internal/model"
	"gorm.io/gorm"
)

type CategoryRepoItf interface {
	CreateCategory(ctx context.Context, category *model.Category) error
	GetAll(ctx context.Context) ([]model.Category, error)
	GetByID(ctx context.Context, categoryID string) (*model.Category, error)
	GetByName(ctx context.Context, categoryName string) (*model.Category, error)
	UpdateCategory(ctx context.Context, category *model.Category, categoryID string) error
	DeleteCategory(ctx context.Context, categoryID string) error
}

type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{db}
}

func (r *CategoryRepo) CreateCategory(ctx context.Context, category *model.Category) error {
	err := r.db.WithContext(ctx).Create(category).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepo) GetAll(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.WithContext(ctx).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, categoryID string) (*model.Category, error) {
	var category model.Category
	err := r.db.WithContext(ctx).First(&category, "id = ?", categoryID).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepo) GetByName(ctx context.Context, categoryName string) (*model.Category, error) {
	var category model.Category
	err := r.db.WithContext(ctx).First(&category, "name = ?", categoryName).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepo) UpdateCategory(ctx context.Context, category *model.Category, categoryID string) error {
	if err := r.db.WithContext(ctx).Save(category).Error; err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepo) DeleteCategory(ctx context.Context, categoryID string) error {
	if err := r.db.WithContext(ctx).Delete(&model.Category{}, "id = ?", categoryID).Error; err != nil {
		return err
	}

	return nil
}
