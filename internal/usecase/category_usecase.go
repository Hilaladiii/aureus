package usecase

import (
	"context"
	"errors"

	"github.com/Hilaladiii/aureus/internal/model"
	"github.com/Hilaladiii/aureus/internal/repository"
	"github.com/Hilaladiii/aureus/pkg/exception"
	"gorm.io/gorm"
)

type CategoryUsecaseItf interface {
	CreateCategory(ctx context.Context, req *model.CategoryCreateRequest) (model.CategoryResource, error)
	UpdateCategory(ctx context.Context, req *model.CategoryUpdateRequest, categoryID string) (model.CategoryResource, error)
	DeleteCategory(ctx context.Context, categoryID string) (model.CategoryResource, error)
	GetAll(ctx context.Context) ([]model.CategoryResource, error)
	GetByID(ctx context.Context, categoryID string) (model.CategoryResource, error)
}

type CategoryUsecase struct {
	categoryRepo repository.CategoryRepoItf
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepoItf) *CategoryUsecase {
	return &CategoryUsecase{
		categoryRepo: categoryRepo,
	}
}

func (u *CategoryUsecase) CreateCategory(ctx context.Context, req *model.CategoryCreateRequest) (model.CategoryResource, error) {
	category, err := u.categoryRepo.GetByName(ctx, req.Name)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.CategoryResource{}, err
	}
	if category != nil {
		return model.CategoryResource{}, exception.NewBadRequestError("category already exists")
	}

	newCategory := model.Category{
		Name:        req.Name,
		Description: &req.Description,
	}

	err = u.categoryRepo.CreateCategory(ctx, &newCategory)
	if err != nil {
		return model.CategoryResource{}, err
	}

	return newCategory.Resource(), nil
}

func (u *CategoryUsecase) UpdateCategory(ctx context.Context, req *model.CategoryUpdateRequest, categoryID string) (model.CategoryResource, error) {
	category, err := u.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.CategoryResource{}, exception.NewNotFoundError("category not found")
		}
		return model.CategoryResource{}, err
	}

	if req.Name != "" {
		category.Name = req.Name
	}

	if req.Description != "" {
		category.Description = &req.Description
	}

	err = u.categoryRepo.UpdateCategory(ctx, category, categoryID)
	if err != nil {
		return model.CategoryResource{}, err
	}

	return category.Resource(), nil
}

func (u *CategoryUsecase) DeleteCategory(ctx context.Context, categoryID string) (model.CategoryResource, error) {
	category, err := u.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.CategoryResource{}, exception.NewNotFoundError("category not found")
		}
		return model.CategoryResource{}, nil
	}

	if category == nil {
		return model.CategoryResource{}, errors.New("category doesn't exists")
	}

	if err := u.categoryRepo.DeleteCategory(ctx, categoryID); err != nil {
		return model.CategoryResource{}, nil
	}

	return category.Resource(), nil
}

func (u *CategoryUsecase) GetAll(ctx context.Context) ([]model.CategoryResource, error) {
	categories, err := u.categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return model.CategoryResources(categories), nil
}

func (u *CategoryUsecase) GetByID(ctx context.Context, categoryID string) (model.CategoryResource, error) {
	category, err := u.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.CategoryResource{}, exception.NewNotFoundError("category not found")
		}
		return model.CategoryResource{}, err
	}

	return category.Resource(), nil
}
