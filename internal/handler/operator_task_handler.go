package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// OperatorTaskHandler handles HTTP requests for operator tasks: the maintenance tasks that apply
// to a customer's machines of a model. Only the superadmin writes them; everyone reads their
// company's (service.OperatorTaskService decides which).
type OperatorTaskHandler struct {
	service *service.OperatorTaskService
}

// NewOperatorTaskHandler creates an OperatorTaskHandler that uses operatorTaskService.
func NewOperatorTaskHandler(operatorTaskService *service.OperatorTaskService) *OperatorTaskHandler {
	return &OperatorTaskHandler{service: operatorTaskService}
}

// GetAll handles GET /api/v1/operator-tasks[?userId=<customer id>].
//
// @Summary      Get operator tasks
// @Description  Returns the operator tasks the user may read: the superadmin every customer's (?userId= narrows to one), everyone else their company's. Each row has the customer, the machine type (model with controller and brand; modelId null = every machine of the customer) and the task.
// @Description
// @Description  Used by pages: Operator tasks
// @Tags         Operator tasks
// @ID           getOperatorTasks
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "Customer (superadmin)"
// @Success      200 {object} response.Envelope{data=[]model.OperatorTask}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Another company's customer"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/operator-tasks [get]
func (operatorTaskHandler *OperatorTaskHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	customerID, err := queryID(request, "userId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	tasks, err := operatorTaskHandler.service.GetAll(request.Context(), actor, customerID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, tasks)
}

// GetByID handles GET /api/v1/operator-tasks/{id}.
//
// @Summary      Get operator task by id
// @Description  Returns the operator task when the user may read it (same rule as the list), else "Operator task not found".
// @Description
// @Description  Used by pages: Operator tasks
// @Tags         Operator tasks
// @ID           getOperatorTaskById
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Operator task id"
// @Success      200 {object} response.Envelope{data=model.OperatorTask}
// @Failure      400 {object} response.Envelope "Please enter valid id"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      404 {object} response.Envelope "Operator task not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/operator-tasks/{id} [get]
func (operatorTaskHandler *OperatorTaskHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	task, err := operatorTaskHandler.service.GetByID(request.Context(), actor, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, task)
}

// Create handles POST /api/v1/operator-tasks.
//
// @Summary      Add operator task
// @Description  Adds a task for a customer's machines of a model (modelId null = every machine of the customer). The customer must be a customer, the model a catalog model and the task an active task of the master list. A task is added once per customer and machine type.
// @Description  Roles: superadmin (others get 403).
// @Description  Validation: Please select customer; Please select valid customer; Please select valid machine type; Please select task; Please select valid task; This task is already added for this customer and machine type.
// @Description
// @Description  Used by pages: Operator tasks
// @Tags         Operator tasks
// @ID           createOperatorTask
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.CreateOperatorTaskRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.OperatorTask}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/operator-tasks [post]
func (operatorTaskHandler *OperatorTaskHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var taskRequest dto.CreateOperatorTaskRequest
	if err := decodeAndValidate(writer, request, &taskRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	task, err := operatorTaskHandler.service.Create(request.Context(), taskRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, task)
}

// Update handles PUT /api/v1/operator-tasks/{id}.
//
// @Summary      Update operator task
// @Description  Replaces the customer, machine type and task (same rules as Add).
// @Description  Roles: superadmin (others get 403).
// @Description  Validation: as Add operator task; Operator task not found.
// @Description
// @Description  Used by pages: Operator tasks
// @Tags         Operator tasks
// @ID           updateOperatorTask
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Operator task id"
// @Param        request body dto.UpdateOperatorTaskRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.OperatorTask}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Operator task not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/operator-tasks/{id} [put]
func (operatorTaskHandler *OperatorTaskHandler) Update(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	var taskRequest dto.UpdateOperatorTaskRequest
	if err := decodeAndValidate(writer, request, &taskRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	task, err := operatorTaskHandler.service.Update(request.Context(), id, taskRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, task)
}

// Delete handles DELETE /api/v1/operator-tasks/{id}.
//
// @Summary      Delete operator task
// @Description  Soft-deletes the operator task. The task stays in the master list.
// @Description  Roles: superadmin (others get 403).
// @Description
// @Description  Used by pages: Operator tasks
// @Tags         Operator tasks
// @ID           deleteOperatorTask
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Operator task id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Please enter valid id"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Operator task not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/operator-tasks/{id} [delete]
func (operatorTaskHandler *OperatorTaskHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	if err := operatorTaskHandler.service.Delete(request.Context(), id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Operator task deleted")
}
