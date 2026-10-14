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

// MachineModelService holds the business rules for machine models (models table).
// Every model belongs to an active controller; its name is unique per controller.
type MachineModelService struct {
	machineModelRepository *repository.MachineModelRepository
	controllerService      *ControllerService
}

// NewMachineModelService creates a MachineModelService backed by machineModelRepository
// that checks the selected controller with controllerService.
func NewMachineModelService(machineModelRepository *repository.MachineModelRepository, controllerService *ControllerService) *MachineModelService {
	return &MachineModelService{machineModelRepository: machineModelRepository, controllerService: controllerService}
}

// GetAll returns the models. A non-nil controllerID keeps only that controller's models.
func (machineModelService *MachineModelService) GetAll(ctx context.Context, controllerID *int64) ([]model.MachineModel, error) {
	return machineModelService.machineModelRepository.FindAll(ctx, controllerID)
}

func (machineModelService *MachineModelService) GetByID(ctx context.Context, id int64) (model.MachineModel, error) {
	machineModel, err := machineModelService.machineModelRepository.FindByID(ctx, id)
	return machineModel, toMachineModelResponseError(err)
}

func (machineModelService *MachineModelService) Create(ctx context.Context, modelRequest dto.MachineModelRequest) (model.MachineModel, error) {
	if err := machineModelService.checkController(ctx, modelRequest.ControllerID); err != nil {
		return model.MachineModel{}, err
	}
	machineModel, err := machineModelService.machineModelRepository.Create(ctx, model.MachineModel{
		ControllerID: modelRequest.ControllerID,
		Name:         modelRequest.Name,
	})
	return machineModel, toMachineModelResponseError(err)
}

// Update replaces the model's controller and name. Its controller cannot change while
// machines use it, because each machine stores the controller of its model.
func (machineModelService *MachineModelService) Update(ctx context.Context, id int64, modelRequest dto.MachineModelRequest) (model.MachineModel, error) {
	saved, err := machineModelService.GetByID(ctx, id)
	if err != nil {
		return model.MachineModel{}, err
	}
	if err := machineModelService.checkController(ctx, modelRequest.ControllerID); err != nil {
		return model.MachineModel{}, err
	}
	if saved.ControllerID != modelRequest.ControllerID {
		hasMachines, err := machineModelService.machineModelRepository.HasMachines(ctx, id)
		if err != nil {
			return model.MachineModel{}, err
		}
		if hasMachines {
			return model.MachineModel{}, response.NewError(http.StatusBadRequest, "Machines use this model, so its controller cannot be changed")
		}
	}

	machineModel, err := machineModelService.machineModelRepository.Update(ctx, model.MachineModel{
		ID:           id,
		ControllerID: modelRequest.ControllerID,
		Name:         modelRequest.Name,
	})
	return machineModel, toMachineModelResponseError(err)
}

// Delete soft-deletes the model. A model used by machines is kept.
func (machineModelService *MachineModelService) Delete(ctx context.Context, id int64) error {
	hasMachines, err := machineModelService.machineModelRepository.HasMachines(ctx, id)
	if err != nil {
		return err
	}
	if hasMachines {
		return response.NewError(http.StatusBadRequest, "Machines use this model. Change or delete those machines first")
	}
	return toMachineModelResponseError(machineModelService.machineModelRepository.Delete(ctx, id))
}

// checkController returns "Please select valid controller" unless controllerID is an
// active controller.
func (machineModelService *MachineModelService) checkController(ctx context.Context, controllerID int64) error {
	if _, err := machineModelService.controllerService.GetByID(ctx, controllerID); err != nil {
		return asInvalidSelection(err, "Controller not found", "Please select valid controller")
	}
	return nil
}

// toMachineModelResponseError maps repository errors to model responses. Other errors
// pass through and are sent as ServerError.
func toMachineModelResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Model not found")
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NewError(http.StatusBadRequest, "This controller already has a model with this name")
	case errors.Is(err, repository.ErrInvalidController):
		return response.NewError(http.StatusBadRequest, "Please select valid controller")
	default:
		return err
	}
}
