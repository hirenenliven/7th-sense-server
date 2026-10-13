package service

import (
	"context"
	"errors"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// defaultProductImage is stored when a product is saved without an image.
const defaultProductImage = "no_image_available.png"

// ProductService holds the business rules for products.
type ProductService struct {
	productRepository *repository.ProductRepository
}

// NewProductService creates a ProductService backed by productRepository.
func NewProductService(productRepository *repository.ProductRepository) *ProductService {
	return &ProductService{productRepository: productRepository}
}

func (productService *ProductService) GetAll(ctx context.Context) ([]model.Product, error) {
	return productService.productRepository.FindAll(ctx)
}

func (productService *ProductService) GetByID(ctx context.Context, id int64) (model.Product, error) {
	product, err := productService.productRepository.FindByID(ctx, id)
	return product, toProductResponseError(err)
}

func (productService *ProductService) Create(ctx context.Context, productRequest dto.CreateProductRequest) (model.Product, error) {
	product, err := productService.productRepository.Create(ctx, model.Product{
		Name:        productRequest.Name,
		Code:        productRequest.Code,
		CategoryID:  productRequest.CategoryID,
		Image:       productImageOrDefault(productRequest.Image),
		Description: productRequest.Description,
	})
	return product, toProductResponseError(err)
}

func (productService *ProductService) Update(ctx context.Context, id int64, productRequest dto.UpdateProductRequest) (model.Product, error) {
	product, err := productService.productRepository.Update(ctx, model.Product{
		ID:          id,
		Name:        productRequest.Name,
		Code:        productRequest.Code,
		CategoryID:  productRequest.CategoryID,
		Image:       productImageOrDefault(productRequest.Image),
		Description: productRequest.Description,
	})
	return product, toProductResponseError(err)
}

func (productService *ProductService) Delete(ctx context.Context, id int64) error {
	return toProductResponseError(productService.productRepository.Delete(ctx, id))
}

// productImageOrDefault returns image, or the default product image when it is empty.
func productImageOrDefault(image string) string {
	if image == "" {
		return defaultProductImage
	}
	return image
}

// toProductResponseError maps repository errors to product responses. Other errors
// pass through and are sent as ServerError.
func toProductResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.ProductNotFound
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NameExists
	case errors.Is(err, repository.ErrDuplicateCode):
		return response.CodeExists
	case errors.Is(err, repository.ErrInvalidCategory):
		return response.ValidCategory
	default:
		return err
	}
}
