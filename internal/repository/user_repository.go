package repository

import (
	"context"

	"7th-sense-server/internal/model"
)

// UserRepository is the data-access contract for users.
type UserRepository interface {
	FindAll(ctx context.Context) ([]model.User, error)
	FindByID(ctx context.Context, id string) (model.User, error)
	FindByEmail(ctx context.Context, email string) (model.User, error)
	Create(ctx context.Context, user model.User) (model.User, error)
	Update(ctx context.Context, user model.User) (model.User, error)
	Delete(ctx context.Context, id string) error
}
