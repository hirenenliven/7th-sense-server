package service

import (
	"context"
	"errors"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// CategoryService holds the business rules for categories.
type CategoryService struct {
	categoryRepository *repository.CategoryRepository
}

// NewCategoryService creates a CategoryService backed by categoryRepository.
func NewCategoryService(categoryRepository *repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepository: categoryRepository}
}

func (categoryService *CategoryService) GetAll(ctx context.Context) ([]model.Category, error) {
	return categoryService.categoryRepository.FindAll(ctx)
}

func (categoryService *CategoryService) GetByID(ctx context.Context, id int64) (model.Category, error) {
	category, err := categoryService.categoryRepository.FindByID(ctx, id)
	return category, toCategoryResponseError(err)
}

func (categoryService *CategoryService) Create(ctx context.Context, categoryRequest dto.CreateCategoryRequest) (model.Category, error) {
	category, err := categoryService.categoryRepository.Create(ctx, model.Category{
		Name:        categoryRequest.Name,
		Code:        categoryRequest.Code,
		Description: categoryRequest.Description,
	})
	return category, toCategoryResponseError(err)
}

func (categoryService *CategoryService) Update(ctx context.Context, id int64, categoryRequest dto.UpdateCategoryRequest) (model.Category, error) {
	category, err := categoryService.categoryRepository.Update(ctx, model.Category{
		ID:          id,
		Name:        categoryRequest.Name,
		Code:        categoryRequest.Code,
		Description: categoryRequest.Description,
	})
	return category, toCategoryResponseError(err)
}

func (categoryService *CategoryService) Delete(ctx context.Context, id int64) error {
	return toCategoryResponseError(categoryService.categoryRepository.Delete(ctx, id))
}

// toCategoryResponseError maps repository errors to category responses. Other errors
// pass through and are sent as ServerError.
func toCategoryResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.CategoryNotFound
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NameExists
	case errors.Is(err, repository.ErrDuplicateCode):
		return response.CodeExists
	default:
		return err
	}
}
