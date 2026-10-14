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

// MachineService holds the business rules for machines. Every machine belongs to a
// customer (user_id); its brand, controller, model and line must exist, the controller
// must be of the machine's brand, the model of its controller and the line one of the
// customer's lines.
type MachineService struct {
	machineRepository *repository.MachineRepository
	userRepository    *repository.UserRepository
	modelService      *MachineModelService
	brandService      *BrandService
	controllerService *ControllerService
	userLineService   *UserLineService
}

// NewMachineService creates a MachineService that stores machines with machineRepository
// and checks their customer, model, brand, controller and line with the other arguments.
func NewMachineService(
	machineRepository *repository.MachineRepository,
	userRepository *repository.UserRepository,
	modelService *MachineModelService,
	brandService *BrandService,
	controllerService *ControllerService,
	userLineService *UserLineService,
) *MachineService {
	return &MachineService{
		machineRepository: machineRepository,
		userRepository:    userRepository,
		modelService:      modelService,
		brandService:      brandService,
		controllerService: controllerService,
		userLineService:   userLineService,
	}
}

// GetAll returns the machines. A non-nil customerID keeps only that customer's machines.
func (machineService *MachineService) GetAll(ctx context.Context, customerID *int64) ([]model.Machine, error) {
	return machineService.machineRepository.FindAll(ctx, customerID)
}

// GetByID returns the machine. A non-nil customerID finds it only when that customer owns it.
func (machineService *MachineService) GetByID(ctx context.Context, customerID *int64, id int64) (model.Machine, error) {
	machine, err := machineService.machineRepository.FindByID(ctx, customerID, id)
	return machine, toMachineResponseError(err)
}

func (machineService *MachineService) Create(ctx context.Context, machineRequest dto.CreateMachineRequest) (model.Machine, error) {
	if err := machineService.checkReferences(ctx, machineRequest.MachineFields); err != nil {
		return model.Machine{}, err
	}
	machine := machineFromFields(machineRequest.MachineFields)
	machine.Password = machineRequest.Password
	machine.Status = dto.MachineStatusActive

	id, err := machineService.machineRepository.Create(ctx, machine)
	if err != nil {
		return model.Machine{}, toMachineResponseError(err)
	}
	return machineService.GetByID(ctx, nil, id)
}

// Update replaces the form fields. A blank password or status keeps the saved one.
func (machineService *MachineService) Update(ctx context.Context, id int64, machineRequest dto.UpdateMachineRequest) (model.Machine, error) {
	if err := machineService.checkReferences(ctx, machineRequest.MachineFields); err != nil {
		return model.Machine{}, err
	}
	machine := machineFromFields(machineRequest.MachineFields)
	machine.ID = id
	machine.Password = machineRequest.Password
	machine.Status = machineRequest.Status

	if err := machineService.machineRepository.Update(ctx, machine); err != nil {
		return model.Machine{}, toMachineResponseError(err)
	}
	return machineService.GetByID(ctx, nil, id)
}

func (machineService *MachineService) Delete(ctx context.Context, id int64) error {
	return toMachineResponseError(machineService.machineRepository.Delete(ctx, id))
}

// checkReferences returns a "Please select valid …" error unless the customer, brand,
// controller, model and line of fields exist and fit together.
func (machineService *MachineService) checkReferences(ctx context.Context, fields dto.MachineFields) error {
	isCustomer, err := machineService.userRepository.IsCustomer(ctx, fields.UserID)
	if err != nil {
		return err
	}
	if !isCustomer {
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	}

	if _, err := machineService.brandService.GetByID(ctx, fields.BrandID); err != nil {
		return asInvalidSelection(err, "Brand not found", "Please select valid brand")
	}

	controller, err := machineService.controllerService.GetByID(ctx, fields.ControllerID)
	if err != nil {
		return asInvalidSelection(err, "Controller not found", "Please select valid controller")
	}
	if controller.BrandID == nil || *controller.BrandID != fields.BrandID {
		return response.NewError(http.StatusBadRequest, "Please select a controller of the selected brand")
	}

	machineModel, err := machineService.modelService.GetByID(ctx, fields.ModelID)
	if err != nil {
		return asInvalidSelection(err, "Model not found", "Please select valid model")
	}
	if machineModel.ControllerID != fields.ControllerID {
		return response.NewError(http.StatusBadRequest, "Please select a model of the selected controller")
	}

	if fields.LineID != nil {
		if _, err := machineService.userLineService.GetByID(ctx, fields.UserID, *fields.LineID); err != nil {
			return asInvalidSelection(err, "Line not found", "Please select valid line")
		}
	}
	return nil
}

// asInvalidSelection turns the not-found response of a lookup into a 400 with message;
// other errors pass through.
func asInvalidSelection(err error, notFoundMessage, message string) error {
	if err.Error() == notFoundMessage {
		return response.NewError(http.StatusBadRequest, message)
	}
	return err
}

// machineFromFields copies the form fields into a machine; blank optional dates become NULL.
func machineFromFields(fields dto.MachineFields) model.Machine {
	controllerVersionType := fields.ControllerVersionType
	return model.Machine{
		UserID:                &fields.UserID,
		ModelID:               &fields.ModelID,
		Name:                  fields.Name,
		SerialNumber:          fields.SerialNumber,
		BrandID:               &fields.BrandID,
		ControllerID:          &fields.ControllerID,
		ControllerVersion:     fields.ControllerVersion,
		AccuracyID:            fields.AccuracyID,
		LineID:                fields.LineID,
		IsInsured:             fields.IsInsured,
		MaintenanceRequired:   fields.MaintenanceRequired,
		InstallationDate:      fields.InstallationDate,
		BillDate:              optionalText(fields.BillDate),
		WarrantyDate:          optionalText(fields.WarrantyDate),
		HourlyCost:            fields.HourlyCost,
		AMC:                   fields.AMC,
		MachineSystemType:     fields.MachineSystemType,
		SecurityMode:          fields.SecurityMode,
		LocalIP:               fields.LocalIP,
		LiveIP:                fields.LiveIP,
		Port:                  fields.Port,
		Username:              fields.Username,
		ControllerVersionType: &controllerVersionType,
		Longitude:             fields.Longitude,
		Latitude:              fields.Latitude,
		HealthIDs:             fields.HealthIDs,
	}
}

// optionalText returns nil for "" so the column is stored as NULL.
func optionalText(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// toMachineResponseError maps repository errors to machine responses. Other errors
// pass through and are sent as ServerError.
func toMachineResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Machine not found")
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NewError(http.StatusBadRequest, "This customer already has a machine with this name")
	case errors.Is(err, repository.ErrInvalidUser):
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	case errors.Is(err, repository.ErrInvalidModel):
		return response.NewError(http.StatusBadRequest, "Please select valid model")
	case errors.Is(err, repository.ErrInvalidBrand):
		return response.NewError(http.StatusBadRequest, "Please select valid brand")
	case errors.Is(err, repository.ErrInvalidController):
		return response.NewError(http.StatusBadRequest, "Please select valid controller")
	default:
		return err
	}
}
