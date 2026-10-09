package service

import (
	"context"
	"errors"
	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// ItemService holds the business rules for items.
type ItemService struct {
	repo repository.ItemRepository
}

// NewItemService creates an ItemService backed by repo.
func NewItemService(repo repository.ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) GetAll(ctx context.Context) ([]model.Item, error) {
	return s.repo.FindAll(ctx)
}

func (s *ItemService) GetByID(ctx context.Context, id string) (model.Item, error) {
	item, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return model.Item{}, response.ErrItemNotFound
	}
	return item, err
}

func (s *ItemService) Create(ctx context.Context, req dto.ItemRequest) (model.Item, error) {
	return s.repo.Create(ctx, model.Item{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
}

func (s *ItemService) Update(ctx context.Context, id string, req dto.ItemRequest) (model.Item, error) {
	item, err := s.repo.Update(ctx, model.Item{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if errors.Is(err, repository.ErrNotFound) {
		return model.Item{}, response.ErrItemNotFound
	}
	return item, err
}

func (s *ItemService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return response.ErrItemNotFound
	}
	return err
}
