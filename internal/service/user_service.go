package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// UserService holds the business rules for users.
type UserService struct {
	userRepository *repository.UserRepository
}

// NewUserService creates a UserService backed by userRepository.
func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (userService *UserService) GetAll(ctx context.Context) ([]model.User, error) {
	return userService.userRepository.FindAll(ctx)
}

func (userService *UserService) GetByID(ctx context.Context, id int64) (model.User, error) {
	user, err := userService.userRepository.FindByID(ctx, id)
	return user, toUserResponseError(err)
}

func (userService *UserService) Create(ctx context.Context, userRequest dto.CreateUserRequest) (model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := userService.userRepository.Create(ctx, model.User{
		Username:  userRequest.Username,
		Email:     userRequest.Email,
		Password:  string(hashedPassword),
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
		Mobile:    userRequest.Mobile,
		RoleID:    userRequest.RoleID,
		BranchID:  userRequest.BranchID,
		CompanyID: userRequest.CompanyID,
	})
	return user, toUserResponseError(err)
}

func (userService *UserService) Update(ctx context.Context, id int64, userRequest dto.UpdateUserRequest) (model.User, error) {
	user, err := userService.userRepository.Update(ctx, model.User{
		ID:        id,
		Username:  userRequest.Username,
		Email:     userRequest.Email,
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
		Mobile:    userRequest.Mobile,
		RoleID:    userRequest.RoleID,
		BranchID:  userRequest.BranchID,
		CompanyID: userRequest.CompanyID,
		Status:    userRequest.Status,
	})
	return user, toUserResponseError(err)
}

func (userService *UserService) Delete(ctx context.Context, id int64) error {
	return toUserResponseError(userService.userRepository.Delete(ctx, id))
}

// toUserResponseError maps repository errors to user responses. Other errors
// pass through and are sent as ServerError.
func toUserResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.UserNotFound
	case errors.Is(err, repository.ErrDuplicateEmail):
		return response.EmailExists
	case errors.Is(err, repository.ErrDuplicateUsername):
		return response.UsernameExists
	default:
		return err
	}
}
