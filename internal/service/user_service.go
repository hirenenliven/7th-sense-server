package service

import (
	"context"
	"errors"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// UserService holds the business rules for users.
type UserService struct {
	repo repository.UserRepository
}

// NewUserService creates a UserService backed by repo.
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAll(ctx context.Context) ([]model.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *UserService) GetByID(ctx context.Context, id string) (model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return model.User{}, response.ErrUserNotFound
	}
	return user, err
}

func (s *UserService) Create(ctx context.Context, req dto.UserRequest) (model.User, error) {
	if err := s.ensureEmailAvailable(ctx, req.Email, ""); err != nil {
		return model.User{}, err
	}
	return s.repo.Create(ctx, model.User{Name: req.Name, Email: req.Email})
}

func (s *UserService) Update(ctx context.Context, id string, req dto.UserRequest) (model.User, error) {
	if _, err := s.GetByID(ctx, id); err != nil {
		return model.User{}, err
	}
	if err := s.ensureEmailAvailable(ctx, req.Email, id); err != nil {
		return model.User{}, err
	}

	user, err := s.repo.Update(ctx, model.User{ID: id, Name: req.Name, Email: req.Email})
	if errors.Is(err, repository.ErrNotFound) {
		return model.User{}, response.ErrUserNotFound
	}
	return user, err
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return response.ErrUserNotFound
	}
	return err
}

// ensureEmailAvailable fails if email belongs to a user other than exceptID.
func (s *UserService) ensureEmailAvailable(ctx context.Context, email, exceptID string) error {
	existing, err := s.repo.FindByEmail(ctx, email)
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return nil
	case err != nil:
		return err
	case existing.ID != exceptID:
		return response.ErrEmailExists
	default:
		return nil
	}
}
