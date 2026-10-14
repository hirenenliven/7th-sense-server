package service

import (
	"context"
	"errors"
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// ControllerService holds the business rules for controllers. A controller belongs to
// a brand; its machine models (models table) belong to it.
type ControllerService struct {
	controllerRepository *repository.ControllerRepository
	brandService         *BrandService
}

// NewControllerService creates a ControllerService backed by controllerRepository that
// checks the selected brand with brandService.
func NewControllerService(controllerRepository *repository.ControllerRepository, brandService *BrandService) *ControllerService {
	return &ControllerService{controllerRepository: controllerRepository, brandService: brandService}
}

// GetAll returns the controllers. A non-nil brandID keeps only that brand's controllers.
func (controllerService *ControllerService) GetAll(ctx context.Context, brandID *int64) ([]model.Controller, error) {
	return controllerService.controllerRepository.FindAll(ctx, brandID)
}

func (controllerService *ControllerService) GetByID(ctx context.Context, id int64) (model.Controller, error) {
	controller, err := controllerService.controllerRepository.FindByID(ctx, id)
	return controller, toControllerResponseError(err)
}

func (controllerService *ControllerService) Create(ctx context.Context, controllerRequest dto.CreateControllerRequest) (model.Controller, error) {
	if err := controllerService.checkBrand(ctx, *controllerRequest.BrandID); err != nil {
		return model.Controller{}, err
	}
	controller, err := controllerService.controllerRepository.Create(ctx, model.Controller{
		Name:        controllerRequest.Name,
		BrandID:     controllerRequest.BrandID,
		Description: controllerRequest.Description,
	})
	return controller, toControllerResponseError(err)
}

// Update replaces the controller's fields. Its brand cannot change while machines use
// it, because each machine stores the brand of its controller.
func (controllerService *ControllerService) Update(ctx context.Context, id int64, controllerRequest dto.UpdateControllerRequest) (model.Controller, error) {
	saved, err := controllerService.GetByID(ctx, id)
	if err != nil {
		return model.Controller{}, err
	}
	if err := controllerService.checkBrand(ctx, *controllerRequest.BrandID); err != nil {
		return model.Controller{}, err
	}
	if saved.BrandID == nil || *saved.BrandID != *controllerRequest.BrandID {
		hasMachines, err := controllerService.controllerRepository.HasMachines(ctx, id)
		if err != nil {
			return model.Controller{}, err
		}
		if hasMachines {
			return model.Controller{}, response.NewError(http.StatusBadRequest, "Machines use this controller, so its brand cannot be changed")
		}
	}

	controller, err := controllerService.controllerRepository.Update(ctx, model.Controller{
		ID:          id,
		Name:        controllerRequest.Name,
		BrandID:     controllerRequest.BrandID,
		Description: controllerRequest.Description,
	})
	return controller, toControllerResponseError(err)
}

// Delete soft-deletes the controller. A controller that still has models or is used
// by machines is kept: those must be deleted or changed first.
func (controllerService *ControllerService) Delete(ctx context.Context, id int64) error {
	hasModels, err := controllerService.controllerRepository.HasModels(ctx, id)
	if err != nil {
		return err
	}
	if hasModels {
		return response.NewError(http.StatusBadRequest, "This controller has models. Delete its models first")
	}
	hasMachines, err := controllerService.controllerRepository.HasMachines(ctx, id)
	if err != nil {
		return err
	}
	if hasMachines {
		return response.NewError(http.StatusBadRequest, "Machines use this controller. Change or delete those machines first")
	}
	return toControllerResponseError(controllerService.controllerRepository.Delete(ctx, id))
}

// checkBrand returns "Please enter valid brand" unless brandID is an active brand.
func (controllerService *ControllerService) checkBrand(ctx context.Context, brandID int64) error {
	if _, err := controllerService.brandService.GetByID(ctx, brandID); err != nil {
		return asInvalidSelection(err, "Brand not found", response.ValidBrand.Message)
	}
	return nil
}

// toControllerResponseError maps repository errors to controller responses. Other errors
// pass through and are sent as ServerError.
func toControllerResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Controller not found")
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NewError(http.StatusBadRequest, "This brand already has a controller with this name")
	case errors.Is(err, repository.ErrInvalidBrand):
		return response.ValidBrand
	default:
		return err
	}
}
