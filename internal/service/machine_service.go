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
// customer (user_id); its brand, controller and model must exist, the controller must
// be of the machine's brand and the model of its controller. Its location (optional levels)
// is one chain of the customer's organization, and the machine is saved at the lowest level
// named (migration 040).
type MachineService struct {
	machineRepository *repository.MachineRepository
	userRepository    *repository.UserRepository
	modelService      *MachineModelService
	brandService      *BrandService
	controllerService *ControllerService
	accuracyService   *MachineAccuracyService
}

// NewMachineService creates a MachineService that stores machines with machineRepository
// and checks their customer, model, brand, controller and accuracy with the other
// arguments.
func NewMachineService(
	machineRepository *repository.MachineRepository,
	userRepository *repository.UserRepository,
	modelService *MachineModelService,
	brandService *BrandService,
	controllerService *ControllerService,
	accuracyService *MachineAccuracyService,
) *MachineService {
	return &MachineService{
		machineRepository: machineRepository,
		userRepository:    userRepository,
		modelService:      modelService,
		brandService:      brandService,
		controllerService: controllerService,
		accuracyService:   accuracyService,
	}
}

// GetAll returns the machines actor may access (see machineScope). A non-nil
// customerID keeps only that customer's machines.
func (machineService *MachineService) GetAll(ctx context.Context, actor model.AuthUser, customerID *int64) ([]model.Machine, error) {
	scope, err := machineScope(actor, customerID)
	if err != nil {
		return nil, err
	}
	machines, err := machineService.machineRepository.FindAll(ctx, scope)
	if err != nil {
		return nil, err
	}
	for index := range machines {
		machines[index].CustomerFilters = model.SortOrganizationLevels(machines[index].CustomerFilters)
	}
	return machines, nil
}

// GetLocationFilters returns the organization levels a customer's machines are placed at (its
// filters, in hierarchy order): the Add / Edit Machine form shows exactly these. The superadmin
// names the customer (customerID, required); everyone else gets their own company's, and naming
// another customer is Forbidden.
func (machineService *MachineService) GetLocationFilters(ctx context.Context, actor model.AuthUser, customerID *int64) ([]string, error) {
	company := actor.Company()
	if actor.Role == model.RoleSuperAdmin {
		company = customerID
	} else if customerID != nil && (company == nil || *customerID != *company) {
		return nil, response.Forbidden
	}
	if company == nil {
		return nil, response.NewError(http.StatusBadRequest, "Please select customer")
	}
	isCustomer, err := machineService.userRepository.IsCustomer(ctx, *company)
	if err != nil {
		return nil, err
	}
	if !isCustomer {
		return nil, response.NewError(http.StatusBadRequest, "Please select valid customer")
	}
	filters, err := machineService.machineRepository.FindCustomerFilters(ctx, *company)
	if err != nil {
		return nil, err
	}
	return model.SortOrganizationLevels(filters), nil
}

// GetByID returns the machine when actor may access it (see machineScope), else
// "Machine not found". A non-nil customerID finds it only when that customer owns it.
func (machineService *MachineService) GetByID(ctx context.Context, actor model.AuthUser, customerID *int64, id int64) (model.Machine, error) {
	scope, err := machineScope(actor, customerID)
	if err != nil {
		return model.Machine{}, err
	}
	return machineService.find(ctx, scope, id)
}

// find returns machine id when it is in scope.
func (machineService *MachineService) find(ctx context.Context, scope model.MachineScope, id int64) (model.Machine, error) {
	machine, err := machineService.machineRepository.FindByID(ctx, scope, id)
	machine.CustomerFilters = model.SortOrganizationLevels(machine.CustomerFilters)
	return machine, toMachineResponseError(err)
}

// CheckOwnedBy returns "Please select valid machines" unless every one of ids is
// a machine of customer ownerID.
func (machineService *MachineService) CheckOwnedBy(ctx context.Context, ownerID int64, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	count, err := machineService.machineRepository.CountOwnedBy(ctx, ownerID, ids)
	if err != nil {
		return err
	}
	if count != len(ids) {
		return response.NewError(http.StatusBadRequest, "Please select valid machines")
	}
	return nil
}

// ownCompanyMachine makes the machine of fields belong to actor's company when
// actor is not the superadmin. Only managers reach here (managersOnly); one without
// a company gets Forbidden.
func ownCompanyMachine(actor model.AuthUser, fields *dto.MachineFields) error {
	if actor.Role == model.RoleSuperAdmin {
		return nil
	}
	company := actor.Company()
	if company == nil {
		return response.Forbidden
	}
	fields.UserID = *company
	return nil
}

// machineScope returns the machines actor may read. The superadmin reads every
// machine (customerID narrows to one customer). A customer and its admins read
// every machine of their company; other users only the machines assigned to them
// (user_machine_access) and, with all machines on, the ones under their machine filter
// (user_machine_filters). Asking for another customer's machines is Forbidden.
func machineScope(actor model.AuthUser, customerID *int64) (model.MachineScope, error) {
	if actor.Role == model.RoleSuperAdmin {
		return model.MachineScope{OwnerID: customerID}, nil
	}
	company := actor.Company()
	if customerID != nil && (company == nil || *customerID != *company) {
		return model.MachineScope{}, response.Forbidden
	}
	if company != nil && model.SeesAllCompanyMachines(actor.Role) {
		return model.MachineScope{OwnerID: company}, nil
	}
	return model.MachineScope{OwnerID: company, AssignedUserID: &actor.ID}, nil
}

// Create adds a machine. A customer or admin (actor) adds it to their own company,
// whatever customer the request names; the superadmin to the customer named.
func (machineService *MachineService) Create(ctx context.Context, actor model.AuthUser, machineRequest dto.CreateMachineRequest) (model.Machine, error) {
	if err := ownCompanyMachine(actor, &machineRequest.MachineFields); err != nil {
		return model.Machine{}, err
	}
	if err := machineService.checkReferences(ctx, machineRequest.MachineFields); err != nil {
		return model.Machine{}, err
	}
	machine := machineFromFields(machineRequest.MachineFields)
	machine.Password = machineRequest.Password
	machine.Status = dto.MachineStatusActive

	id, err := machineService.machineRepository.Create(ctx, machine, machineRequest.Location.Deepest())
	if err != nil {
		return model.Machine{}, toMachineResponseError(err)
	}
	return machineService.find(ctx, model.MachineScope{}, id)
}

// Update replaces the form fields. A blank password or status keeps the saved one.
// A customer or admin (actor) changes only their company's machines ("Machine not
// found" for others) and cannot move one to another company.
func (machineService *MachineService) Update(ctx context.Context, actor model.AuthUser, id int64, machineRequest dto.UpdateMachineRequest) (model.Machine, error) {
	if err := ownCompanyMachine(actor, &machineRequest.MachineFields); err != nil {
		return model.Machine{}, err
	}
	if actor.Role != model.RoleSuperAdmin {
		if _, err := machineService.find(ctx, model.MachineScope{OwnerID: actor.Company()}, id); err != nil {
			return model.Machine{}, err
		}
	}
	if err := machineService.checkReferences(ctx, machineRequest.MachineFields); err != nil {
		return model.Machine{}, err
	}
	machine := machineFromFields(machineRequest.MachineFields)
	machine.ID = id
	machine.Password = machineRequest.Password
	machine.Status = machineRequest.Status

	if err := machineService.machineRepository.Update(ctx, machine, machineRequest.Location.Deepest()); err != nil {
		return model.Machine{}, toMachineResponseError(err)
	}
	return machineService.find(ctx, model.MachineScope{}, id)
}

func (machineService *MachineService) Delete(ctx context.Context, id int64) error {
	return toMachineResponseError(machineService.machineRepository.Delete(ctx, id))
}

// checkReferences returns a "Please select valid …" error unless the customer, brand,
// controller, model, location and accuracy of fields exist and fit together.
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

	if err := machineService.checkLocation(ctx, fields.UserID, fields.Location); err != nil {
		return err
	}

	if fields.AccuracyID != nil {
		if err := machineService.accuracyService.checkExists(ctx, *fields.AccuracyID, response.SelectValidMachineAccuracy); err != nil {
			return err
		}
	}
	return nil
}

// checkLocation returns "Please select valid location" unless the levels location names are one
// chain of the customer's organization. Every level is optional (none = no location), and levels
// the customer no longer uses may be named too: a machine keeps the place it was saved at when the
// customer's filters change.
func (machineService *MachineService) checkLocation(ctx context.Context, customerID int64, location model.OrganizationPath) error {
	if location.IsEmpty() {
		return nil
	}
	inOrganization, err := machineService.machineRepository.IsPathInOrganization(ctx, customerID, location)
	if err != nil {
		return err
	}
	if !inOrganization {
		return response.NewError(http.StatusBadRequest, "Please select valid location")
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
