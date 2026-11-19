package dto

import (
	"regexp"
	"strings"

	"7th-sense-server/internal/model"
	"7th-sense-server/internal/response"
)

// PlcNodeIDMaxLength is the size of machine_variables.address.
const PlcNodeIDMaxLength = 255

// plcNodeIDPattern is an OPC UA node id: an optional namespace, then a string, numeric, GUID or
// opaque identifier, e.g. "ns=2;s=/Random/Plc/DB3300.DBX3.0".
var plcNodeIDPattern = regexp.MustCompile(`^(ns=\d+;)?[sigb]=\S.*$`)

// UpdatePlcNodeIDRequest is the body for PUT /machines/{id}/plc-variables/{variableId}: the node
// id the machine reads the variable from. The variable's default node id clears the override.
type UpdatePlcNodeIDRequest struct {
	NodeID string `json:"nodeId"`
}

// Validate trims the node id and checks its length and OPC UA form.
func (request *UpdatePlcNodeIDRequest) Validate() error {
	request.NodeID = strings.TrimSpace(request.NodeID)
	if request.NodeID == "" {
		return response.EnterPlcNodeID
	}
	if len(request.NodeID) > PlcNodeIDMaxLength {
		return response.PlcNodeIDTooLong
	}
	if !plcNodeIDPattern.MatchString(request.NodeID) {
		return response.InvalidPlcNodeID
	}
	return nil
}

// PLC node id statuses (MachinePlcVariableResponse.NodeStatus).
const (
	PlcNodeDefault = "default" // reads the default node id (cnc_variable_addresses)
	PlcNodeUpdated = "updated" // reads the machine's own node id (machine_variables)
)

// MachinePlcVariableResponse is one PLC variable of a machine: the node id the machine reads it
// from (nodeId), the default (defaultNodeId) and nodeStatus. isDefault = the machine has no node id
// of its own, so it reads the default.
type MachinePlcVariableResponse struct {
	VariableID    int    `json:"variableId"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Unit          string `json:"unit"`
	Category      string `json:"category"`
	DataType      string `json:"dataType"`
	Status        string `json:"status"`
	DefaultNodeID string `json:"defaultNodeId"`
	NodeID        string `json:"nodeId"`
	IsDefault     bool   `json:"isDefault"`
	NodeStatus    string `json:"nodeStatus"`
}

// NewMachinePlcVariableResponse builds the response row of variable.
func NewMachinePlcVariableResponse(variable model.MachinePlcVariable) MachinePlcVariableResponse {
	nodeID, status := variable.CustomNodeID, PlcNodeUpdated
	if nodeID == "" {
		nodeID, status = variable.DefaultNodeID, PlcNodeDefault
	}
	return MachinePlcVariableResponse{
		VariableID:    variable.ID,
		Code:          variable.Code,
		Name:          variable.Name,
		Unit:          variable.Unit,
		Category:      variable.Category,
		DataType:      variable.DataType,
		Status:        variable.Status,
		DefaultNodeID: variable.DefaultNodeID,
		NodeID:        nodeID,
		IsDefault:     variable.CustomNodeID == "",
		NodeStatus:    status,
	}
}

// MachinePlcVariablesResponse is GET /machines/{id}/plc-variables: every PLC variable (a plc
// address for any engine) with the node id the machine uses; defaultCount + updatedCount = all.
type MachinePlcVariablesResponse struct {
	MachineID             int64                        `json:"machineId"`
	MachineName           string                       `json:"machineName"`
	ControllerName        string                       `json:"controllerName"`
	ControllerVersionType *int                         `json:"controllerVersionType"`
	DefaultCount          int                          `json:"defaultCount"`
	UpdatedCount          int                          `json:"updatedCount"`
	Variables             []MachinePlcVariableResponse `json:"variables"`
}

// NewMachinePlcVariablesResponse builds the response for machine and its PLC variables.
func NewMachinePlcVariablesResponse(machine model.Machine, variables []model.MachinePlcVariable) MachinePlcVariablesResponse {
	list := MachinePlcVariablesResponse{
		MachineID:             machine.ID,
		MachineName:           machine.Name,
		ControllerName:        machine.ControllerName,
		ControllerVersionType: machine.ControllerVersionType,
		Variables:             make([]MachinePlcVariableResponse, len(variables)),
	}
	for i, variable := range variables {
		row := NewMachinePlcVariableResponse(variable)
		list.Variables[i] = row
		if row.IsDefault {
			list.DefaultCount++
		} else {
			list.UpdatedCount++
		}
	}
	return list
}
