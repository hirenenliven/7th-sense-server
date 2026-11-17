package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// MaintenanceTaskHandler handles HTTP requests for the maintenance task master list (the tasks an
// operator task picks from).
type MaintenanceTaskHandler struct {
	service *service.MaintenanceTaskService
}

// NewMaintenanceTaskHandler creates a MaintenanceTaskHandler that uses taskService.
func NewMaintenanceTaskHandler(taskService *service.MaintenanceTaskService) *MaintenanceTaskHandler {
	return &MaintenanceTaskHandler{service: taskService}
}

// GetAll handles GET /api/v1/maintenance-tasks.
//
// @Summary      Get maintenance tasks
// @Description  Returns the master list of operator tasks (shared by every customer), ordered by description, active and inactive, each with usageCount (the operator tasks that use it). Every logged-in user reads it.
// @Description
// @Description  Used by pages: Operator tasks, Task list
// @Tags         Operator tasks
// @ID           getMaintenanceTasks
// @Produce      json
// @Security     cookieAuth
// @Success      200 {object} response.Envelope{data=[]model.MaintenanceTask}
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/maintenance-tasks [get]
func (taskHandler *MaintenanceTaskHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	tasks, err := taskHandler.service.GetAll(request.Context())
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, tasks)
}

// Create handles POST /api/v1/maintenance-tasks: a new task for the master list.
//
// @Summary      Add maintenance task
// @Description  Adds an active task to the master list (Task List screen, or the Task field's New task on the Add Operator Task form). Spaces are trimmed and joined; the description is unique, ignoring case.
// @Description  Roles: superadmin (others get 403).
// @Description  Validation: Please enter task; Task must be at most 100 characters; This task already exists.
// @Description
// @Description  Used by pages: Operator tasks, Task list
// @Tags         Operator tasks
// @ID           createMaintenanceTask
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.CreateMaintenanceTaskRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.MaintenanceTask}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/maintenance-tasks [post]
func (taskHandler *MaintenanceTaskHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var taskRequest dto.CreateMaintenanceTaskRequest
	if err := decodeAndValidate(writer, request, &taskRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	task, err := taskHandler.service.Create(request.Context(), taskRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, task)
}

// Update handles PUT /api/v1/maintenance-tasks/{id}.
//
// @Summary      Update maintenance task
// @Description  Changes the description and status (active / inactive; blank keeps it). An inactive task is left out of the Task dropdown; operator tasks that use it keep it. The new description shows on every operator task that uses it.
// @Description  Roles: superadmin (others get 403).
// @Description  Validation: Please enter task; Task must be at most 100 characters; Please enter valid status; This task already exists; Task not found.
// @Description
// @Description  Used by pages: Task list
// @Tags         Operator tasks
// @ID           updateMaintenanceTask
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Task id"
// @Param        request body dto.UpdateMaintenanceTaskRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.MaintenanceTask}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Task not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/maintenance-tasks/{id} [put]
func (taskHandler *MaintenanceTaskHandler) Update(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	var taskRequest dto.UpdateMaintenanceTaskRequest
	if err := decodeAndValidate(writer, request, &taskRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	task, err := taskHandler.service.Update(request.Context(), id, taskRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, task)
}

// Delete handles DELETE /api/v1/maintenance-tasks/{id}.
//
// @Summary      Delete maintenance task
// @Description  Soft-deletes the task. A task used by operator tasks is refused (400): remove it from them first, or make it inactive.
// @Description  Roles: superadmin (others get 403).
// @Description
// @Description  Used by pages: Task list
// @Tags         Operator tasks
// @ID           deleteMaintenanceTask
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Task id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Task in use"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Task not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/maintenance-tasks/{id} [delete]
func (taskHandler *MaintenanceTaskHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	if err := taskHandler.service.Delete(request.Context(), id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Task deleted")
}
