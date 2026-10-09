package repository

import (
	"context"

	"7th-sense-server/internal/model"
)

// ItemRepository is the data-access contract for items.
type ItemRepository interface {
	FindAll(ctx context.Context) ([]model.Item, error)
	FindByID(ctx context.Context, id string) (model.Item, error)
	Create(ctx context.Context, item model.Item) (model.Item, error)
	Update(ctx context.Context, item model.Item) (model.Item, error)
	Delete(ctx context.Context, id string) error
}
