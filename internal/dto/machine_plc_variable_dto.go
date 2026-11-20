package dto

import (
	"regexp"
	"strings"
	"time"

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
	return checkPlcNodeID(request.NodeID)
}

// checkPlcNodeID checks a trimmed node id's length and OPC UA form.
func checkPlcNodeID(nodeID string) error {
	if nodeID == "" {
		return response.EnterPlcNodeID
	}
	if len(nodeID) > PlcNodeIDMaxLength {
		return response.PlcNodeIDTooLong
	}
	if !plcNodeIDPattern.MatchString(nodeID) {
		return response.InvalidPlcNodeID
	}
	return nil
}

// PlcNodeTestMaxNodes is the most node ids one test reads: the largest page of the client's list.
const PlcNodeTestMaxNodes = 100

// TestPlcNodeIDsRequest is the body for POST /machines/{id}/plc-variables/test: node ids to read
// once from the machine, saved or only typed.
type TestPlcNodeIDsRequest struct {
	NodeIDs []string `json:"nodeIds"`
}

// Validate trims every node id, checks each like UpdatePlcNodeIDRequest and drops repeats.
func (request *TestPlcNodeIDsRequest) Validate() error {
	if len(request.NodeIDs) == 0 {
		return response.EnterPlcNodeID
	}
	if len(request.NodeIDs) > PlcNodeTestMaxNodes {
		return response.TooManyPlcNodeIDs
	}
	seen := make(map[string]bool, len(request.NodeIDs))
	nodeIDs := make([]string, 0, len(request.NodeIDs))
	for _, nodeID := range request.NodeIDs {
		nodeID = strings.TrimSpace(nodeID)
		if err := checkPlcNodeID(nodeID); err != nil {
			return err
		}
		if !seen[nodeID] {
			seen[nodeID] = true
			nodeIDs = append(nodeIDs, nodeID)
		}
	}
	request.NodeIDs = nodeIDs
	return nil
}

// PlcNodeTestResponse is POST /machines/{id}/plc-variables/test: whether the server reached the
// machine's OPC UA server (tested = false when it did not try, e.g. not a Siemens controller) and
// each node id's value or the controller's refusal. message / hint say the outcome in words.
type PlcNodeTestResponse struct {
	Endpoint       string              `json:"endpoint"`
	SecurityMode   string              `json:"securityMode"`
	SecurityPolicy string              `json:"securityPolicy"`
	Tested         bool                `json:"tested"`
	Connected      bool                `json:"connected"`
	Message        string              `json:"message"`
	Hint           string              `json:"hint"`
	TestedAt       time.Time           `json:"testedAt"`
	DurationMs     int64               `json:"durationMs"`
	Results        []PlcNodeTestResult `json:"results"`
}

// PlcNodeTestResult is one node id read once: ok with its value (a number, booleans 1 / 0, or
// text) and the controller's data type, else status = why not (e.g. BadNodeIdUnknown).
type PlcNodeTestResult struct {
	NodeID     string     `json:"nodeId"`
	OK         bool       `json:"ok"`
	Status     string     `json:"status"`
	Value      any        `json:"value"`
	ValueType  string     `json:"valueType"`
	SourceTime *time.Time `json:"sourceTime"`
}

// PlcReaderResponse is what the machine's OPC UA PLC reader does (ADR 0021): enabled =
// PLC_ENABLED, supported = a Siemens controller, state = its worker's (connecting, connected,
// failed, no_variables; "" when no worker reports the machine) and its last error.
type PlcReaderResponse struct {
	Enabled   bool   `json:"enabled"`
	Supported bool   `json:"supported"`
	State     string `json:"state"`
	LastError string `json:"lastError"`
}

// PLC node id statuses (MachinePlcVariableResponse.NodeStatus).
const (
	PlcNodeDefault = "default" // reads the default node id (cnc_variable_addresses)
	PlcNodeUpdated = "updated" // reads the machine's own node id (machine_variables)
)

// MachinePlcVariableResponse is one PLC variable of a machine: the node id the machine reads it
// from (nodeId), the default (defaultNodeId) and nodeStatus. isDefault = the machine has no node id
// of its own, so it reads the default. readerNodeId is the node id the PLC reader subscribed to
// ("" = none; it differs from nodeId until the reader subscribes again after a save) and
// readerStatus the controller's refusal of it ("" = accepted).
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
	ReaderNodeID  string `json:"readerNodeId"`
	ReaderStatus  string `json:"readerStatus"`
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
// reader is what the machine's PLC reader does.
type MachinePlcVariablesResponse struct {
	MachineID             int64                        `json:"machineId"`
	MachineName           string                       `json:"machineName"`
	ControllerName        string                       `json:"controllerName"`
	ControllerVersionType *int                         `json:"controllerVersionType"`
	DefaultCount          int                          `json:"defaultCount"`
	UpdatedCount          int                          `json:"updatedCount"`
	Reader                PlcReaderResponse            `json:"reader"`
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
