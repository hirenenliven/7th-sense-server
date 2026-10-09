package repository

import (
	"context"
	"slices"
	"strconv"
	"sync"
	"time"

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

type inMemoryItemRepository struct {
	mu     sync.RWMutex
	items  map[string]model.Item
	nextID int
}

// NewInMemoryItemRepository returns an ItemRepository that keeps data in memory.
func NewInMemoryItemRepository() ItemRepository {
	return &inMemoryItemRepository{items: make(map[string]model.Item), nextID: 1}
}

func (r *inMemoryItemRepository) FindAll(ctx context.Context) ([]model.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]model.Item, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}
	slices.SortFunc(items, func(a, b model.Item) int { return a.CreatedAt.Compare(b.CreatedAt) })
	return items, nil
}

func (r *inMemoryItemRepository) FindByID(ctx context.Context, id string) (model.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[id]
	if !ok {
		return model.Item{}, ErrNotFound
	}
	return item, nil
}

func (r *inMemoryItemRepository) Create(ctx context.Context, item model.Item) (model.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	item.ID = strconv.Itoa(r.nextID)
	item.CreatedAt = now
	item.UpdatedAt = now
	r.nextID++
	r.items[item.ID] = item
	return item, nil
}

func (r *inMemoryItemRepository) Update(ctx context.Context, item model.Item) (model.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.items[item.ID]
	if !ok {
		return model.Item{}, ErrNotFound
	}
	item.CreatedAt = existing.CreatedAt
	item.UpdatedAt = time.Now().UTC()
	r.items[item.ID] = item
	return item, nil
}

func (r *inMemoryItemRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return ErrNotFound
	}
	delete(r.items, id)
	return nil
}
