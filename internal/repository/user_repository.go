package repository

import (
	"context"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"7th-sense-server/internal/model"
)

// UserRepository is the data-access contract for users. Swap the in-memory
// implementation for a database one without touching services.
type UserRepository interface {
	FindAll(ctx context.Context) ([]model.User, error)
	FindByID(ctx context.Context, id string) (model.User, error)
	FindByEmail(ctx context.Context, email string) (model.User, error)
	Create(ctx context.Context, user model.User) (model.User, error)
	Update(ctx context.Context, user model.User) (model.User, error)
	Delete(ctx context.Context, id string) error
}

type inMemoryUserRepository struct {
	mu     sync.RWMutex
	users  map[string]model.User
	nextID int
}

// NewInMemoryUserRepository returns a UserRepository that keeps data in memory.
func NewInMemoryUserRepository() UserRepository {
	return &inMemoryUserRepository{users: make(map[string]model.User), nextID: 1}
}

func (r *inMemoryUserRepository) FindAll(ctx context.Context) ([]model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	slices.SortFunc(users, func(a, b model.User) int { return a.CreatedAt.Compare(b.CreatedAt) })
	return users, nil
}

func (r *inMemoryUserRepository) FindByID(ctx context.Context, id string) (model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return model.User{}, ErrNotFound
	}
	return user, nil
}

func (r *inMemoryUserRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if strings.EqualFold(user.Email, email) {
			return user, nil
		}
	}
	return model.User{}, ErrNotFound
}

func (r *inMemoryUserRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	user.ID = strconv.Itoa(r.nextID)
	user.CreatedAt = now
	user.UpdatedAt = now
	r.nextID++
	r.users[user.ID] = user
	return user, nil
}

func (r *inMemoryUserRepository) Update(ctx context.Context, user model.User) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.users[user.ID]
	if !ok {
		return model.User{}, ErrNotFound
	}
	user.CreatedAt = existing.CreatedAt
	user.UpdatedAt = time.Now().UTC()
	r.users[user.ID] = user
	return user, nil
}

func (r *inMemoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return ErrNotFound
	}
	delete(r.users, id)
	return nil
}
