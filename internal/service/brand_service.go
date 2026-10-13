package service

import (
	"context"
	"errors"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// BrandService holds the business rules for brands.
type BrandService struct {
	brandRepository *repository.BrandRepository
}

// NewBrandService creates a BrandService backed by brandRepository.
func NewBrandService(brandRepository *repository.BrandRepository) *BrandService {
	return &BrandService{brandRepository: brandRepository}
}

func (brandService *BrandService) GetAll(ctx context.Context) ([]model.Brand, error) {
	return brandService.brandRepository.FindAll(ctx)
}

func (brandService *BrandService) GetByID(ctx context.Context, id int64) (model.Brand, error) {
	brand, err := brandService.brandRepository.FindByID(ctx, id)
	return brand, toBrandResponseError(err)
}

func (brandService *BrandService) Create(ctx context.Context, brandRequest dto.CreateBrandRequest) (model.Brand, error) {
	brand, err := brandService.brandRepository.Create(ctx, model.Brand{
		Name:        brandRequest.Name,
		Code:        brandRequest.Code,
		Description: brandRequest.Description,
	})
	return brand, toBrandResponseError(err)
}

func (brandService *BrandService) Update(ctx context.Context, id int64, brandRequest dto.UpdateBrandRequest) (model.Brand, error) {
	brand, err := brandService.brandRepository.Update(ctx, model.Brand{
		ID:          id,
		Name:        brandRequest.Name,
		Code:        brandRequest.Code,
		Description: brandRequest.Description,
	})
	return brand, toBrandResponseError(err)
}

func (brandService *BrandService) Delete(ctx context.Context, id int64) error {
	return toBrandResponseError(brandService.brandRepository.Delete(ctx, id))
}

// toBrandResponseError maps repository errors to brand responses. Other errors
// pass through and are sent as ServerError.
func toBrandResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.BrandNotFound
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NameExists
	case errors.Is(err, repository.ErrDuplicateCode):
		return response.CodeExists
	default:
		return err
	}
}
