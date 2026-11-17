package service

import (
	"context"
	"errors"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// OperatorTaskService holds the rules of operator tasks (server ADR 0015): which maintenance
// tasks apply to a customer's machines of a model (the Machine Type; none = every machine).
// Only the superadmin adds, edits and deletes them (superAdminOnly routes); everyone reads their
// company's (companyOwner), so a future operator checklist can show them. A task is added once per
// customer and machine type.
type OperatorTaskService struct {
	operatorTaskRepository *repository.OperatorTaskRepository
	userRepository         *repository.UserRepository
	machineModelService    *MachineModelService
	maintenanceTaskService *MaintenanceTaskService
}

// NewOperatorTaskService creates an OperatorTaskService that stores operator tasks with
// operatorTaskRepository and checks customers, models and tasks with the others.
func NewOperatorTaskService(operatorTaskRepository *repository.OperatorTaskRepository, userRepository *repository.UserRepository,
	machineModelService *MachineModelService, maintenanceTaskService *MaintenanceTaskService) *OperatorTaskService {
	return &OperatorTaskService{
		operatorTaskRepository: operatorTaskRepository,
		userRepository:         userRepository,
		machineModelService:    machineModelService,
		maintenanceTaskService: maintenanceTaskService,
	}
}

// GetAll returns the operator tasks actor may read (see companyOwner). A non-nil customerID keeps
// only that customer's.
func (operatorTaskService *OperatorTaskService) GetAll(ctx context.Context, actor model.AuthUser, customerID *int64) ([]model.OperatorTask, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return nil, err
	}
	return operatorTaskService.operatorTaskRepository.FindAll(ctx, ownerID)
}

// GetByID returns the operator task when actor may read it, else OperatorTaskNotFound.
func (operatorTaskService *OperatorTaskService) GetByID(ctx context.Context, actor model.AuthUser, id int64) (model.OperatorTask, error) {
	ownerID, err := companyOwner(actor, nil)
	if err != nil {
		return model.OperatorTask{}, err
	}
	task, err := operatorTaskService.operatorTaskRepository.FindByID(ctx, ownerID, id)
	return task, toOperatorTaskResponseError(err)
}

// Create adds the operator task after checking what it names, and returns it with its names.
func (operatorTaskService *OperatorTaskService) Create(ctx context.Context, taskRequest dto.CreateOperatorTaskRequest) (model.OperatorTask, error) {
	if err := operatorTaskService.checkReferences(ctx, taskRequest.OperatorTaskFields); err != nil {
		return model.OperatorTask{}, err
	}
	id, err := operatorTaskService.operatorTaskRepository.Create(ctx, operatorTaskFromFields(0, taskRequest.OperatorTaskFields))
	if err != nil {
		return model.OperatorTask{}, toOperatorTaskResponseError(err)
	}
	created, err := operatorTaskService.operatorTaskRepository.FindByID(ctx, nil, id)
	return created, toOperatorTaskResponseError(err)
}

// Update replaces the customer, machine type and task of operator task id.
func (operatorTaskService *OperatorTaskService) Update(ctx context.Context, id int64, taskRequest dto.UpdateOperatorTaskRequest) (model.OperatorTask, error) {
	if err := operatorTaskService.checkReferences(ctx, taskRequest.OperatorTaskFields); err != nil {
		return model.OperatorTask{}, err
	}
	if err := operatorTaskService.operatorTaskRepository.Update(ctx, operatorTaskFromFields(id, taskRequest.OperatorTaskFields)); err != nil {
		return model.OperatorTask{}, toOperatorTaskResponseError(err)
	}
	updated, err := operatorTaskService.operatorTaskRepository.FindByID(ctx, nil, id)
	return updated, toOperatorTaskResponseError(err)
}

// Delete soft-deletes the operator task.
func (operatorTaskService *OperatorTaskService) Delete(ctx context.Context, id int64) error {
	return toOperatorTaskResponseError(operatorTaskService.operatorTaskRepository.Delete(ctx, id))
}

// checkReferences checks the customer is a customer, the model (when given) is a catalog model
// and the task an active task of the list.
func (operatorTaskService *OperatorTaskService) checkReferences(ctx context.Context, fields dto.OperatorTaskFields) error {
	isCustomer, err := operatorTaskService.userRepository.IsCustomer(ctx, fields.UserID)
	if err != nil {
		return err
	}
	if !isCustomer {
		return response.SelectValidCustomer
	}
	if fields.ModelID != nil {
		_, err := operatorTaskService.machineModelService.GetByID(ctx, *fields.ModelID)
		// GetByID answers a response ("Model not found") for a missing model; anything else is a failure
		var apiResponse *response.Response
		switch {
		case errors.As(err, &apiResponse):
			return response.SelectValidMachineType
		case err != nil:
			return err
		}
	}
	return operatorTaskService.maintenanceTaskService.checkActive(ctx, fields.MaintenanceTaskID)
}

// operatorTaskFromFields returns the operator_tasks row of fields (id 0 for a new one).
func operatorTaskFromFields(id int64, fields dto.OperatorTaskFields) model.OperatorTask {
	return model.OperatorTask{ID: id, UserID: fields.UserID, ModelID: fields.ModelID, MaintenanceTaskID: fields.MaintenanceTaskID}
}

// toOperatorTaskResponseError maps repository errors to operator task responses. Other errors
// pass through and are sent as ServerError.
func toOperatorTaskResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.OperatorTaskNotFound
	case errors.Is(err, repository.ErrDuplicateName):
		return response.OperatorTaskExists
	case errors.Is(err, repository.ErrInvalidUser):
		return response.SelectValidCustomer
	case errors.Is(err, repository.ErrInvalidModel):
		return response.SelectValidMachineType
	case errors.Is(err, repository.ErrInvalidParent):
		return response.SelectValidTask
	default:
		return err
	}
}
