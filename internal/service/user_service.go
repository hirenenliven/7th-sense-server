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
	userRepository repository.UserRepository
}

// NewUserService creates a UserService backed by userRepository.
func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (userService *UserService) GetAll(ctx context.Context) ([]model.User, error) {
	return userService.userRepository.FindAll(ctx)
}

func (userService *UserService) GetByID(ctx context.Context, id string) (model.User, error) {
	user, err := userService.userRepository.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return model.User{}, response.UserNotFound
	}
	return user, err
}

func (userService *UserService) Create(ctx context.Context, userRequest dto.UserRequest) (model.User, error) {
	if err := userService.ensureEmailAvailable(ctx, userRequest.Email, ""); err != nil {
		return model.User{}, err
	}
	return userService.userRepository.Create(ctx, model.User{Name: userRequest.Name, Email: userRequest.Email})
}

func (userService *UserService) Update(ctx context.Context, id string, userRequest dto.UserRequest) (model.User, error) {
	if _, err := userService.GetByID(ctx, id); err != nil {
		return model.User{}, err
	}
	if err := userService.ensureEmailAvailable(ctx, userRequest.Email, id); err != nil {
		return model.User{}, err
	}

	user, err := userService.userRepository.Update(ctx, model.User{ID: id, Name: userRequest.Name, Email: userRequest.Email})
	if errors.Is(err, repository.ErrNotFound) {
		return model.User{}, response.UserNotFound
	}
	return user, err
}

func (userService *UserService) Delete(ctx context.Context, id string) error {
	err := userService.userRepository.Delete(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return response.UserNotFound
	}
	return err
}

// ensureEmailAvailable fails if email belongs to a user other than exceptID.
func (userService *UserService) ensureEmailAvailable(ctx context.Context, email, exceptID string) error {
	existing, err := userService.userRepository.FindByEmail(ctx, email)
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return nil
	case err != nil:
		return err
	case existing.ID != exceptID:
		return response.EmailExists
	default:
		return nil
	}
}
