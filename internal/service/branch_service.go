package service

import (
	"context"
	"errors"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// BranchService holds the business rules for branches.
type BranchService struct {
	branchRepository *repository.BranchRepository
}

// NewBranchService creates a BranchService backed by branchRepository.
func NewBranchService(branchRepository *repository.BranchRepository) *BranchService {
	return &BranchService{branchRepository: branchRepository}
}

func (branchService *BranchService) GetAll(ctx context.Context) ([]model.Branch, error) {
	return branchService.branchRepository.FindAll(ctx)
}

func (branchService *BranchService) GetByID(ctx context.Context, id int64) (model.Branch, error) {
	branch, err := branchService.branchRepository.FindByID(ctx, id)
	return branch, toBranchResponseError(err)
}

func (branchService *BranchService) Create(ctx context.Context, branchRequest dto.CreateBranchRequest) (model.Branch, error) {
	branch, err := branchService.branchRepository.Create(ctx, model.Branch{
		Code:      branchRequest.Code,
		Name:      branchRequest.Name,
		Email:     branchRequest.Email,
		Address:   branchRequest.Address,
		City:      branchRequest.City,
		StateID:   branchRequest.StateID,
		CustomerID: branchRequest.CustomerID,
		GSTNumber: branchRequest.GSTNumber,
		PANNumber: branchRequest.PANNumber,
		ColorCode: branchRequest.ColorCode,
		Status:    branchRequest.Status,
	})
	return branch, toBranchResponseError(err)
}

func (branchService *BranchService) Update(ctx context.Context, id int64, branchRequest dto.UpdateBranchRequest) (model.Branch, error) {
	branch, err := branchService.branchRepository.Update(ctx, model.Branch{
		ID:        id,
		Code:      branchRequest.Code,
		Name:      branchRequest.Name,
		Email:     branchRequest.Email,
		Address:   branchRequest.Address,
		City:      branchRequest.City,
		StateID:   branchRequest.StateID,
		CustomerID: branchRequest.CustomerID,
		GSTNumber: branchRequest.GSTNumber,
		PANNumber: branchRequest.PANNumber,
		ColorCode: branchRequest.ColorCode,
		Status:    branchRequest.Status,
	})
	return branch, toBranchResponseError(err)
}

func (branchService *BranchService) Delete(ctx context.Context, id int64) error {
	return toBranchResponseError(branchService.branchRepository.Delete(ctx, id))
}

// toBranchResponseError maps repository errors to branch responses. Other errors
// pass through and are sent as ServerError.
func toBranchResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.BranchNotFound
	case errors.Is(err, repository.ErrDuplicateCode):
		return response.CodeExists
	default:
		return err
	}
}
