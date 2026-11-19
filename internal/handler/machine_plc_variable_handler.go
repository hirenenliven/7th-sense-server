package handler

import (
	"net/http"
	"strconv"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// MachinePlcVariableHandler handles the superadmin's PLC node ids of one machine
// (service.MachinePlcVariableService).
type MachinePlcVariableHandler struct {
	service *service.MachinePlcVariableService
}

// NewMachinePlcVariableHandler creates a MachinePlcVariableHandler that uses plcVariableService.
func NewMachinePlcVariableHandler(plcVariableService *service.MachinePlcVariableService) *MachinePlcVariableHandler {
	return &MachinePlcVariableHandler{service: plcVariableService}
}

// GetAll handles GET /api/v1/machines/{id}/plc-variables.
//
// @Summary      List machine PLC node IDs
// @Description  Every PLC variable (a cnc_variables row with a plc address in cnc_variable_addresses for any CNC engine) with the OPC UA node id the machine reads it from. `defaultNodeId` = the plc address for the machine's engine and controller version, else the best one of another engine (PLC defaults exist for Siemens only) - the collector uses the same. `nodeStatus`: default (no node id of its own, reads the default) or updated (`nodeId` is the machine's own).
// @Description  Roles: superadmin (others get 403).
// @Description
// @Description  Used by pages: Machine PLC node IDs
// @Tags         Machines
// @ID           getMachinePlcVariables
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Machine id"
// @Success      200 {object} response.Envelope{data=dto.MachinePlcVariablesResponse}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machines/{id}/plc-variables [get]
func (plcVariableHandler *MachinePlcVariableHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
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

	variables, err := plcVariableHandler.service.GetAll(request.Context(), actor, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, variables)
}

// Update handles PUT /api/v1/machines/{id}/plc-variables/{variableId}.
//
// @Summary      Update machine PLC node ID
// @Description  Sets the OPC UA node id the machine reads one PLC variable from (`ns=<n>;s=…`, also i= / g= / b=, at most 255 characters). The variable's default node id resets it to the default instead of storing a copy. Stored in machine_variables (source plc); the machine's PLC worker (OPC UA, ADR 0021) subscribes again at once. PLC variables are read over OPC UA from Siemens machines only (PLC_ENABLED).
// @Description  Roles: superadmin (others get 403).
// @Description
// @Description  Used by pages: Machine PLC node IDs
// @Tags         Machines
// @ID           updateMachinePlcVariable
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Machine id"
// @Param        variableId path int true "Variable id (cnc_variables)"
// @Param        body body dto.UpdatePlcNodeIDRequest true "Node id"
// @Success      200 {object} response.Envelope{data=dto.MachinePlcVariableResponse}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machines/{id}/plc-variables/{variableId} [put]
func (plcVariableHandler *MachinePlcVariableHandler) Update(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, variableID, err := plcVariablePath(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	var body dto.UpdatePlcNodeIDRequest
	if err := decodeAndValidate(writer, request, &body); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	variable, err := plcVariableHandler.service.Update(request.Context(), actor, id, variableID, body)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, variable)
}

// Reset handles DELETE /api/v1/machines/{id}/plc-variables/{variableId}.
//
// @Summary      Reset machine PLC node ID
// @Description  The machine reads one PLC variable from its controller's default node id again (drops the machine's own node id; nothing happens when the default is already used).
// @Description  Roles: superadmin (others get 403).
// @Description
// @Description  Used by pages: Machine PLC node IDs
// @Tags         Machines
// @ID           resetMachinePlcVariable
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Machine id"
// @Param        variableId path int true "Variable id (cnc_variables)"
// @Success      200 {object} response.Envelope{data=dto.MachinePlcVariableResponse}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machines/{id}/plc-variables/{variableId} [delete]
func (plcVariableHandler *MachinePlcVariableHandler) Reset(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, variableID, err := plcVariablePath(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	variable, err := plcVariableHandler.service.Reset(request.Context(), actor, id, variableID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, variable)
}

// ResetAll handles DELETE /api/v1/machines/{id}/plc-variables.
//
// @Summary      Reset all machine PLC node IDs
// @Description  Every PLC variable of the machine reads its controller's default node id again. Returns the list as GET does.
// @Description  Roles: superadmin (others get 403).
// @Description
// @Description  Used by pages: Machine PLC node IDs
// @Tags         Machines
// @ID           resetMachinePlcVariables
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Machine id"
// @Success      200 {object} response.Envelope{data=dto.MachinePlcVariablesResponse}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machines/{id}/plc-variables [delete]
func (plcVariableHandler *MachinePlcVariableHandler) ResetAll(writer http.ResponseWriter, request *http.Request) {
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

	variables, err := plcVariableHandler.service.ResetAll(request.Context(), actor, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, variables)
}

// plcVariablePath reads the {id} (machine) and {variableId} path values as positive integers.
func plcVariablePath(request *http.Request) (int64, int, error) {
	id, err := pathID(request)
	if err != nil {
		return 0, 0, err
	}
	variableID, err := strconv.Atoi(request.PathValue("variableId"))
	if err != nil || variableID <= 0 {
		return 0, 0, response.NewError(http.StatusBadRequest, "Please enter valid variableId")
	}
	return id, variableID, nil
}
