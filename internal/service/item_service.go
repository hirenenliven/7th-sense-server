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
	itemRepository repository.ItemRepository
}

// NewItemService creates an ItemService backed by itemRepository.
func NewItemService(itemRepository repository.ItemRepository) *ItemService {
	return &ItemService{itemRepository: itemRepository}
}

func (itemService *ItemService) GetAll(ctx context.Context) ([]model.Item, error) {
	return itemService.itemRepository.FindAll(ctx)
}

func (itemService *ItemService) GetByID(ctx context.Context, id string) (model.Item, error) {
	item, err := itemService.itemRepository.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return model.Item{}, response.ItemNotFound
	}
	return item, err
}

func (itemService *ItemService) Create(ctx context.Context, itemRequest dto.ItemRequest) (model.Item, error) {
	return itemService.itemRepository.Create(ctx, model.Item{
		Name:        itemRequest.Name,
		Description: itemRequest.Description,
		Price:       itemRequest.Price,
	})
}

func (itemService *ItemService) Update(ctx context.Context, id string, itemRequest dto.ItemRequest) (model.Item, error) {
	item, err := itemService.itemRepository.Update(ctx, model.Item{
		ID:          id,
		Name:        itemRequest.Name,
		Description: itemRequest.Description,
		Price:       itemRequest.Price,
	})
	if errors.Is(err, repository.ErrNotFound) {
		return model.Item{}, response.ItemNotFound
	}
	return item, err
}

func (itemService *ItemService) Delete(ctx context.Context, id string) error {
	err := itemService.itemRepository.Delete(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return response.ItemNotFound
	}
	return err
}
