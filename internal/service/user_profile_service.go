package service

import (
	"context"
	"errors"

	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// UserProfileService holds the business rules for user profiles.
type UserProfileService struct {
	userProfileRepository *repository.UserProfileRepository
}

// NewUserProfileService creates a UserProfileService backed by userProfileRepository.
func NewUserProfileService(userProfileRepository *repository.UserProfileRepository) *UserProfileService {
	return &UserProfileService{userProfileRepository: userProfileRepository}
}

func (userProfileService *UserProfileService) GetByUserID(ctx context.Context, userID int64) (model.UserProfile, error) {
	profile, err := userProfileService.userProfileRepository.FindByUserID(ctx, userID)
	return profile, toUserProfileResponseError(err)
}

func (userProfileService *UserProfileService) Create(ctx context.Context, profile model.UserProfile) (model.UserProfile, error) {
	createdProfile, err := userProfileService.userProfileRepository.Create(ctx, profile)
	return createdProfile, toUserProfileResponseError(err)
}

func (userProfileService *UserProfileService) Update(ctx context.Context, profile model.UserProfile) (model.UserProfile, error) {
	updatedProfile, err := userProfileService.userProfileRepository.Update(ctx, profile)
	return updatedProfile, toUserProfileResponseError(err)
}

func (userProfileService *UserProfileService) Delete(ctx context.Context, userID int64) error {
	return toUserProfileResponseError(userProfileService.userProfileRepository.Delete(ctx, userID))
}

// toUserProfileResponseError maps repository errors to profile responses. Other errors
// pass through and are sent as ServerError.
func toUserProfileResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.UserProfileNotFound
	default:
		return err
	}
}
