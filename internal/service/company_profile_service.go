package service

import (
	"context"
	"errors"
	"net/http"

	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// CompanyProfileService holds the business rules for company profiles.
type CompanyProfileService struct {
	companyProfileRepository *repository.CompanyProfileRepository
}

// NewCompanyProfileService creates a CompanyProfileService backed by companyProfileRepository.
func NewCompanyProfileService(companyProfileRepository *repository.CompanyProfileRepository) *CompanyProfileService {
	return &CompanyProfileService{companyProfileRepository: companyProfileRepository}
}

func (companyProfileService *CompanyProfileService) GetByUserID(ctx context.Context, userID int64) (model.CompanyProfile, error) {
	profile, err := companyProfileService.companyProfileRepository.FindByUserID(ctx, userID)
	return profile, toCompanyProfileResponseError(err)
}

func (companyProfileService *CompanyProfileService) Create(ctx context.Context, profile model.CompanyProfile) (model.CompanyProfile, error) {
	createdProfile, err := companyProfileService.companyProfileRepository.Create(ctx, profile)
	return createdProfile, toCompanyProfileResponseError(err)
}

func (companyProfileService *CompanyProfileService) Update(ctx context.Context, profile model.CompanyProfile) (model.CompanyProfile, error) {
	updatedProfile, err := companyProfileService.companyProfileRepository.Update(ctx, profile)
	return updatedProfile, toCompanyProfileResponseError(err)
}

func (companyProfileService *CompanyProfileService) Delete(ctx context.Context, userID int64) error {
	return toCompanyProfileResponseError(companyProfileService.companyProfileRepository.Delete(ctx, userID))
}

// toCompanyProfileResponseError maps repository errors to profile responses. Other errors
// pass through and are sent as ServerError.
func toCompanyProfileResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Company profile not found")
	default:
		return err
	}
}
