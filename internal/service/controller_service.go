package service

import (
	"context"
	"errors"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// ControllerService holds the business rules for controllers.
type ControllerService struct {
	controllerRepository *repository.ControllerRepository
}

// NewControllerService creates a ControllerService backed by controllerRepository.
func NewControllerService(controllerRepository *repository.ControllerRepository) *ControllerService {
	return &ControllerService{controllerRepository: controllerRepository}
}

func (controllerService *ControllerService) GetAll(ctx context.Context) ([]model.Controller, error) {
	return controllerService.controllerRepository.FindAll(ctx)
}

func (controllerService *ControllerService) GetByID(ctx context.Context, id int64) (model.Controller, error) {
	controller, err := controllerService.controllerRepository.FindByID(ctx, id)
	return controller, toControllerResponseError(err)
}

func (controllerService *ControllerService) Create(ctx context.Context, controllerRequest dto.CreateControllerRequest) (model.Controller, error) {
	controller, err := controllerService.controllerRepository.Create(ctx, model.Controller{
		Name:        controllerRequest.Name,
		Code:        controllerRequest.Code,
		Variant:     controllerRequest.Variant,
		BrandID:     controllerRequest.BrandID,
		Description: controllerRequest.Description,
	})
	return controller, toControllerResponseError(err)
}

func (controllerService *ControllerService) Update(ctx context.Context, id int64, controllerRequest dto.UpdateControllerRequest) (model.Controller, error) {
	controller, err := controllerService.controllerRepository.Update(ctx, model.Controller{
		ID:          id,
		Name:        controllerRequest.Name,
		Code:        controllerRequest.Code,
		Variant:     controllerRequest.Variant,
		BrandID:     controllerRequest.BrandID,
		Description: controllerRequest.Description,
	})
	return controller, toControllerResponseError(err)
}

func (controllerService *ControllerService) Delete(ctx context.Context, id int64) error {
	return toControllerResponseError(controllerService.controllerRepository.Delete(ctx, id))
}

// toControllerResponseError maps repository errors to controller responses. Other errors
// pass through and are sent as ServerError.
func toControllerResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.ControllerNotFound
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NameExists
	case errors.Is(err, repository.ErrDuplicateCode):
		return response.CodeExists
	case errors.Is(err, repository.ErrInvalidBrand):
		return response.ValidBrand
	default:
		return err
	}
}
