package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// MaintenanceTaskService holds the rules of the maintenance task master list (server ADR 0015):
// every logged-in user reads it; the superadmin adds, edits and deletes tasks (superAdminOnly
// routes) on the Task List screen, or adds one from the Add Operator Task form. Descriptions are
// unique, ignoring case. A task used by operator tasks cannot be deleted; making it inactive takes
// it out of the Task dropdown instead.
type MaintenanceTaskService struct {
	taskRepository *repository.MaintenanceTaskRepository
}

// NewMaintenanceTaskService creates a MaintenanceTaskService that stores tasks with taskRepository.
func NewMaintenanceTaskService(taskRepository *repository.MaintenanceTaskRepository) *MaintenanceTaskService {
	return &MaintenanceTaskService{taskRepository: taskRepository}
}

// GetAll returns every task, ordered by description.
func (taskService *MaintenanceTaskService) GetAll(ctx context.Context) ([]model.MaintenanceTask, error) {
	return taskService.taskRepository.FindAll(ctx)
}

// Create adds an active task; a description already in the list answers MaintenanceTaskExists.
func (taskService *MaintenanceTaskService) Create(ctx context.Context, taskRequest dto.CreateMaintenanceTaskRequest) (model.MaintenanceTask, error) {
	task, err := taskService.taskRepository.Create(ctx, taskRequest.Description, dto.MaintenanceTaskStatusActive)
	return task, toMaintenanceTaskResponseError(err)
}

// Update changes the description and status (blank keeps it) of task id and returns it.
func (taskService *MaintenanceTaskService) Update(ctx context.Context, id int64, taskRequest dto.UpdateMaintenanceTaskRequest) (model.MaintenanceTask, error) {
	if err := taskService.taskRepository.Update(ctx, id, taskRequest.Description, taskRequest.Status); err != nil {
		return model.MaintenanceTask{}, toMaintenanceTaskResponseError(err)
	}
	task, err := taskService.taskRepository.FindByID(ctx, id)
	return task, toMaintenanceTaskResponseError(err)
}

// Delete soft-deletes task id unless operator tasks use it (they would lose their task).
func (taskService *MaintenanceTaskService) Delete(ctx context.Context, id int64) error {
	task, err := taskService.taskRepository.FindByID(ctx, id)
	if err != nil {
		return toMaintenanceTaskResponseError(err)
	}
	if task.UsageCount > 0 {
		return response.NewError(http.StatusBadRequest,
			fmt.Sprintf("This task is used by %d operator task(s). Remove it from them first, or make it inactive", task.UsageCount))
	}
	return toMaintenanceTaskResponseError(taskService.taskRepository.Delete(ctx, id))
}

// toMaintenanceTaskResponseError maps repository errors to maintenance task responses. Other
// errors pass through and are sent as ServerError.
func toMaintenanceTaskResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.MaintenanceTaskMissing
	case errors.Is(err, repository.ErrDuplicateName):
		return response.MaintenanceTaskExists
	default:
		return err
	}
}

// checkActive returns SelectValidTask unless id is a task of the list that is active.
func (taskService *MaintenanceTaskService) checkActive(ctx context.Context, id int64) error {
	task, err := taskService.taskRepository.FindByID(ctx, id)
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.SelectValidTask
	case err != nil:
		return err
	case task.Status != dto.MaintenanceTaskStatusActive:
		return response.SelectValidTask
	}
	return nil
}
